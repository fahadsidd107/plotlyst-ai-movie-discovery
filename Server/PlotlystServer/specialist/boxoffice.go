package specialist

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// DemoBoxOffice is deliberately synthetic. It exists so every UI state can be
// developed without presenting made-up figures as reporting data.
type DemoBoxOffice struct{}

func (DemoBoxOffice) USSnapshot(_ context.Context, period string, asOf time.Time) (BoxOfficeSnapshot, error) {
	if period != "daily" && period != "weekly" && period != "weekend" {
		return BoxOfficeSnapshot{}, ErrNoCoverage
	}
	end := asOf.UTC().Truncate(24 * time.Hour)
	days := 1
	if period == "weekend" {
		days = 3
	} else if period == "weekly" {
		days = 7
	}
	return BoxOfficeSnapshot{
		Period: period, PeriodStart: end.AddDate(0, 0, -(days - 1)), PeriodEnd: end, AsOf: end, Currency: "USD", Demo: true,
		Entries: []BoxOfficeEntry{
			{Rank: 1, Title: "Example Premiere", WeekendGross: 18400000, DomesticGross: 57200000, Theaters: 3820},
			{Rank: 2, Title: "The Second Reel", WeekendGross: 12100000, DomesticGross: 33900000, Theaters: 3410},
			{Rank: 3, Title: "Midnight Matinee", WeekendGross: 7800000, DomesticGross: 21400000, Theaters: 2860},
		},
		Source: SourceMeta{Provider: "Plotlyst demo fixture", License: "DEMO DATA — NOT VERIFIED OR FOR REPORTING", Attribution: "Synthetic values for development only", RetrievedAt: time.Now().UTC(), ReportingPeriod: period, Coverage: "Illustrative fallback; titles and values are synthetic.", Status: StatusDemo},
	}, nil
}

// OpusDataProvider isolates the licensed transport. BaseURL must be the exact
// endpoint supplied by the customer's OpusData contract; no endpoint is guessed.
type OpusDataProvider struct {
	BaseURL, APIKey string
	Client          *http.Client
}

func (p OpusDataProvider) USSnapshot(ctx context.Context, period string, asOf time.Time) (BoxOfficeSnapshot, error) {
	if p.BaseURL == "" || p.APIKey == "" {
		return BoxOfficeSnapshot{}, ErrNotConfigured
	}
	u, err := url.Parse(strings.TrimRight(p.BaseURL, "/") + "/box-office/us")
	if err != nil {
		return BoxOfficeSnapshot{}, err
	}
	q := u.Query()
	q.Set("period", period)
	q.Set("as_of", asOf.Format("2006-01-02"))
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return BoxOfficeSnapshot{}, err
	}
	req.Header.Set("Authorization", "Bearer "+p.APIKey)
	req.Header.Set("Accept", "application/json")
	client := p.Client
	if client == nil {
		client = &http.Client{Timeout: 8 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return BoxOfficeSnapshot{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return BoxOfficeSnapshot{}, fmt.Errorf("opusdata returned %d", resp.StatusCode)
	}
	var out BoxOfficeSnapshot
	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&out); err != nil {
		return BoxOfficeSnapshot{}, err
	}
	if err := validateSnapshot(out); err != nil {
		return BoxOfficeSnapshot{}, err
	}
	out.Demo = false
	out.Source = SourceMeta{Provider: "OpusData / The Numbers", SourceURL: u.String(), License: "Licensed OpusData API", Attribution: "Box-office data supplied under an OpusData license", RetrievedAt: time.Now().UTC(), ReportingPeriod: period, Coverage: "Licensed US box-office reporting snapshot.", Status: StatusAvailable}
	return out, nil
}

func validateSnapshot(v BoxOfficeSnapshot) error {
	if v.Currency != "USD" || len(v.Entries) == 0 {
		return errors.New("invalid OpusData snapshot")
	}
	seen := map[int]bool{}
	for _, e := range v.Entries {
		if e.Rank < 1 || seen[e.Rank] || e.WeekendGross < 0 || e.DomesticGross < 0 {
			return errors.New("invalid OpusData entry")
		}
		seen[e.Rank] = true
	}
	return nil
}

// Marshal helper used by contract tests and documented sample adapters.
func EncodeSnapshot(v BoxOfficeSnapshot) []byte {
	b := bytes.Buffer{}
	_ = json.NewEncoder(&b).Encode(v)
	return b.Bytes()
}
func ParseAsOf(raw string) (time.Time, error) {
	if raw == "" {
		return time.Now().UTC(), nil
	}
	return time.Parse("2006-01-02", raw)
}
func Int64(raw string) int64 { v, _ := strconv.ParseInt(raw, 10, 64); return v }
