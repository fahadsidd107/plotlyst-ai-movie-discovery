package contracts

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestOpenAPIContractIsValidAndCoversCriticalRoutes(t *testing.T) {
	raw, err := os.ReadFile("../../../docs/openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var spec struct {
		OpenAPI    string                    `yaml:"openapi"`
		Info       map[string]any            `yaml:"info"`
		Paths      map[string]map[string]any `yaml:"paths"`
		Components map[string]any            `yaml:"components"`
	}
	if err := yaml.Unmarshal(raw, &spec); err != nil {
		t.Fatalf("invalid OpenAPI YAML: %v", err)
	}
	if spec.OpenAPI != "3.1.0" || spec.Info["title"] != "Plotlyst API" {
		t.Fatalf("unexpected contract header: %#v", spec)
	}
	required := map[string][]string{"/health/live": {"get"}, "/api/catalog/movies": {"get"}, "/api/discovery/search": {"get"}, "/api/ai/discover": {"post"}, "/api/analytics/cinema-dna": {"get"}, "/api/story-dna/movies/{id}": {"get"}, "/api/story-dna/discover": {"post"}, "/api/me/spoiler-preferences": {"get", "put"}, "/api/community/reviews/{reviewId}/reveal": {"get"}, "/api/debate/movies/{id}": {"get"}, "/api/movie-night/rooms": {"post"}, "/api/specialist/box-office": {"get"}, "/api/me/profile": {"get", "put"}}
	for path, methods := range required {
		ops, ok := spec.Paths[path]
		if !ok {
			t.Errorf("critical path missing: %s", path)
			continue
		}
		for _, method := range methods {
			if _, ok := ops[method]; !ok {
				t.Errorf("%s %s missing", method, path)
			}
		}
	}
	if len(spec.Components) == 0 {
		t.Fatal("components/schemas must be documented")
	}
}
