package database

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestEnsureIndexesAndUniqueRepositoryConstraints(t *testing.T) {
	uri := os.Getenv("TEST_MONGODB_URI")
	if uri == "" {
		t.Skip("TEST_MONGODB_URI not configured")
	}
	os.Setenv("DATABASE_NAME", "plotlyst_integration")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := Connect(ctx, uri)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(ctx)
	db := client.Database("plotlyst_integration")
	defer db.Drop(ctx)
	if err := EnsureIndexes(ctx, client); err != nil {
		t.Fatal(err)
	}
	col := db.Collection("ratings")
	doc := bson.M{"user_id": "u1", "media_type": "movie", "media_id": 42, "stars": 5}
	if _, err = col.InsertOne(ctx, doc); err != nil {
		t.Fatal(err)
	}
	if _, err = col.InsertOne(ctx, doc); err == nil {
		t.Fatal("compound ownership index accepted duplicate rating")
	}
}
