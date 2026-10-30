package debate

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Service struct {
	Reviews, Ratings, Cache, AIPrivacy *mongo.Collection
	Local                              Analyzer
	AI                                 Analyzer
	Timeout                            time.Duration
	mu                                 sync.Mutex
	running                            map[int]bool
}

func (s *Service) Get(ctx context.Context, movieID int, bucket string) (Summary, error) {
	rows, err := s.reviewRows(ctx, movieID)
	if err != nil {
		return Summary{}, err
	}
	distribution, err := s.ratingDistribution(ctx, movieID)
	if err != nil {
		return Summary{}, err
	}
	base := Summary{MovieID: movieID, ReviewCount: len(rows), MinimumReviews: MinimumReviews, RatingDistribution: distribution, Polarization: Polarize(distribution)}
	if len(rows) < MinimumReviews {
		base.Status = "insufficient"
		base.AnalysisMethod = "none"
		return base, nil
	}
	var cached Summary
	if s.Cache.FindOne(ctx, bson.M{"movie_id": movieID}).Decode(&cached) == nil {
		pending := len(rows) - cached.SourceReviewCount
		if pending < 0 {
			pending = 0
		}
		cached.ReviewCount = len(rows)
		cached.RatingDistribution = distribution
		cached.Polarization = base.Polarization
		cached.NewReviewsPending = pending
		changed := len(rows) <= cached.SourceReviewCount && latestUpdate(rows).After(cached.SourceUpdatedAt)
		if pending < 3 && !changed {
			cached.Status = "available"
			return Filter(cached, bucket), nil
		}
		cached.Status = "stale"
		s.refreshAsync(movieID, rows, distribution)
		return Filter(cached, bucket), nil
	}
	created, err := s.analyze(ctx, movieID, rows, distribution)
	if err != nil {
		return Summary{}, err
	}
	return Filter(created, bucket), nil
}

func (s *Service) analyze(ctx context.Context, movieID int, rows []ReviewRecord, distribution map[string]int) (Summary, error) {
	analyzer := s.Local
	if analyzer == nil {
		analyzer = DeterministicAnalyzer{}
	}
	analysisRows := rows
	if s.AI != nil {
		eligible := s.aiEligible(ctx, rows)
		if len(eligible) >= MinimumReviews {
			analyzer = s.AI
			analysisRows = eligible
		}
	}
	loved, disliked, method, err := analyzer.Analyze(ctx, analysisRows)
	if err != nil {
		return Summary{}, err
	}
	summary := Summary{MovieID: movieID, Status: "available", ReviewCount: len(rows), MinimumReviews: MinimumReviews, RatingDistribution: distribution, Polarization: Polarize(distribution), Loved: loved, Disliked: disliked, AnalysisMethod: method, GeneratedAt: time.Now().UTC(), AnalyzedReviewCount: len(analysisRows), SourceReviewCount: len(rows), SourceUpdatedAt: latestUpdate(rows)}
	_, err = s.Cache.UpdateOne(ctx, bson.M{"movie_id": movieID}, bson.M{"$set": summary}, options.UpdateOne().SetUpsert(true))
	return summary, err
}

func (s *Service) refreshAsync(movieID int, rows []ReviewRecord, distribution map[string]int) {
	s.mu.Lock()
	if s.running == nil {
		s.running = map[int]bool{}
	}
	if s.running[movieID] {
		s.mu.Unlock()
		return
	}
	s.running[movieID] = true
	s.mu.Unlock()
	go func() {
		defer func() { s.mu.Lock(); delete(s.running, movieID); s.mu.Unlock() }()
		timeout := s.Timeout
		if timeout <= 0 {
			timeout = 15 * time.Second
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, _ = s.analyze(ctx, movieID, rows, distribution)
	}()
}

func (s *Service) reviewRows(ctx context.Context, movieID int) ([]ReviewRecord, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"media_type": "movie", "media_id": movieID, "status": "VISIBLE", "parent_id": nil, "deleted_at": nil}}},
		{{Key: "$lookup", Value: bson.M{"from": "ratings", "let": bson.M{"review_user": "$user_id"}, "pipeline": bson.A{bson.M{"$match": bson.M{"$expr": bson.M{"$and": bson.A{bson.M{"$eq": bson.A{"$user_id", "$$review_user"}}, bson.M{"$eq": bson.A{"$media_type", "movie"}}, bson.M{"$eq": bson.A{"$media_id", movieID}}}}}}}, "as": "rating"}}},
		{{Key: "$set", Value: bson.M{"stars": bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$rating.stars", 0}}, 0}}}}},
		{{Key: "$project", Value: bson.M{"rating": 0}}}, {{Key: "$sort", Value: bson.M{"created_at": -1}}},
	}
	cursor, err := s.Reviews.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	rows := []ReviewRecord{}
	err = cursor.All(ctx, &rows)
	return rows, err
}

func (s *Service) ratingDistribution(ctx context.Context, movieID int) (map[string]int, error) {
	pipeline := mongo.Pipeline{{{Key: "$match", Value: bson.M{"media_type": "movie", "media_id": movieID}}}, {{Key: "$group", Value: bson.M{"_id": "$stars", "count": bson.M{"$sum": 1}}}}}
	cursor, err := s.Ratings.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var rows []struct {
		ID    int `bson:"_id"`
		Count int `bson:"count"`
	}
	if err = cursor.All(ctx, &rows); err != nil {
		return nil, err
	}
	out := map[string]int{"1": 0, "2": 0, "3": 0, "4": 0, "5": 0}
	for _, row := range rows {
		out[fmt.Sprint(row.ID)] = row.Count
	}
	return out, nil
}

func (s *Service) aiEligible(ctx context.Context, rows []ReviewRecord) []ReviewRecord {
	ids := []string{}
	for _, row := range rows {
		ids = append(ids, row.UserID)
	}
	cursor, err := s.AIPrivacy.Find(ctx, bson.M{"user_id": bson.M{"$in": ids}, "ai_enabled": true, "allow_activity_for_explanations": true})
	if err != nil {
		return nil
	}
	defer cursor.Close(ctx)
	var policies []struct {
		UserID string `bson:"user_id"`
	}
	_ = cursor.All(ctx, &policies)
	allowed := map[string]bool{}
	for _, policy := range policies {
		allowed[policy.UserID] = true
	}
	eligible := []ReviewRecord{}
	for _, row := range rows {
		if allowed[row.UserID] && !row.Spoiler {
			eligible = append(eligible, row)
		}
	}
	return eligible
}

func latestUpdate(rows []ReviewRecord) time.Time {
	var latest time.Time
	for _, row := range rows {
		candidate := row.UpdatedAt
		if candidate.IsZero() {
			candidate = row.CreatedAt
		}
		if candidate.After(latest) {
			latest = candidate
		}
	}
	return latest
}
