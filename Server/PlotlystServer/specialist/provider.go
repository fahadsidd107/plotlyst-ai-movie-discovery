package specialist

import (
	"context"
	"errors"
	"time"
)

var ErrNotConfigured = errors.New("provider not configured")
var ErrNoCoverage = errors.New("provider has no coverage")

type Status string

const (
	StatusAvailable   Status = "available"
	StatusPartial     Status = "partial"
	StatusUnavailable Status = "unavailable"
)

type SourceMeta struct {
	Provider        string    `json:"provider" bson:"provider"`
	SourceURL       string    `json:"source_url,omitempty" bson:"source_url,omitempty"`
	License         string    `json:"license,omitempty" bson:"license,omitempty"`
	Attribution     string    `json:"attribution,omitempty" bson:"attribution,omitempty"`
	RetrievedAt     time.Time `json:"retrieved_at,omitempty" bson:"retrieved_at,omitempty"`
	ReportingPeriod string    `json:"reporting_period,omitempty" bson:"reporting_period,omitempty"`
	Coverage        string    `json:"coverage,omitempty" bson:"coverage,omitempty"`
	Status          Status    `json:"status" bson:"status"`
}
type MediaRef struct {
	MediaType string `json:"media_type" bson:"media_type"`
	MediaID   int    `json:"media_id" bson:"media_id"`
}
type Offer struct {
	ProviderID   int    `json:"provider_id"`
	ProviderName string `json:"provider_name"`
	LogoPath     string `json:"logo_path"`
	Type         string `json:"type"`
}
type AvailabilityResult struct {
	MediaRef
	Region string     `json:"region"`
	Link   string     `json:"link,omitempty"`
	Offers []Offer    `json:"offers"`
	Source SourceMeta `json:"source"`
}
type AwardOrganization struct {
	ID      string     `json:"id" bson:"organization_id"`
	Name    string     `json:"name" bson:"name"`
	Website string     `json:"website,omitempty" bson:"website,omitempty"`
	Source  SourceMeta `json:"source" bson:"source"`
}
type AwardCeremony struct {
	ID             string          `json:"id" bson:"ceremony_id"`
	OrganizationID string          `json:"organization_id" bson:"organization_id"`
	Name           string          `json:"name" bson:"name"`
	Year           int             `json:"year" bson:"year"`
	Edition        int             `json:"edition,omitempty" bson:"edition,omitempty"`
	Categories     []AwardCategory `json:"categories" bson:"categories"`
	Source         SourceMeta      `json:"source" bson:"source"`
}
type AwardCategory struct {
	Name     string         `json:"name" bson:"name"`
	Nominees []AwardNominee `json:"nominees" bson:"nominees"`
}
type AwardNominee struct {
	Name      string `json:"name" bson:"name"`
	Winner    bool   `json:"winner" bson:"winner"`
	MediaType string `json:"media_type,omitempty" bson:"media_type,omitempty"`
	MediaID   int    `json:"media_id,omitempty" bson:"media_id,omitempty"`
	PersonID  int    `json:"person_id,omitempty" bson:"person_id,omitempty"`
}
type BoxOfficeEntry struct {
	Rank          int    `json:"rank" bson:"rank"`
	Title         string `json:"title" bson:"title"`
	MediaID       int    `json:"media_id,omitempty" bson:"media_id,omitempty"`
	WeekendGross  int64  `json:"weekend_gross,omitempty" bson:"weekend_gross,omitempty"`
	DomesticGross int64  `json:"domestic_gross,omitempty" bson:"domestic_gross,omitempty"`
	Theaters      int    `json:"theaters,omitempty" bson:"theaters,omitempty"`
}
type BoxOfficeSnapshot struct {
	Period      string           `json:"period" bson:"period"`
	PeriodStart time.Time        `json:"period_start" bson:"period_start"`
	PeriodEnd   time.Time        `json:"period_end" bson:"period_end"`
	AsOf        time.Time        `json:"as_of" bson:"as_of"`
	Currency    string           `json:"currency" bson:"currency"`
	Entries     []BoxOfficeEntry `json:"entries" bson:"entries"`
	Source      SourceMeta       `json:"source" bson:"source"`
}
type EditorialRecord struct {
	ID                 string `json:"id" bson:"external_id"`
	MediaRef           `bson:",inline"`
	ClassificationType string     `json:"classification_type" bson:"classification_type"`
	Value              string     `json:"value" bson:"value"`
	Label              string     `json:"label" bson:"label"`
	Region             string     `json:"region,omitempty" bson:"region,omitempty"`
	Notes              string     `json:"notes,omitempty" bson:"notes,omitempty"`
	Source             SourceMeta `json:"source" bson:"source"`
	ImportedAt         time.Time  `json:"imported_at" bson:"imported_at"`
}

type AvailabilityProvider interface {
	Get(context.Context, MediaRef, string) (AvailabilityResult, error)
}
type AwardsProvider interface {
	Organizations(context.Context) ([]AwardOrganization, SourceMeta, error)
	Ceremony(context.Context, string, int) (AwardCeremony, error)
}
type BoxOfficeProvider interface {
	USSnapshot(context.Context, string, time.Time) (BoxOfficeSnapshot, error)
}
