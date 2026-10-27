package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func EnsureIndexes(ctx context.Context, client *mongo.Client) error {
	db := client.Database(os.Getenv("DATABASE_NAME"))
	indexes := map[string][]mongo.IndexModel{
		"users":            {{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true).SetName("unique_email")}},
		"profiles":         {{Keys: bson.D{{Key: "user_id", Value: 1}}, Options: options.Index().SetUnique(true)}, {Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)}},
		"ratings":          {{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "media_type", Value: 1}, {Key: "media_id", Value: 1}}, Options: options.Index().SetUnique(true)}, {Keys: bson.D{{Key: "media_type", Value: 1}, {Key: "media_id", Value: 1}}}},
		"reviews":          {{Keys: bson.D{{Key: "media_type", Value: 1}, {Key: "media_id", Value: 1}, {Key: "status", Value: 1}, {Key: "created_at", Value: -1}}}, {Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: -1}}}},
		"review_votes":     {{Keys: bson.D{{Key: "review_id", Value: 1}, {Key: "user_id", Value: 1}}, Options: options.Index().SetUnique(true)}},
		"review_reports":   {{Keys: bson.D{{Key: "review_id", Value: 1}, {Key: "reporter_user_id", Value: 1}}, Options: options.Index().SetUnique(true)}},
		"watchlist_items":  {{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "media_type", Value: 1}, {Key: "media_id", Value: 1}}, Options: options.Index().SetUnique(true)}},
		"watch_history":    {{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "media_type", Value: 1}, {Key: "media_id", Value: 1}}, Options: options.Index().SetUnique(true)}},
		"playlists":        {{Keys: bson.D{{Key: "owner_user_id", Value: 1}, {Key: "updated_at", Value: -1}}}},
		"playlist_items":   {{Keys: bson.D{{Key: "playlist_id", Value: 1}, {Key: "media_type", Value: 1}, {Key: "media_id", Value: 1}}, Options: options.Index().SetUnique(true)}},
		"episode_progress": {{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "show_id", Value: 1}, {Key: "season_number", Value: 1}, {Key: "episode_number", Value: 1}}, Options: options.Index().SetUnique(true)}},
	}
	for name, models := range indexes {
		if _, err := db.Collection(name).Indexes().CreateMany(ctx, models); err != nil {
			return err
		}
	}
	return nil
}
