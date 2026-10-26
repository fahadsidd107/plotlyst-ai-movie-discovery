package main

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"log"
	"plotlyst/server/config"
	"plotlyst/server/database"
	"plotlyst/server/migration"
	"plotlyst/server/models"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := database.Connect(ctx, cfg.MongoDBURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	collection := database.OpenCollection("movies", client)
	cursor, err := collection.Find(ctx, bson.M{"editorial_rating.schema_version": bson.M{"$ne": 1}})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	var changed, skipped int
	for cursor.Next(ctx) {
		var movie models.Movie
		if err := cursor.Decode(&movie); err != nil {
			skipped++
			continue
		}
		stars, ok := migration.LegacyRankingToStars(movie.Ranking.RankingValue, movie.Ranking.RankingName)
		if !ok {
			skipped++
			continue
		}
		_, err = collection.UpdateOne(ctx, bson.M{"_id": movie.ID, "editorial_rating.schema_version": bson.M{"$ne": 1}}, bson.M{"$set": bson.M{"editorial_rating": models.EditorialRating{Stars: stars, Scale: 5, Source: "legacy_editorial_ranking", SchemaVersion: 1}}})
		if err != nil {
			skipped++
			continue
		}
		changed++
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("rating migration complete: updated=%d skipped=%d", changed, skipped)
}
