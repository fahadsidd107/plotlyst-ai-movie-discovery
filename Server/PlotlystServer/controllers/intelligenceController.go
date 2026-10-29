package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/net/websocket"
	"plotlyst/server/apiresponse"
	"plotlyst/server/database"
	"plotlyst/server/intelligence"
)

type IntelligenceController struct {
	client  *mongo.Client
	service *intelligence.Service
}
type aiPrivacy struct {
	UserID                       string    `json:"-" bson:"user_id"`
	AIEnabled                    bool      `json:"ai_enabled" bson:"ai_enabled"`
	PersonalizationEnabled       bool      `json:"personalization_enabled" bson:"personalization_enabled"`
	AllowActivityForExplanations bool      `json:"allow_activity_for_explanations" bson:"allow_activity_for_explanations"`
	UpdatedAt                    time.Time `json:"updated_at" bson:"updated_at"`
}

func NewIntelligenceController(client *mongo.Client, service *intelligence.Service) *IntelligenceController {
	return &IntelligenceController{client: client, service: service}
}
func (h *IntelligenceController) col(name string) *mongo.Collection {
	return database.OpenCollection(name, h.client)
}
func (h *IntelligenceController) Discover(mode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var b struct {
			Query   string `json:"query"`
			Memory  string `json:"memory"`
			Region  string `json:"region"`
			Consent bool   `json:"consent_to_personalization"`
		}
		if c.ShouldBindJSON(&b) != nil {
			apiresponse.Error(c, 400, "invalid_request", "A JSON discovery request is required")
			return
		}
		q := b.Query
		if mode == "memory" && q == "" {
			q = b.Memory
		}
		out, err := h.service.Discover(c.Request.Context(), q, mode, b.Region)
		if err != nil {
			apiresponse.Error(c, 400, "discovery_failed", err.Error())
			return
		}
		if b.Consent {
			var privacy aiPrivacy
			if h.col("ai_privacy").FindOne(c.Request.Context(), bson.M{"user_id": c.GetString("userId")}).Decode(&privacy) == nil && privacy.PersonalizationEnabled && privacy.AllowActivityForExplanations {
				five, _ := h.col("ratings").CountDocuments(c.Request.Context(), bson.M{"user_id": c.GetString("userId"), "stars": 5})
				watched, _ := h.col("watch_history").CountDocuments(c.Request.Context(), bson.M{"user_id": c.GetString("userId")})
				for n := range out.Items {
					out.Items[n].Reasons = append(out.Items[n].Reasons, "Personalized with your Plotlyst activity: "+strconv.FormatInt(five, 10)+" five-star ratings and "+strconv.FormatInt(watched, 10)+" watched titles.")
				}
				out.Provenance.Explanation = "grounded-activity-template-v1"
			}
		}
		c.JSON(200, gin.H{"data": out})
	}
}
func (h *IntelligenceController) GetAIPrivacy(c *gin.Context) {
	var p aiPrivacy
	if h.col("ai_privacy").FindOne(c.Request.Context(), bson.M{"user_id": c.GetString("userId")}).Decode(&p) != nil {
		p = aiPrivacy{UserID: c.GetString("userId")}
	}
	c.JSON(200, gin.H{"data": p})
}
func (h *IntelligenceController) PutAIPrivacy(c *gin.Context) {
	var p aiPrivacy
	if c.ShouldBindJSON(&p) != nil {
		apiresponse.Error(c, 400, "invalid_privacy", "Valid AI privacy settings are required")
		return
	}
	p.UserID = c.GetString("userId")
	p.UpdatedAt = time.Now().UTC()
	_, err := h.col("ai_privacy").UpdateOne(c.Request.Context(), bson.M{"user_id": p.UserID}, bson.M{"$set": p}, options.UpdateOne().SetUpsert(true))
	if err != nil {
		apiresponse.Error(c, 500, "privacy_update_failed", "Unable to update AI privacy settings")
		return
	}
	c.JSON(200, gin.H{"data": p})
}

