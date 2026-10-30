package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	RoleUser      = "USER"
	RoleModerator = "MODERATOR"
	RoleAdmin     = "ADMIN"
)

type PrivacySettings struct {
	Profile  string `json:"profile" bson:"profile"`
	Library  string `json:"library" bson:"library"`
	Activity string `json:"activity" bson:"activity"`
}

type Profile struct {
	ID              bson.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID          string          `json:"user_id" bson:"user_id"`
	Username        string          `json:"username" bson:"username"`
	DisplayName     string          `json:"display_name" bson:"display_name"`
	Bio             string          `json:"bio" bson:"bio"`
	AvatarURL       string          `json:"avatar_url" bson:"avatar_url"`
	Locale          string          `json:"locale" bson:"locale"`
	Region          string          `json:"region" bson:"region"`
	Timezone        string          `json:"timezone" bson:"timezone"`
	ShowSpoilers    bool            `json:"show_spoilers" bson:"show_spoilers"`
	Privacy         PrivacySettings `json:"privacy" bson:"privacy"`
	FavouriteGenres []Genre         `json:"favourite_genres" bson:"favourite_genres"`
	CreatedAt       time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" bson:"updated_at"`
}

type MediaRef struct {
	MediaType string `json:"media_type" bson:"media_type"`
	MediaID   int    `json:"media_id" bson:"media_id"`
}

type Rating struct {
	ID        bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    string        `json:"user_id" bson:"user_id"`
	MediaRef  `bson:",inline"`
	Stars     int       `json:"stars" bson:"stars"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type Review struct {
	ID                bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID            string        `json:"user_id" bson:"user_id"`
	AuthorName        string        `json:"author_name" bson:"author_name"`
	MediaRef          `bson:",inline"`
	ParentID          *bson.ObjectID `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	Title             string         `json:"title" bson:"title"`
	Body              string         `json:"body" bson:"body"`
	Spoiler           bool           `json:"spoiler" bson:"spoiler"`
	SpoilerSource     string         `json:"spoiler_source,omitempty" bson:"spoiler_source,omitempty"`
	SpoilerConfidence float64        `json:"spoiler_confidence,omitempty" bson:"spoiler_confidence,omitempty"`
	SpoilerHidden     bool           `json:"spoiler_hidden,omitempty" bson:"-"`
	SpoilerReason     string         `json:"spoiler_reason,omitempty" bson:"-"`
	Status            string         `json:"status" bson:"status"`
	HelpfulCount      int64          `json:"helpful_count" bson:"helpful_count"`
	CreatedAt         time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at" bson:"updated_at"`
	DeletedAt         *time.Time     `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type LibraryItem struct {
	ID         bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID     string        `json:"user_id" bson:"user_id"`
	MediaRef   `bson:",inline"`
	Title      string     `json:"title,omitempty" bson:"title,omitempty"`
	PosterPath string     `json:"poster_path,omitempty" bson:"poster_path,omitempty"`
	Status     string     `json:"status,omitempty" bson:"status,omitempty"`
	WatchedAt  *time.Time `json:"watched_at,omitempty" bson:"watched_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" bson:"updated_at"`
}

type Playlist struct {
	ID          bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OwnerUserID string        `json:"owner_user_id" bson:"owner_user_id"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	Private     bool          `json:"private" bson:"private"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
}

type PlaylistItem struct {
	ID         bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PlaylistID bson.ObjectID `json:"playlist_id" bson:"playlist_id"`
	MediaRef   `bson:",inline"`
	Title      string    `json:"title,omitempty" bson:"title,omitempty"`
	PosterPath string    `json:"poster_path,omitempty" bson:"poster_path,omitempty"`
	Position   int       `json:"position" bson:"position"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
}

type EpisodeProgress struct {
	ID              bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID          string        `json:"user_id" bson:"user_id"`
	ShowID          int           `json:"show_id" bson:"show_id"`
	SeasonNumber    int           `json:"season_number" bson:"season_number"`
	EpisodeNumber   int           `json:"episode_number" bson:"episode_number"`
	Status          string        `json:"status" bson:"status"`
	PositionSeconds int           `json:"position_seconds" bson:"position_seconds"`
	DurationSeconds int           `json:"duration_seconds" bson:"duration_seconds"`
	CompletedAt     *time.Time    `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	UpdatedAt       time.Time     `json:"updated_at" bson:"updated_at"`
}

type ReviewVote struct {
	ReviewID  bson.ObjectID `bson:"review_id"`
	UserID    string        `bson:"user_id"`
	CreatedAt time.Time     `bson:"created_at"`
}
type ReviewReport struct {
	ID             bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ReviewID       bson.ObjectID `json:"review_id" bson:"review_id"`
	ReporterUserID string        `json:"reporter_user_id" bson:"reporter_user_id"`
	Reason         string        `json:"reason" bson:"reason"`
	Status         string        `json:"status" bson:"status"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
}
