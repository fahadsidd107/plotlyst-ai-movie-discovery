package storydna

import "time"

type Scores struct {
	EmotionalWarmth int `json:"emotional_warmth" bson:"emotional_warmth"`
	Tension         int `json:"tension" bson:"tension"`
	Darkness        int `json:"darkness" bson:"darkness"`
	Hope            int `json:"hope" bson:"hope"`
	Pacing          int `json:"pacing" bson:"pacing"`
	Complexity      int `json:"complexity" bson:"complexity"`
	VisualStyle     int `json:"visual_style" bson:"visual_style"`
	Dialogue        int `json:"dialogue" bson:"dialogue"`
}

type SpoilerTraits struct {
	EndingShape string `json:"ending_shape,omitempty" bson:"ending_shape,omitempty"`
	Resolution  string `json:"resolution,omitempty" bson:"resolution,omitempty"`
}

type Provenance struct {
	Method        string    `json:"method" bson:"method"`
	Source        string    `json:"source" bson:"source"`
	Policy        string    `json:"policy" bson:"policy"`
	Confidence    float64   `json:"confidence" bson:"confidence"`
	GeneratedAt   time.Time `json:"generated_at" bson:"generated_at"`
	SchemaVersion int       `json:"schema_version" bson:"schema_version"`
}

type Profile struct {
	MovieID     int           `json:"movie_id" bson:"movie_id"`
	Title       string        `json:"title" bson:"title"`
	PosterPath  string        `json:"poster_path,omitempty" bson:"poster_path,omitempty"`
	ReleaseYear int           `json:"release_year,omitempty" bson:"release_year,omitempty"`
	Structures  []string      `json:"structures" bson:"structures"`
	Styles      []string      `json:"styles" bson:"styles"`
	Tones       []string      `json:"tones" bson:"tones"`
	Scores      Scores        `json:"scores" bson:"scores"`
	Spoiler     SpoilerTraits `json:"spoiler_traits,omitempty" bson:"spoiler_traits,omitempty"`
	Provenance  Provenance    `json:"provenance" bson:"provenance"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
}

type Weights struct {
	Mood       float64 `json:"mood"`
	Pacing     float64 `json:"pacing"`
	Complexity float64 `json:"complexity"`
	Structure  float64 `json:"structure"`
}

type Match struct {
	MovieID     int      `json:"movie_id"`
	Title       string   `json:"title"`
	PosterPath  string   `json:"poster_path,omitempty"`
	ReleaseYear int      `json:"release_year,omitempty"`
	Similarity  float64  `json:"similarity"`
	Reasons     []string `json:"reasons"`
	Profile     Profile  `json:"profile"`
}

type DiscoveryRequest struct {
	Query           string  `json:"query"`
	SeedMovieID     int     `json:"seed_movie_id"`
	Weights         Weights `json:"weights"`
	IncludeSpoilers bool    `json:"include_spoilers"`
}
