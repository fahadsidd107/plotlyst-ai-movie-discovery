package controllers

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"plotlyst/server/apiresponse"
	"plotlyst/server/database"
	"plotlyst/server/models"
	"plotlyst/server/utils"
)

type CommunityController struct{ client *mongo.Client }

func NewCommunityController(client *mongo.Client) *CommunityController {
	return &CommunityController{client: client}
}
func (cc *CommunityController) collection(name string) *mongo.Collection {
	return database.OpenCollection(name, cc.client)
}
func requestContext(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), 8*time.Second)
}

func cleanText(value string, max int) string {
	value = strings.TrimSpace(strings.Map(func(r rune) rune {
		if r < 32 && r != '\n' && r != '\t' {
			return -1
		}
		return r
	}, value))
	runes := []rune(value)
	if len(runes) > max {
		value = string(runes[:max])
	}
	return value
}
func validMedia(kind string, id int) bool { return (kind == "movie" || kind == "tv") && id > 0 }
func parseMedia(c *gin.Context) (models.MediaRef, bool) {
	var ref models.MediaRef
	ref.MediaType = c.Param("mediaType")
	if _, err := fmtSscanf(c.Param("mediaId"), &ref.MediaID); err != nil || !validMedia(ref.MediaType, ref.MediaID) {
		apiresponse.Error(c, 400, "invalid_media", "A valid movie or television identifier is required")
		return ref, false
	}
	return ref, true
}

// fmtSscanf is kept tiny to make route validation independently testable.
func fmtSscanf(raw string, target *int) (int, error) { return fmt.Sscanf(raw, "%d", target) }

func (cc *CommunityController) Session(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	ctx, cancel := requestContext(c)
	defer cancel()
	var user models.User
	if err := cc.collection("users").FindOne(ctx, bson.M{"user_id": uid}).Decode(&user); err != nil {
		apiresponse.Error(c, 401, "session_invalid", "The session user no longer exists")
		return
	}
	c.JSON(200, models.UserResponse{UserId: user.UserID, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Role: user.Role, FavouriteGenres: user.FavouriteGenres})
}

func (cc *CommunityController) GetPublicProfile(c *gin.Context) {
	ctx, cancel := requestContext(c)
	defer cancel()
	var profile models.Profile
	if err := cc.collection("profiles").FindOne(ctx, bson.M{"username": strings.ToLower(c.Param("username")), "privacy.profile": bson.M{"$ne": "PRIVATE"}}).Decode(&profile); err != nil {
		apiresponse.Error(c, 404, "profile_not_found", "Profile not found")
		return
	}
	if profile.Privacy.Activity == "PRIVATE" {
		profile.FavouriteGenres = nil
	}
	c.JSON(200, gin.H{"data": profile})
}
func (cc *CommunityController) GetMyProfile(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	ctx, cancel := requestContext(c)
	defer cancel()
	var p models.Profile
	if err := cc.collection("profiles").FindOne(ctx, bson.M{"user_id": uid}).Decode(&p); err != nil {
		apiresponse.Error(c, 404, "profile_not_found", "Create your profile first")
		return
	}
	c.JSON(200, gin.H{"data": p})
}
func (cc *CommunityController) UpsertMyProfile(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	var p models.Profile
	if err := c.ShouldBindJSON(&p); err != nil {
		apiresponse.Error(c, 400, "invalid_request", "Profile payload is invalid")
		return
	}
	p.UserID = uid
	p.Username = strings.ToLower(cleanText(p.Username, 30))
	p.DisplayName = cleanText(p.DisplayName, 80)
	p.Bio = cleanText(p.Bio, 500)
	p.UpdatedAt = time.Now().UTC()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = p.UpdatedAt
	}
	if ok, _ := regexp.MatchString(`^[a-z0-9_]{3,30}$`, p.Username); !ok {
		apiresponse.Error(c, 400, "invalid_username", "Use 3–30 lowercase letters, numbers or underscores")
		return
	}
	if p.Region == "" {
		p.Region = "DE"
	}
	if p.Locale == "" {
		p.Locale = "en"
	}
	if p.Privacy.Profile == "" {
		p.Privacy = models.PrivacySettings{Profile: "PUBLIC", Library: "PRIVATE", Activity: "PRIVATE"}
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	_, err := cc.collection("profiles").UpdateOne(ctx, bson.M{"user_id": uid}, bson.M{"$set": p}, options.UpdateOne().SetUpsert(true))
	if mongo.IsDuplicateKeyError(err) {
		apiresponse.Error(c, 409, "username_taken", "That username is already in use")
		return
	}
	if err != nil {
		apiresponse.Error(c, 500, "database_error", "Unable to save profile")
		return
	}
	c.JSON(200, gin.H{"data": p})
}

