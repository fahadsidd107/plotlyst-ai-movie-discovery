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

API versioning and DTO envelopes will be introduced alongside the TMDB catalogue adapter, with compatibility maintained for the current client during migration.
