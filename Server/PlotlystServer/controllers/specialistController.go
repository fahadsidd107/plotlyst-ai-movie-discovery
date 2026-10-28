package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"plotlyst/server/apiresponse"
	"plotlyst/server/database"
	"plotlyst/server/specialist"
)

type SpecialistController struct {
	client       *mongo.Client
	availability specialist.AvailabilityProvider
}

func NewSpecialistController(client *mongo.Client, availability specialist.AvailabilityProvider) *SpecialistController {
	return &SpecialistController{client: client, availability: availability}
}
func (s *SpecialistController) col(name string) *mongo.Collection {
	return database.OpenCollection(name, s.client)
}
func unavailable(provider, coverage string) specialist.SourceMeta {
	return specialist.SourceMeta{Provider: provider, Status: specialist.StatusUnavailable, Coverage: coverage}
}

func (s *SpecialistController) Providers(c *gin.Context) {
	c.JSON(200, gin.H{"data": []gin.H{
		{"kind": "availability", "provider": "TMDB / JustWatch", "status": "available", "license": "TMDB API terms; JustWatch attribution required"},
		{"kind": "awards", "provider": "Verified admin imports", "status": "partial", "license": "Per-import source manifest"},
		{"kind": "box_office", "provider": "Licensed admin imports", "status": "partial", "license": "Per-import source manifest"},
		{"kind": "editorial", "provider": "Editorial source manifests", "status": "partial", "license": "Per-import source manifest"},
	}})
}
func (s *SpecialistController) Availability(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	kind := c.Param("mediaType")
	if err != nil || id < 1 || (kind != "movie" && kind != "tv") {
		apiresponse.Error(c, 400, "invalid_media", "Valid movie or television media is required")
		return
	}
	result, err := s.availability.Get(c.Request.Context(), specialist.MediaRef{MediaType: kind, MediaID: id}, c.DefaultQuery("region", "DE"))
	if err != nil {
		apiresponse.Error(c, 502, "availability_upstream_error", "Availability provider is temporarily unavailable")
		return
	}
	c.JSON(200, gin.H{"data": result})
}
func (s *SpecialistController) AwardOrganizations(c *gin.Context) {
	ctx, cancel := requestContext(c)
	defer cancel()
	cur, err := s.col("award_organizations").Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		apiresponse.Error(c, 500, "database_error", "Unable to load awards")
		return
	}
	rows := []specialist.AwardOrganization{}
	defer cur.Close(ctx)
	cur.All(ctx, &rows)
	source := unavailable("Verified awards imports", "No verified awards source has been imported yet.")
	if len(rows) > 0 {
		source = specialist.SourceMeta{Provider: "Verified awards imports", Status: specialist.StatusAvailable, Coverage: "Organizations with traceable source manifests."}
	}
	c.JSON(200, gin.H{"data": rows, "source": source})
}
func (s *SpecialistController) AwardCeremony(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		apiresponse.Error(c, 400, "invalid_year", "A valid ceremony year is required")
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	var row specialist.AwardCeremony
	if s.col("award_ceremonies").FindOne(ctx, bson.M{"organization_id": c.Param("organization"), "year": year}).Decode(&row) != nil {
		c.JSON(200, gin.H{"data": nil, "source": unavailable("Verified awards imports", "No verified ceremony data is available for this organization and year.")})
		return
	}
	c.JSON(200, gin.H{"data": row, "source": row.Source})
}
func (s *SpecialistController) BoxOffice(c *gin.Context) {
	period := c.DefaultQuery("period", "weekend")
	if period != "weekend" && period != "weekly" && period != "daily" {
		apiresponse.Error(c, 400, "invalid_period", "Period must be daily, weekend, or weekly")
		return
	}
	asOf := time.Now().UTC()
	if raw := c.Query("date"); raw != "" {
		var err error
		asOf, err = time.Parse("2006-01-02", raw)
		if err != nil {
			apiresponse.Error(c, 400, "invalid_date", "Date must use YYYY-MM-DD")
			return
		}
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	var row specialist.BoxOfficeSnapshot
	err := s.col("box_office_snapshots").FindOne(ctx, bson.M{"period": period, "period_end": bson.M{"$lte": asOf}}, options.FindOne().SetSort(bson.D{{Key: "period_end", Value: -1}})).Decode(&row)
	if err != nil {
		c.JSON(200, gin.H{"data": nil, "source": unavailable("Licensed box-office import", "No licensed US box-office snapshot is configured for this reporting period.")})
		return
	}
	c.JSON(200, gin.H{"data": row, "source": row.Source})
}
func (s *SpecialistController) Editorial(c *gin.Context) {
	filter := bson.M{}
	if kind := c.Query("media_type"); kind != "" {
		filter["media_type"] = kind
	}
	if raw := c.Query("media_id"); raw != "" {
		if id, e := strconv.Atoi(raw); e == nil {
			filter["media_id"] = id
		}
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	cur, _ := s.col("editorial_records").Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "imported_at", Value: -1}}).SetLimit(200))
	rows := []specialist.EditorialRecord{}
	if cur != nil {
		defer cur.Close(ctx)
		cur.All(ctx, &rows)
	}
	source := unavailable("Editorial imports", "No matching editorial classification has been imported.")
	if len(rows) > 0 {
		source = specialist.SourceMeta{Provider: "Editorial imports", Status: specialist.StatusAvailable, Coverage: "Imported classifications with record-level provenance."}
	}
	c.JSON(200, gin.H{"data": rows, "source": source})
}

