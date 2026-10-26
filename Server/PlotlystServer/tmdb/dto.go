package tmdb

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Company struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	LogoPath      string `json:"logo_path"`
	OriginCountry string `json:"origin_country"`
}
type Collection struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	PosterPath   string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`
}
type Provider struct {
	ProviderID      int    `json:"provider_id"`
	ProviderName    string `json:"provider_name"`
	LogoPath        string `json:"logo_path"`
	DisplayPriority int    `json:"display_priority"`
}
type ProviderRegion struct {
	Link     string     `json:"link"`
	Flatrate []Provider `json:"flatrate"`
	Rent     []Provider `json:"rent"`
	Buy      []Provider `json:"buy"`
	Free     []Provider `json:"free"`
}
type WatchProviders struct {
	Results map[string]ProviderRegion `json:"results"`
}

type Credit struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	Character    string `json:"character,omitempty"`
	Job          string `json:"job,omitempty"`
	Department   string `json:"department,omitempty"`
	ProfilePath  string `json:"profile_path"`
	Order        int    `json:"order,omitempty"`
}
type Credits struct {
	Cast []Credit `json:"cast"`
	Crew []Credit `json:"crew"`
}

type MediaSummary struct {
	ID                 int     `json:"id"`
	MediaType          string  `json:"media_type,omitempty"`
	Title              string  `json:"title,omitempty"`
	Name               string  `json:"name,omitempty"`
	Overview           string  `json:"overview"`
	PosterPath         string  `json:"poster_path"`
	BackdropPath       string  `json:"backdrop_path"`
	ProfilePath        string  `json:"profile_path,omitempty"`
	KnownForDepartment string  `json:"known_for_department,omitempty"`
	ReleaseDate        string  `json:"release_date,omitempty"`
	FirstAirDate       string  `json:"first_air_date,omitempty"`
	GenreIDs           []int   `json:"genre_ids"`
	VoteAverage        float64 `json:"vote_average"`
	VoteCount          int     `json:"vote_count"`
	Popularity         float64 `json:"popularity"`
}
type MovieSummary struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title"`
	Overview      string  `json:"overview"`
	PosterPath    string  `json:"poster_path"`
	BackdropPath  string  `json:"backdrop_path"`
	ReleaseDate   string  `json:"release_date"`
	GenreIDs      []int   `json:"genre_ids"`
	VoteAverage   float64 `json:"vote_average"`
	VoteCount     int     `json:"vote_count"`
	Popularity    float64 `json:"popularity"`
}
type TVSummary struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	OriginalName  string   `json:"original_name"`
	Overview      string   `json:"overview"`
	PosterPath    string   `json:"poster_path"`
	BackdropPath  string   `json:"backdrop_path"`
	FirstAirDate  string   `json:"first_air_date"`
	OriginCountry []string `json:"origin_country"`
	GenreIDs      []int    `json:"genre_ids"`
	VoteAverage   float64  `json:"vote_average"`
	VoteCount     int      `json:"vote_count"`
	Popularity    float64  `json:"popularity"`
}
type Page[T any] struct {
	Page         int `json:"page"`
	Results      []T `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type ReleaseDate struct {
	Certification string `json:"certification"`
	ISO31661      string `json:"iso_3166_1"`
	ReleaseDate   string `json:"release_date"`
	Type          int    `json:"type"`
}
type ReleaseCountry struct {
	ISO31661     string        `json:"iso_3166_1"`
	ReleaseDates []ReleaseDate `json:"release_dates"`
}
type ReleaseDates struct {
	Results []ReleaseCountry `json:"results"`
}
type ContentRating struct {
	Rating   string `json:"rating"`
	ISO31661 string `json:"iso_3166_1"`
}
type ContentRatings struct {
	Results []ContentRating `json:"results"`
}

type MovieDetail struct {
	ID                  int                `json:"id"`
	IMDbID              string             `json:"imdb_id"`
	Title               string             `json:"title"`
	OriginalTitle       string             `json:"original_title"`
	Tagline             string             `json:"tagline"`
	Overview            string             `json:"overview"`
	Status              string             `json:"status"`
	Runtime             int                `json:"runtime"`
	ReleaseDate         string             `json:"release_date"`
	PosterPath          string             `json:"poster_path"`
	BackdropPath        string             `json:"backdrop_path"`
	VoteAverage         float64            `json:"vote_average"`
	VoteCount           int                `json:"vote_count"`
	Genres              []Genre            `json:"genres"`
	ProductionCompanies []Company          `json:"production_companies"`
	BelongsToCollection *Collection        `json:"belongs_to_collection"`
	Credits             Credits            `json:"credits"`
	ReleaseDates        ReleaseDates       `json:"release_dates"`
	Recommendations     Page[MediaSummary] `json:"recommendations"`
	WatchProviders      WatchProviders     `json:"watch/providers"`
}

type SeasonSummary struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Overview     string  `json:"overview"`
	AirDate      string  `json:"air_date"`
	EpisodeCount int     `json:"episode_count"`
	PosterPath   string  `json:"poster_path"`
	SeasonNumber int     `json:"season_number"`
	VoteAverage  float64 `json:"vote_average"`
}
type TVDetail struct {
	ID                  int                `json:"id"`
	Name                string             `json:"name"`
	OriginalName        string             `json:"original_name"`
	Tagline             string             `json:"tagline"`
	Overview            string             `json:"overview"`
	Status              string             `json:"status"`
	Type                string             `json:"type"`
	FirstAirDate        string             `json:"first_air_date"`
	LastAirDate         string             `json:"last_air_date"`
	NumberOfSeasons     int                `json:"number_of_seasons"`
	NumberOfEpisodes    int                `json:"number_of_episodes"`
	EpisodeRunTime      []int              `json:"episode_run_time"`
	PosterPath          string             `json:"poster_path"`
	BackdropPath        string             `json:"backdrop_path"`
	VoteAverage         float64            `json:"vote_average"`
	VoteCount           int                `json:"vote_count"`
	Genres              []Genre            `json:"genres"`
	Networks            []Company          `json:"networks"`
	ProductionCompanies []Company          `json:"production_companies"`
	Seasons             []SeasonSummary    `json:"seasons"`
	Credits             Credits            `json:"credits"`
	Recommendations     Page[MediaSummary] `json:"recommendations"`
	WatchProviders      WatchProviders     `json:"watch/providers"`
	ContentRatings      ContentRatings     `json:"content_ratings"`
}

type EpisodeSummary struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Overview      string   `json:"overview"`
	AirDate       string   `json:"air_date"`
	EpisodeNumber int      `json:"episode_number"`
	SeasonNumber  int      `json:"season_number"`
	Runtime       int      `json:"runtime"`
	StillPath     string   `json:"still_path"`
	VoteAverage   float64  `json:"vote_average"`
	Crew          []Credit `json:"crew"`
	GuestStars    []Credit `json:"guest_stars"`
}
type SeasonDetail struct {
	ID           int              `json:"id"`
	Name         string           `json:"name"`
	Overview     string           `json:"overview"`
	AirDate      string           `json:"air_date"`
	PosterPath   string           `json:"poster_path"`
	SeasonNumber int              `json:"season_number"`
	Episodes     []EpisodeSummary `json:"episodes"`
	Credits      Credits          `json:"credits"`
}
type EpisodeDetail struct {
	EpisodeSummary
	Credits Credits `json:"credits"`
}

type PersonSummary struct {
	ID                 int            `json:"id"`
	Name               string         `json:"name"`
	KnownForDepartment string         `json:"known_for_department"`
	ProfilePath        string         `json:"profile_path"`
	Popularity         float64        `json:"popularity"`
	KnownFor           []MediaSummary `json:"known_for"`
}
type CombinedCredit struct {
	ID           int     `json:"id"`
	MediaType    string  `json:"media_type"`
	Title        string  `json:"title,omitempty"`
	Name         string  `json:"name,omitempty"`
	Character    string  `json:"character,omitempty"`
	Job          string  `json:"job,omitempty"`
	PosterPath   string  `json:"poster_path"`
	ReleaseDate  string  `json:"release_date,omitempty"`
	FirstAirDate string  `json:"first_air_date,omitempty"`
	VoteAverage  float64 `json:"vote_average"`
}
type CombinedCredits struct {
	Cast []CombinedCredit `json:"cast"`
	Crew []CombinedCredit `json:"crew"`
}
type PersonDetail struct {
	ID                 int             `json:"id"`
	IMDbID             string          `json:"imdb_id"`
	Name               string          `json:"name"`
	Biography          string          `json:"biography"`
	Birthday           string          `json:"birthday"`
	Deathday           string          `json:"deathday"`
	PlaceOfBirth       string          `json:"place_of_birth"`
	KnownForDepartment string          `json:"known_for_department"`
	ProfilePath        string          `json:"profile_path"`
	Popularity         float64         `json:"popularity"`
	CombinedCredits    CombinedCredits `json:"combined_credits"`
}