func (cc *CommunityController) PutRating(c *gin.Context) {
	ref, ok := parseMedia(c)
	if !ok {
		return
	}
	uid, _ := utils.GetUserIdFromContext(c)
	var body struct {
		Stars int `json:"stars"`
	}
	if c.ShouldBindJSON(&body) != nil || body.Stars < 1 || body.Stars > 5 {
		apiresponse.Error(c, 400, "invalid_rating", "Stars must be between 1 and 5")
		return
	}
	now := time.Now().UTC()
	ctx, cancel := requestContext(c)
	defer cancel()
	_, err := cc.collection("ratings").UpdateOne(ctx, bson.M{"user_id": uid, "media_type": ref.MediaType, "media_id": ref.MediaID}, bson.M{"$set": bson.M{"stars": body.Stars, "updated_at": now}, "$setOnInsert": bson.M{"created_at": now}}, options.UpdateOne().SetUpsert(true))
	if err != nil {
		apiresponse.Error(c, 500, "database_error", "Unable to save rating")
		return
	}
	c.JSON(200, gin.H{"stars": body.Stars})
}
func (cc *CommunityController) DeleteRating(c *gin.Context) {
	ref, ok := parseMedia(c)
	if !ok {
		return
	}
	uid, _ := utils.GetUserIdFromContext(c)
	ctx, cancel := requestContext(c)
	defer cancel()
	cc.collection("ratings").DeleteOne(ctx, bson.M{"user_id": uid, "media_type": ref.MediaType, "media_id": ref.MediaID})
	c.Status(204)
}
func (cc *CommunityController) RatingSummary(c *gin.Context) {
	ref, ok := parseMedia(c)
	if !ok {
		return
	}
	uid, _ := utils.GetUserIdFromContext(c)
	ctx, cancel := requestContext(c)
	defer cancel()
	cur, err := cc.collection("ratings").Find(ctx, bson.M{"media_type": ref.MediaType, "media_id": ref.MediaID})
	if err != nil {
		apiresponse.Error(c, 500, "database_error", "Unable to load ratings")
		return
	}
	defer cur.Close(ctx)
	dist := map[string]int{"1": 0, "2": 0, "3": 0, "4": 0, "5": 0}
	count, total, viewer := 0, 0, 0
	for cur.Next(ctx) {
		var r models.Rating
		if cur.Decode(&r) == nil {
			count++
			total += r.Stars
			dist[strconv.Itoa(r.Stars)]++
			if r.UserID == uid {
				viewer = r.Stars
			}
		}
	}
	avg := 0.0
	if count > 0 {
		avg = float64(total) / float64(count)
	}
	c.JSON(200, gin.H{"average": avg, "count": count, "distribution": dist, "viewer_rating": viewer})
}