type specialistImport struct {
	Source        specialist.SourceMeta          `json:"source"`
	Organizations []specialist.AwardOrganization `json:"organizations"`
	Ceremonies    []specialist.AwardCeremony     `json:"ceremonies"`
	BoxOffice     []specialist.BoxOfficeSnapshot `json:"snapshots"`
	Editorial     []specialist.EditorialRecord   `json:"records"`
}

func validateSource(src specialist.SourceMeta) bool {
	return strings.TrimSpace(src.Provider) != "" && strings.HasPrefix(src.SourceURL, "https://") && strings.TrimSpace(src.License) != ""
}
func checksum(v any) string {
	raw, _ := json.Marshal(v)
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
func (s *SpecialistController) Import(kind string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload specialistImport
		if c.ShouldBindJSON(&payload) != nil || !validateSource(payload.Source) {
			apiresponse.Error(c, 400, "invalid_source_manifest", "Provider, HTTPS source URL and license are required")
			return
		}
		payload.Source.RetrievedAt = time.Now().UTC()
		payload.Source.Status = specialist.StatusAvailable
		if c.Query("dry_run") == "true" {
			c.JSON(200, gin.H{"valid": true, "checksum": checksum(payload)})
			return
		}
		ctx, cancel := requestContext(c)
		defer cancel()
		count := 0
		switch kind {
		case "awards":
			for _, v := range payload.Organizations {
				v.Source = payload.Source
				s.col("award_organizations").UpdateOne(ctx, bson.M{"organization_id": v.ID}, bson.M{"$set": v}, options.UpdateOne().SetUpsert(true))
				count++
			}
			for _, v := range payload.Ceremonies {
				v.Source = payload.Source
				s.col("award_ceremonies").UpdateOne(ctx, bson.M{"organization_id": v.OrganizationID, "year": v.Year}, bson.M{"$set": v}, options.UpdateOne().SetUpsert(true))
				count++
			}
		case "box-office":
			for _, v := range payload.BoxOffice {
				v.Source = payload.Source
				v.Currency = "USD"
				s.col("box_office_snapshots").UpdateOne(ctx, bson.M{"source.provider": payload.Source.Provider, "period": v.Period, "period_end": v.PeriodEnd}, bson.M{"$set": v}, options.UpdateOne().SetUpsert(true))
				count++
			}
		case "editorial":
			for _, v := range payload.Editorial {
				if v.MediaType != "movie" && v.MediaType != "tv" {
					continue
				}
				v.Source = payload.Source
				v.ImportedAt = time.Now().UTC()
				s.col("editorial_records").UpdateOne(ctx, bson.M{"source.provider": payload.Source.Provider, "external_id": v.ID, "classification_type": v.ClassificationType, "region": v.Region}, bson.M{"$set": v}, options.UpdateOne().SetUpsert(true))
				count++
			}
		}
		c.JSON(http.StatusCreated, gin.H{"imported": count, "checksum": checksum(payload)})
	}
}