func (h *IntelligenceController) CinemaDNA(c *gin.Context) {
	uid := c.GetString("userId")
	ctx, cancel := requestContext(c)
	defer cancel()
	type starRow struct {
		ID    int   `bson:"_id" json:"stars"`
		Count int64 `bson:"count" json:"count"`
	}
	stars := []starRow{}
	cur, _ := h.col("ratings").Aggregate(ctx, mongo.Pipeline{{{Key: "$match", Value: bson.M{"user_id": uid}}}, {{Key: "$group", Value: bson.M{"_id": "$stars", "count": bson.M{"$sum": 1}}}}, {{Key: "$sort", Value: bson.M{"_id": 1}}}})
	if cur != nil {
		defer cur.Close(ctx)
		_ = cur.All(ctx, &stars)
	}
	ratings, _ := h.col("ratings").CountDocuments(ctx, bson.M{"user_id": uid})
	history, _ := h.col("watch_history").CountDocuments(ctx, bson.M{"user_id": uid})
	watchlist, _ := h.col("watchlist_items").CountDocuments(ctx, bson.M{"user_id": uid})
	reviews, _ := h.col("reviews").CountDocuments(ctx, bson.M{"user_id": uid, "deleted_at": nil})
	completed, _ := h.col("episode_progress").CountDocuments(ctx, bson.M{"user_id": uid, "status": "COMPLETED"})
	episodes, _ := h.col("episode_progress").CountDocuments(ctx, bson.M{"user_id": uid})
	c.JSON(200, gin.H{"data": gin.H{"totals": gin.H{"ratings": ratings, "watched": history, "watchlist": watchlist, "reviews": reviews, "episodes": episodes, "episodes_completed": completed}, "rating_distribution": stars, "coverage": gin.H{"activity_records": ratings + history + episodes, "enriched_records": 0, "note": "Genre, runtime and language DNA will appear as activity metadata is enriched."}}, "source": gin.H{"provider": "Plotlyst MongoDB aggregations", "status": "available", "generated_at": time.Now().UTC()}})
}