func (cc *CommunityController) ListReviews(c *gin.Context) {
	ref, ok := parseMedia(c)
	if !ok {
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	cur, err := cc.collection("reviews").Find(ctx, bson.M{"media_type": ref.MediaType, "media_id": ref.MediaID, "status": "VISIBLE", "parent_id": nil}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(100))
	if err != nil {
		apiresponse.Error(c, 500, "database_error", "Unable to load reviews")
		return
	}
	defer cur.Close(ctx)
	rows := []models.Review{}
	cur.All(ctx, &rows)
	c.JSON(200, gin.H{"data": rows})
}
func (cc *CommunityController) CreateReview(c *gin.Context) {
	ref, ok := parseMedia(c)
	if !ok {
		return
	}
	uid, _ := utils.GetUserIdFromContext(c)
	var body models.Review
	if c.ShouldBindJSON(&body) != nil {
		apiresponse.Error(c, 400, "invalid_request", "Review payload is invalid")
		return
	}
	body.ID = bson.NewObjectID()
	body.UserID = uid
	body.MediaRef = ref
	body.Title = cleanText(body.Title, 120)
	body.Body = cleanText(body.Body, 5000)
	if len([]rune(body.Body)) < 10 {
		apiresponse.Error(c, 400, "review_too_short", "Review must be at least 10 characters")
		return
	}
	body.Status = "VISIBLE"
	body.CreatedAt = time.Now().UTC()
	body.UpdatedAt = body.CreatedAt
	ctx, cancel := requestContext(c)
	defer cancel()
	if _, err := cc.collection("reviews").InsertOne(ctx, body); err != nil {
		apiresponse.Error(c, 500, "database_error", "Unable to publish review")
		return
	}
	c.JSON(201, gin.H{"data": body})
}

func (cc *CommunityController) CreateReply(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	parent, err := bson.ObjectIDFromHex(c.Param("reviewId"))
	if err != nil {
		apiresponse.Error(c, 400, "invalid_review", "Review identifier is invalid")
		return
	}
	var body struct {
		Body    string `json:"body"`
		Spoiler bool   `json:"spoiler"`
	}
	if c.ShouldBindJSON(&body) != nil {
		apiresponse.Error(c, 400, "invalid_request", "Reply payload is invalid")
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	var root models.Review
	if cc.collection("reviews").FindOne(ctx, bson.M{"_id": parent, "parent_id": nil, "status": "VISIBLE"}).Decode(&root) != nil {
		apiresponse.Error(c, 404, "review_not_found", "Review not found")
		return
	}
	reply := models.Review{ID: bson.NewObjectID(), UserID: uid, MediaRef: root.MediaRef, ParentID: &parent, Body: cleanText(body.Body, 2000), Spoiler: body.Spoiler, Status: "VISIBLE", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	if len([]rune(reply.Body)) < 2 {
		apiresponse.Error(c, 400, "reply_too_short", "Reply is too short")
		return
	}
	cc.collection("reviews").InsertOne(ctx, reply)
	c.JSON(201, gin.H{"data": reply})
}
func (cc *CommunityController) UpdateReview(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	id, err := bson.ObjectIDFromHex(c.Param("reviewId"))
	if err != nil {
		apiresponse.Error(c, 400, "invalid_review", "Review identifier is invalid")
		return
	}
	var body struct {
		Title   string `json:"title"`
		Body    string `json:"body"`
		Spoiler bool   `json:"spoiler"`
	}
	if c.ShouldBindJSON(&body) != nil {
		apiresponse.Error(c, 400, "invalid_request", "Review payload is invalid")
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	result, _ := cc.collection("reviews").UpdateOne(ctx, bson.M{"_id": id, "user_id": uid, "deleted_at": nil}, bson.M{"$set": bson.M{"title": cleanText(body.Title, 120), "body": cleanText(body.Body, 5000), "spoiler": body.Spoiler, "updated_at": time.Now().UTC()}})
	if result.MatchedCount == 0 {
		apiresponse.Error(c, 404, "review_not_found", "Review was not found or is not yours")
		return
	}
	c.JSON(200, gin.H{"message": "Review updated"})
}
func (cc *CommunityController) DeleteReview(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	id, err := bson.ObjectIDFromHex(c.Param("reviewId"))
	if err != nil {
		apiresponse.Error(c, 400, "invalid_review", "Review identifier is invalid")
		return
	}
	now := time.Now().UTC()
	ctx, cancel := requestContext(c)
	defer cancel()
	res, _ := cc.collection("reviews").UpdateOne(ctx, bson.M{"_id": id, "user_id": uid, "deleted_at": nil}, bson.M{"$set": bson.M{"deleted_at": now, "status": "REMOVED", "body": "", "title": ""}})
	if res.MatchedCount == 0 {
		apiresponse.Error(c, 404, "review_not_found", "Review was not found or is not yours")
		return
	}
	c.Status(204)
}

func (cc *CommunityController) ToggleHelpful(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	id, err := bson.ObjectIDFromHex(c.Param("reviewId"))
	if err != nil {
		apiresponse.Error(c, 400, "invalid_review", "Review identifier is invalid")
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	var review models.Review
	if cc.collection("reviews").FindOne(ctx, bson.M{"_id": id}).Decode(&review) != nil {
		apiresponse.Error(c, 404, "review_not_found", "Review not found")
		return
	}
	if review.UserID == uid {
		apiresponse.Error(c, 409, "self_vote", "You cannot vote for your own review")
		return
	}
	_, err = cc.collection("review_votes").InsertOne(ctx, models.ReviewVote{ReviewID: id, UserID: uid, CreatedAt: time.Now().UTC()})
	if mongo.IsDuplicateKeyError(err) {
		apiresponse.Error(c, 409, "already_voted", "You already marked this review helpful")
		return
	}
	cc.collection("reviews").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"helpful_count": 1}})
	c.Status(204)
}
func (cc *CommunityController) RemoveHelpful(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	id, err := bson.ObjectIDFromHex(c.Param("reviewId"))
	if err != nil {
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	res, _ := cc.collection("review_votes").DeleteOne(ctx, bson.M{"review_id": id, "user_id": uid})
	if res.DeletedCount > 0 {
		cc.collection("reviews").UpdateOne(ctx, bson.M{"_id": id, "helpful_count": bson.M{"$gt": 0}}, bson.M{"$inc": bson.M{"helpful_count": -1}})
	}
	c.Status(204)
}
func (cc *CommunityController) ReportReview(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	id, err := bson.ObjectIDFromHex(c.Param("reviewId"))
	if err != nil {
		apiresponse.Error(c, 400, "invalid_review", "Review identifier is invalid")
		return
	}
	var b struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&b)
	b.Reason = cleanText(b.Reason, 500)
	if len(b.Reason) < 3 {
		apiresponse.Error(c, 400, "invalid_reason", "A report reason is required")
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	_, err = cc.collection("review_reports").InsertOne(ctx, models.ReviewReport{ID: bson.NewObjectID(), ReviewID: id, ReporterUserID: uid, Reason: b.Reason, Status: "OPEN", CreatedAt: time.Now().UTC()})
	if mongo.IsDuplicateKeyError(err) {
		apiresponse.Error(c, 409, "already_reported", "You already reported this review")
		return
	}
	c.Status(204)
}

func (cc *CommunityController) ListMine(collection string) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, _ := utils.GetUserIdFromContext(c)
		ctx, cancel := requestContext(c)
		defer cancel()
		cur, err := cc.collection(collection).Find(ctx, bson.M{"user_id": uid}, options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}}))
		if err != nil {
			apiresponse.Error(c, 500, "database_error", "Unable to load library")
			return
		}
		defer cur.Close(ctx)
		rows := []models.LibraryItem{}
		cur.All(ctx, &rows)
		c.JSON(200, gin.H{"data": rows})
	}
}
func (cc *CommunityController) PutLibrary(collection string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ref, ok := parseMedia(c)
		if !ok {
			return
		}
		uid, _ := utils.GetUserIdFromContext(c)
		var item models.LibraryItem
		c.ShouldBindJSON(&item)
		item.UserID = uid
		item.MediaRef = ref
		item.Title = cleanText(item.Title, 200)
		now := time.Now().UTC()
		item.UpdatedAt = now
		item.CreatedAt = now
		if collection == "watch_history" {
			item.Status = "WATCHED"
			item.WatchedAt = &now
		}
		ctx, cancel := requestContext(c)
		defer cancel()
		_, err := cc.collection(collection).UpdateOne(ctx, bson.M{"user_id": uid, "media_type": ref.MediaType, "media_id": ref.MediaID}, bson.M{"$set": item}, options.UpdateOne().SetUpsert(true))
		if err != nil {
			apiresponse.Error(c, 500, "database_error", "Unable to update library")
			return
		}
		c.JSON(200, gin.H{"data": item})
	}
}
func (cc *CommunityController) DeleteLibrary(collection string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ref, ok := parseMedia(c)
		if !ok {
			return
		}
		uid, _ := utils.GetUserIdFromContext(c)
		ctx, cancel := requestContext(c)
		defer cancel()
		cc.collection(collection).DeleteOne(ctx, bson.M{"user_id": uid, "media_type": ref.MediaType, "media_id": ref.MediaID})
		c.Status(204)
	}
}

