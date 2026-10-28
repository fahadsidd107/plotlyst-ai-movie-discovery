package specialist

import (
	"context"
	"fmt"
	"strings"
	"time"

	"plotlyst/server/tmdb"
)

type TMDBAvailability struct{ Client *tmdb.Client }

func (p TMDBAvailability) Get(ctx context.Context, ref MediaRef, region string) (AvailabilityResult, error) {
	region = strings.ToUpper(region)
	if len(region) != 2 {
		return AvailabilityResult{}, fmt.Errorf("region must be an ISO 3166-1 alpha-2 code")
	}
	var raw tmdb.WatchProviders
	if err := p.Client.Get(ctx, fmt.Sprintf("/%s/%d/watch/providers", ref.MediaType, ref.MediaID), nil, &raw); err != nil {
		return AvailabilityResult{}, err
	}
	row, ok := raw.Results[region]
	status := StatusAvailable
	coverage := "Country-level streaming, rental and purchase availability supplied by JustWatch via TMDB."
	if !ok {
		status = StatusPartial
		coverage = "No availability listing is currently supplied for this title and region."
	}
	offers := []Offer{}
	appendOffers := func(kind string, items []tmdb.Provider) {
		for _, item := range items {
			offers = append(offers, Offer{ProviderID: item.ProviderID, ProviderName: item.ProviderName, LogoPath: item.LogoPath, Type: kind})
		}
	}
	appendOffers("subscription", row.Flatrate)
	appendOffers("rental", row.Rent)
	appendOffers("purchase", row.Buy)
	appendOffers("free", row.Free)
	return AvailabilityResult{MediaRef: ref, Region: region, Link: row.Link, Offers: offers, Source: SourceMeta{Provider: "TMDB / JustWatch", SourceURL: row.Link, License: "TMDB API terms; JustWatch attribution required", Attribution: "Streaming availability data by JustWatch.", RetrievedAt: time.Now().UTC(), Coverage: coverage, Status: status}}, nil
}