type nightPreferences struct {
	Genres     []string `json:"genres" bson:"genres"`
	Moods      []string `json:"moods" bson:"moods"`
	Languages  []string `json:"languages" bson:"languages"`
	RuntimeMax int      `json:"runtime_max" bson:"runtime_max"`
}
type nightParticipant struct {
	UserID      string           `json:"user_id" bson:"user_id"`
	Preferences nightPreferences `json:"preferences" bson:"preferences"`
	JoinedAt    time.Time        `json:"joined_at" bson:"joined_at"`
}
type movieNight struct {
	ID           bson.ObjectID      `json:"id" bson:"_id,omitempty"`
	Code         string             `json:"code" bson:"code"`
	OwnerID      string             `json:"owner_id" bson:"owner_id"`
	Status       string             `json:"status" bson:"status"`
	Participants []nightParticipant `json:"participants" bson:"participants"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	ExpiresAt    time.Time          `json:"expires_at" bson:"expires_at"`
}

func roomCode() string {
	b := make([]byte, 3)
	_, _ = rand.Read(b)
	return strings.ToUpper(hex.EncodeToString(b))
}
func (h *IntelligenceController) CreateNight(c *gin.Context) {
	now := time.Now().UTC()
	r := movieNight{Code: roomCode(), OwnerID: c.GetString("userId"), Status: "LOBBY", Participants: []nightParticipant{{UserID: c.GetString("userId"), JoinedAt: now}}, CreatedAt: now, ExpiresAt: now.Add(24 * time.Hour)}
	res, err := h.col("movie_nights").InsertOne(c.Request.Context(), r)
	if err != nil {
		apiresponse.Error(c, 500, "room_create_failed", "Unable to create movie night")
		return
	}
	r.ID = res.InsertedID.(bson.ObjectID)
	c.JSON(201, gin.H{"data": r})
}
func (h *IntelligenceController) JoinNight(c *gin.Context) {
	uid := c.GetString("userId")
	code := strings.ToUpper(c.Param("code"))
	ctx, cancel := requestContext(c)
	defer cancel()
	h.col("movie_nights").UpdateOne(ctx, bson.M{"code": code, "participants.user_id": bson.M{"$ne": uid}}, bson.M{"$push": bson.M{"participants": nightParticipant{UserID: uid, JoinedAt: time.Now().UTC()}}})
	h.getNight(c, code)
}
func (h *IntelligenceController) GetNight(c *gin.Context) {
	h.getNight(c, strings.ToUpper(c.Param("code")))
}
func (h *IntelligenceController) getNight(c *gin.Context, code string) {
	var r movieNight
	err := h.col("movie_nights").FindOne(c.Request.Context(), bson.M{"code": code, "participants.user_id": c.GetString("userId")}).Decode(&r)
	if err != nil {
		apiresponse.Error(c, 404, "room_not_found", "Movie night was not found or you are not a participant")
		return
	}
	c.JSON(200, gin.H{"data": r, "intersection": intersect(r.Participants)})
}
func (h *IntelligenceController) Preferences(c *gin.Context) {
	var p nightPreferences
	if c.ShouldBindJSON(&p) != nil {
		apiresponse.Error(c, 400, "invalid_preferences", "Valid preferences are required")
		return
	}
	code := strings.ToUpper(c.Param("code"))
	res, _ := h.col("movie_nights").UpdateOne(c.Request.Context(), bson.M{"code": code, "participants.user_id": c.GetString("userId")}, bson.M{"$set": bson.M{"participants.$.preferences": p}})
	if res == nil || res.MatchedCount == 0 {
		apiresponse.Error(c, 404, "room_not_found", "Movie night was not found")
		return
	}
	h.getNight(c, code)
}
func intersect(ps []nightParticipant) gin.H {
	counts := map[string]int{}
	languageCounts := map[string]int{}
	maxRuntime := 0
	for _, p := range ps {
		for _, g := range p.Preferences.Genres {
			counts[strings.ToLower(g)]++
		}
		for _, language := range p.Preferences.Languages {
			languageCounts[strings.ToLower(language)]++
		}
		if p.Preferences.RuntimeMax > 0 && (maxRuntime == 0 || p.Preferences.RuntimeMax < maxRuntime) {
			maxRuntime = p.Preferences.RuntimeMax
		}
	}
	genres := []string{}
	languages := []string{}
	for g, n := range counts {
		if n == len(ps) {
			genres = append(genres, g)
		}
	}
	for language, n := range languageCounts {
		if n == len(ps) {
			languages = append(languages, language)
		}
	}
	sort.Strings(genres)
	sort.Strings(languages)
	return gin.H{"genres": genres, "languages": languages, "runtime_max": maxRuntime, "participants": len(ps), "method": "hard preference intersection; no AI inference"}
}
func (h *IntelligenceController) NightSocket(c *gin.Context) {
	code := strings.ToUpper(c.Param("code"))
	uid := c.GetString("userId")
	if h.col("movie_nights").FindOne(c.Request.Context(), bson.M{"code": code, "participants.user_id": uid}).Err() != nil {
		apiresponse.Error(c, 403, "room_membership_required", "Join the room before opening its live channel")
		return
	}
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		ws.SetDeadline(time.Now().Add(35 * time.Second))
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			var r movieNight
			if h.col("movie_nights").FindOne(c.Request.Context(), bson.M{"code": code}).Decode(&r) != nil {
				return
			}
			if websocket.JSON.Send(ws, gin.H{"type": "room_snapshot", "data": r, "intersection": intersect(r.Participants)}) != nil {
				return
			}
		}
	}).ServeHTTP(c.Writer, c.Request)
}
