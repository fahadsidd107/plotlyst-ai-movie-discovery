package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Genre struct {
	GenreID   int    `bson:"genre_id" json:"genre_id" validate:"required"`
	GenreName string `bson:"genre_name" json:"genre_name" validate:"required,min=2,max=100"`
}

type Ranking struct {
	RankingValue int    `bson:"ranking_value" json:"ranking_value" validate:"required"`
	RankingName  string `bson:"ranking_name" json:"ranking_name" validate:"required"`
}

type Movie struct {
	ID              bson.ObjectID    `bson:"_id,omitempty" json:"_id,omitempty"`
	ImdbID          string           `bson:"imdb_id" json:"imdb_id" validate:"required"`
	Title           string           `bson:"title" json:"title" validate:"required,min=2,max=500"`
	PosterPath      string           `bson:"poster_path" json:"poster_path" validate:"required,url"`
	Genre           []Genre          `bson:"genre" json:"genre" validate:"required,dive"`
	AdminReview     string           `bson:"admin_review" json:"admin_review"`
	Ranking         Ranking          `bson:"ranking" json:"ranking" validate:"required"`
	EditorialRating *EditorialRating `bson:"editorial_rating,omitempty" json:"editorial_rating,omitempty"`
}

type EditorialRating struct {
	Stars         float64 `bson:"stars" json:"stars"`
	Scale         float64 `bson:"scale" json:"scale"`
	Source        string  `bson:"source" json:"source"`
	SchemaVersion int     `bson:"schema_version" json:"schema_version"`
}
