package spoiler

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	ModeOff      = "OFF"
	ModeStandard = "STANDARD"
	ModeStrict   = "STRICT"
)

type Preference struct {
	UserID    string    `json:"-" bson:"user_id"`
	Mode      string    `json:"mode" bson:"mode"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type Service struct {
	Preferences *mongo.Collection
	Progress    *mongo.Collection
}

func Normalize(mode string) string {
	switch strings.ToUpper(strings.TrimSpace(mode)) {
	case ModeOff:
		return ModeOff
	case ModeStrict:
		return ModeStrict
	default:
		return ModeStandard
	}
}

func (s *Service) Get(ctx context.Context, userID string) Preference {
	preference := Preference{UserID: userID, Mode: ModeStandard}
	if userID != "" {
		_ = s.Preferences.FindOne(ctx, bson.M{"user_id": userID}).Decode(&preference)
	}
	preference.Mode = Normalize(preference.Mode)
	return preference
}

func (s *Service) Put(ctx context.Context, userID, mode string) (Preference, error) {
	preference := Preference{UserID: userID, Mode: Normalize(mode), UpdatedAt: time.Now().UTC()}
	_, err := s.Preferences.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": preference}, options.UpdateOne().SetUpsert(true))
	return preference, err
}

func ConcealReview(mode string, tagged bool) bool {
	mode = Normalize(mode)
	return mode == ModeStrict || (mode == ModeStandard && tagged)
}

// EpisodeVisible permits descriptions only through the user's furthest watched
// episode. OFF disables the firewall; an explicit reveal is handled by callers.
func (s *Service) EpisodeVisible(ctx context.Context, userID, mode string, showID, season, episode int) bool {
	if Normalize(mode) == ModeOff {
		return true
	}
	if userID == "" {
		return false
	}
	var latest struct {
		Season  int `bson:"season_number"`
		Episode int `bson:"episode_number"`
	}
	err := s.Progress.FindOne(ctx, bson.M{"user_id": userID, "show_id": showID, "status": bson.M{"$in": bson.A{"WATCHED", "COMPLETED"}}}, options.FindOne().SetSort(bson.D{{Key: "season_number", Value: -1}, {Key: "episode_number", Value: -1}})).Decode(&latest)
	if err != nil {
		return false
	}
	return season < latest.Season || (season == latest.Season && episode <= latest.Episode)
}