func (cc *CommunityController) ListPlaylists(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	ctx, cancel := requestContext(c)
	defer cancel()
	cur, _ := cc.collection("playlists").Find(ctx, bson.M{"owner_user_id": uid})
	rows := []models.Playlist{}
	if cur != nil {
		defer cur.Close(ctx)
		cur.All(ctx, &rows)
	}
	c.JSON(200, gin.H{"data": rows})
}
func (cc *CommunityController) CreatePlaylist(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	var p models.Playlist
	if c.ShouldBindJSON(&p) != nil {
		apiresponse.Error(c, 400, "invalid_request", "Playlist payload is invalid")
		return
	}
	p.ID = bson.NewObjectID()
	p.OwnerUserID = uid
	p.Name = cleanText(p.Name, 80)
	p.Description = cleanText(p.Description, 500)
	if p.Name == "" {
		apiresponse.Error(c, 400, "name_required", "Playlist name is required")
		return
	}
	p.CreatedAt = time.Now().UTC()
	p.UpdatedAt = p.CreatedAt
	ctx, cancel := requestContext(c)
	defer cancel()
	cc.collection("playlists").InsertOne(ctx, p)
	c.JSON(201, gin.H{"data": p})
}
func (cc *CommunityController) DeletePlaylist(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	id, err := bson.ObjectIDFromHex(c.Param("playlistId"))
	if err != nil {
		apiresponse.Error(c, 400, "invalid_playlist", "Playlist identifier is invalid")
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	res, _ := cc.collection("playlists").DeleteOne(ctx, bson.M{"_id": id, "owner_user_id": uid})
	if res.DeletedCount == 0 {
		apiresponse.Error(c, 404, "playlist_not_found", "Playlist not found or not yours")
		return
	}
	cc.collection("playlist_items").DeleteMany(ctx, bson.M{"playlist_id": id})
	c.Status(204)
}
func (cc *CommunityController) GetPlaylist(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	id, err := bson.ObjectIDFromHex(c.Param("playlistId"))
	if err != nil {
		apiresponse.Error(c, 400, "invalid_playlist", "Playlist identifier is invalid")
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	var p models.Playlist
	if cc.collection("playlists").FindOne(ctx, bson.M{"_id": id, "owner_user_id": uid}).Decode(&p) != nil {
		apiresponse.Error(c, 404, "playlist_not_found", "Playlist not found or not yours")
		return
	}
	cur, _ := cc.collection("playlist_items").Find(ctx, bson.M{"playlist_id": id}, options.Find().SetSort(bson.D{{Key: "position", Value: 1}}))
	items := []models.PlaylistItem{}
	if cur != nil {
		defer cur.Close(ctx)
		cur.All(ctx, &items)
	}
	c.JSON(200, gin.H{"data": p, "items": items})
}
func (cc *CommunityController) UpdatePlaylist(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	id, err := bson.ObjectIDFromHex(c.Param("playlistId"))
	if err != nil {
		return
	}
	var p models.Playlist
	if c.ShouldBindJSON(&p) != nil {
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	res, _ := cc.collection("playlists").UpdateOne(ctx, bson.M{"_id": id, "owner_user_id": uid}, bson.M{"$set": bson.M{"name": cleanText(p.Name, 80), "description": cleanText(p.Description, 500), "private": p.Private, "updated_at": time.Now().UTC()}})
	if res.MatchedCount == 0 {
		apiresponse.Error(c, 404, "playlist_not_found", "Playlist not found or not yours")
		return
	}
	c.Status(204)
}
func (cc *CommunityController) DeletePlaylistItem(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	pid, e1 := bson.ObjectIDFromHex(c.Param("playlistId"))
	iid, e2 := bson.ObjectIDFromHex(c.Param("itemId"))
	if e1 != nil || e2 != nil {
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	if cc.collection("playlists").FindOne(ctx, bson.M{"_id": pid, "owner_user_id": uid}).Err() != nil {
		apiresponse.Error(c, 404, "playlist_not_found", "Playlist not found or not yours")
		return
	}
	cc.collection("playlist_items").DeleteOne(ctx, bson.M{"_id": iid, "playlist_id": pid})
	c.Status(204)
}
func (cc *CommunityController) AddPlaylistItem(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	pid, err := bson.ObjectIDFromHex(c.Param("playlistId"))
	if err != nil {
		apiresponse.Error(c, 400, "invalid_playlist", "Playlist identifier is invalid")
		return
	}
	var item models.PlaylistItem
	if c.ShouldBindJSON(&item) != nil || !validMedia(item.MediaType, item.MediaID) {
		apiresponse.Error(c, 400, "invalid_media", "A valid media item is required")
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	if cc.collection("playlists").FindOne(ctx, bson.M{"_id": pid, "owner_user_id": uid}).Err() != nil {
		apiresponse.Error(c, 404, "playlist_not_found", "Playlist not found or not yours")
		return
	}
	item.ID = bson.NewObjectID()
	item.PlaylistID = pid
	item.Title = cleanText(item.Title, 200)
	item.CreatedAt = time.Now().UTC()
	_, err = cc.collection("playlist_items").InsertOne(ctx, item)
	if mongo.IsDuplicateKeyError(err) {
		apiresponse.Error(c, 409, "duplicate_item", "That title is already in this playlist")
		return
	}
	c.JSON(201, gin.H{"data": item})
}

func (cc *CommunityController) PutProgress(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	var show, season, episode int
	if _, e := fmtSscanf(c.Param("showId"), &show); e != nil {
		return
	}
	fmtSscanf(c.Param("season"), &season)
	fmtSscanf(c.Param("episode"), &episode)
	var p models.EpisodeProgress
	if c.ShouldBindJSON(&p) != nil {
		apiresponse.Error(c, 400, "invalid_request", "Progress payload is invalid")
		return
	}
	if p.PositionSeconds < 0 || p.DurationSeconds < 0 || (p.DurationSeconds > 0 && p.PositionSeconds > p.DurationSeconds) {
		apiresponse.Error(c, 400, "invalid_progress", "Position cannot exceed duration")
		return
	}
	p.UserID = uid
	p.ShowID = show
	p.SeasonNumber = season
	p.EpisodeNumber = episode
	p.UpdatedAt = time.Now().UTC()
	if p.Status == "WATCHED" {
		p.CompletedAt = &p.UpdatedAt
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	_, err := cc.collection("episode_progress").UpdateOne(ctx, bson.M{"user_id": uid, "show_id": show, "season_number": season, "episode_number": episode}, bson.M{"$set": p}, options.UpdateOne().SetUpsert(true))
	if err != nil {
		apiresponse.Error(c, 500, "database_error", "Unable to update progress")
		return
	}
	c.JSON(200, gin.H{"data": p})
}
func (cc *CommunityController) GetProgress(c *gin.Context) {
	uid, _ := utils.GetUserIdFromContext(c)
	var show int
	if _, e := fmtSscanf(c.Param("showId"), &show); e != nil {
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	cur, _ := cc.collection("episode_progress").Find(ctx, bson.M{"user_id": uid, "show_id": show}, options.Find().SetSort(bson.D{{Key: "season_number", Value: 1}, {Key: "episode_number", Value: 1}}))
	rows := []models.EpisodeProgress{}
	if cur != nil {
		defer cur.Close(ctx)
		cur.All(ctx, &rows)
	}
	c.JSON(200, gin.H{"data": rows})
}

func (cc *CommunityController) ModerationReports(c *gin.Context) {
	ctx, cancel := requestContext(c)
	defer cancel()
	cur, _ := cc.collection("review_reports").Find(ctx, bson.M{"status": "OPEN"}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	rows := []models.ReviewReport{}
	if cur != nil {
		defer cur.Close(ctx)
		cur.All(ctx, &rows)
	}
	c.JSON(200, gin.H{"data": rows})
}
func (cc *CommunityController) ModerateReview(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("reviewId"))
	if err != nil {
		return
	}
	var b struct {
		Status string `json:"status"`
	}
	c.ShouldBindJSON(&b)
	if b.Status != "VISIBLE" && b.Status != "HIDDEN" && b.Status != "REMOVED" {
		apiresponse.Error(c, 400, "invalid_status", "Invalid moderation status")
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	cc.collection("reviews").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": b.Status, "updated_at": time.Now().UTC()}})
	c.Status(204)
}
func (cc *CommunityController) SuspendUser(c *gin.Context) {
	id := c.Param("userId")
	var b struct {
		Suspended bool `json:"suspended"`
	}
	if c.ShouldBindJSON(&b) != nil {
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	res, _ := cc.collection("users").UpdateOne(ctx, bson.M{"user_id": id}, bson.M{"$set": bson.M{"suspended": b.Suspended, "updated_at": time.Now().UTC()}})
	if res.MatchedCount == 0 {
		apiresponse.Error(c, 404, "user_not_found", "User not found")
		return
	}
	c.Status(204)
}
