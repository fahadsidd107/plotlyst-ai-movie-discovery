package debate

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const MinimumReviews = 5

type ReviewRecord struct {
	ID         bson.ObjectID `bson:"_id"`
	UserID     string        `bson:"user_id"`
	AuthorName string        `bson:"author_name"`
	Body       string        `bson:"body"`
	Spoiler    bool          `bson:"spoiler"`
	CreatedAt  time.Time     `bson:"created_at"`
	UpdatedAt  time.Time     `bson:"updated_at"`
	Stars      int           `bson:"stars"`
}

type Evidence struct {
	ReviewID      string    `json:"review_id" bson:"review_id"`
	AuthorName    string    `json:"author_name" bson:"author_name"`
	Stars         int       `json:"stars" bson:"stars"`
	Excerpt       string    `json:"excerpt,omitempty" bson:"excerpt,omitempty"`
	SpoilerHidden bool      `json:"spoiler_hidden" bson:"spoiler_hidden"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
}

type Observation struct {
	Theme        string     `json:"theme" bson:"theme"`
	Summary      string     `json:"summary" bson:"summary"`
	MentionCount int        `json:"mention_count" bson:"mention_count"`
	Evidence     []Evidence `json:"evidence" bson:"evidence"`
}

type Polarization struct {
	Score     float64 `json:"score" bson:"score"`
	Polarized bool    `json:"polarized" bson:"polarized"`
	Summary   string  `json:"summary" bson:"summary"`
}

type Summary struct {
	MovieID             int            `json:"movie_id" bson:"movie_id"`
	Status              string         `json:"status" bson:"status"`
	ReviewCount         int            `json:"review_count" bson:"review_count"`
	MinimumReviews      int            `json:"minimum_reviews" bson:"minimum_reviews"`
	RatingDistribution  map[string]int `json:"rating_distribution" bson:"rating_distribution"`
	Polarization        Polarization   `json:"polarization" bson:"polarization"`
	Loved               []Observation  `json:"loved" bson:"loved"`
	Disliked            []Observation  `json:"disliked" bson:"disliked"`
	AnalysisMethod      string         `json:"analysis_method" bson:"analysis_method"`
	GeneratedAt         time.Time      `json:"generated_at" bson:"generated_at"`
	AnalyzedReviewCount int            `json:"analyzed_review_count" bson:"analyzed_review_count"`
	SourceReviewCount   int            `json:"-" bson:"source_review_count"`
	SourceUpdatedAt     time.Time      `json:"source_updated_at" bson:"source_updated_at"`
	NewReviewsPending   int            `json:"new_reviews_pending" bson:"-"`
}
