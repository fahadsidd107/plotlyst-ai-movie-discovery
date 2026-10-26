# Plotlyst API contracts

Step 1 preserves existing successful response bodies so the current React client remains compatible while standardizing every failure response.

## Base URL

Local development: `http://localhost:8080`

## Error envelope

All HTTP errors use this shape:

```json
{
  "error": {
    "code": "movie_not_found",
    "message": "The requested movie was not found",
    "details": null
  }
}
```

- `code` is a stable, machine-readable identifier.
- `message` is safe to display to a user.
- `details` is optional and currently used only for development-safe validation details.
- Internal database errors, secrets, tokens, and provider responses are not exposed.

## Health contracts

### `GET /health/live`

Confirms that the HTTP process is running. It does not check dependencies.

```json
{"service":"plotlyst-api","status":"ok"}
```

### `GET /health/ready`

Pings MongoDB with a two-second deadline. It returns `200` when ready and `503` using the standard error envelope when unavailable.

```json
{"database":"connected","service":"plotlyst-api","status":"ready"}
```

## Current success contracts

The current API predates versioning. These shapes remain stable during Step 1:

- Catalogue endpoints return arrays of movie or genre documents.
- Detail endpoints return a single movie document.
- Login returns the public user profile; tokens remain HTTP-only cookies.
- Registration and movie creation return MongoDB insert results.
- Review updates return `admin_review` and `ranking_name`.
- Session actions return a `message` string.

## TMDB catalogue contracts

The Step 2 catalogue is public and namespaced separately so the legacy MongoDB routes remain compatible:

- `GET /api/catalog/movies?page=1&collection=popular|trending`
- `GET /api/catalog/movies/:id`
- `GET /api/catalog/tv?page=1&collection=popular|trending`
- `GET /api/catalog/tv/:id`
- `GET /api/catalog/tv/:id/seasons/:season`
- `GET /api/catalog/tv/:id/seasons/:season/episodes/:episode`
- `GET /api/catalog/people/trending?page=1`
- `GET /api/catalog/people/born-today?timezone=Europe/Berlin`
- `GET /api/catalog/people/:id`

Success responses use a provider-aware envelope:

```json
{
  "data": {},
  "meta": {
    "source": "TMDB",
    "image_base_url": "https://image.tmdb.org/t/p",
    "attribution": "This product uses the TMDB API but is not endorsed or certified by TMDB."
  }
}
```

List results use distinct movie, television and person summary DTOs. Detail responses use dedicated movie, series, season, episode and person DTOs. TMDB failures are mapped to stable `tmdb_unavailable`, `tmdb_rate_limited`, `tmdb_upstream_error`, or `catalog_item_not_found` errors without leaking the upstream response or API key.
