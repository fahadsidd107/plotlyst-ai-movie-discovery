package specialist

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestDemoBoxOfficeIsAlwaysExplicit(t *testing.T) {
	v, err := (DemoBoxOffice{}).USSnapshot(context.Background(), "weekend", time.Date(2026, 10, 28, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if !v.Demo || v.Source.Status != StatusDemo || v.Source.Provider != "Plotlyst demo fixture" {
		t.Fatalf("demo provenance missing: %#v", v.Source)
	}
	if v.Source.License == "" || len(v.Entries) == 0 {
		t.Fatal("demo warning or fixtures missing")
	}
}

func TestOpusDataProviderMapsLicensedContract(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") != "Bearer licensed-key" {
			t.Errorf("missing auth")
		}
		if r.URL.Query().Get("period") != "weekend" {
			t.Errorf("missing period")
		}
		var body bytes.Buffer
		_ = json.NewEncoder(&body).Encode(BoxOfficeSnapshot{Period: "weekend", Currency: "USD", Entries: []BoxOfficeEntry{{Rank: 1, Title: "Licensed title", WeekendGross: 10}}})
		return &http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(&body)}, nil
	})}
	provider := OpusDataProvider{BaseURL: "https://licensed.example", APIKey: "licensed-key", Client: client}
	snapshot, err := provider.USSnapshot(context.Background(), "weekend", time.Date(2026, 10, 29, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Demo || snapshot.Source.Status != StatusAvailable || snapshot.Source.Provider != "OpusData / The Numbers" {
		t.Fatalf("licensed provenance incorrect: %#v", snapshot)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return fn(request) }
func TestValidateSnapshotRejectsInvalidLicensedPayload(t *testing.T) {
	v := BoxOfficeSnapshot{Currency: "USD", Entries: []BoxOfficeEntry{{Rank: 1, WeekendGross: 2}, {Rank: 1, WeekendGross: 3}}}
	if validateSnapshot(v) == nil {
		t.Fatal("duplicate licensed ranks must be rejected")
	}
}
