# Plotlyst API contracts

## Community and personal library

All `/api/me`, rating-write, review-write and moderation routes use the HTTP-only access cookie. Identity and ownership come from the validated server session, never from a user ID supplied in the body.

| Method | Route | Purpose |
|---|---|---|
| `GET` | `/api/auth/session` | Hydrate the authoritative signed-in user |
| `GET/PUT` | `/api/me/profile` | Read or update profile, avatar URL and preferences |
| `GET` | `/api/profiles/:username` | Read a privacy-shaped public profile |
| `GET` | `/api/community/:mediaType/:mediaId/ratings` | Rating average, count, distribution and viewer rating |
| `PUT/DELETE` | `/api/community/:mediaType/:mediaId/rating` | Upsert or remove the caller's one-to-five-star rating |
| `GET/POST` | `/api/community/:mediaType/:mediaId/reviews` | Read visible reviews or create a programme note |
| `PATCH/DELETE` | `/api/reviews/:reviewId` | Owner-scoped review editing and soft deletion |
| `POST` | `/api/reviews/:reviewId/replies` | Add a one-level reply |
| `PUT/DELETE` | `/api/reviews/:reviewId/helpful` | Add or remove the caller's unique helpful vote |
| `POST` | `/api/reviews/:reviewId/reports` | Privately report a review once |
| `GET/PUT/DELETE` | `/api/me/watchlist[/:mediaType/:mediaId]` | Private watchlist operations |
| `GET/POST/PATCH/DELETE` | `/api/me/playlists[...]` | Owner-scoped custom playlists and items |
| `GET/PUT/DELETE` | `/api/me/history[...]` | Watched-history operations |
| `GET/PUT` | `/api/me/tv/:showId/.../progress` | Show and episode progress |
| `GET/PATCH` | `/api/moderation/...` | MODERATOR or ADMIN report review; suspension is ADMIN-only |

Canonical media references use `media_type` (`movie` or `tv`) plus a positive TMDB `media_id`. User text is stored and rendered as plain text. Duplicate ratings, votes, reports and library entries are prevented by compound indexes.

## Specialist data

| Method | Route | Purpose |
|---|---|---|
| `GET` | `/api/specialist/providers` | Provider configuration and license states |
| `GET` | `/api/specialist/availability/:mediaType/:id?region=DE` | JustWatch-backed subscription, rental, purchase and free offers |
| `GET` | `/api/specialist/awards/organizations` | Verified imported award organizations |
| `GET` | `/api/specialist/awards/:organization/:year` | Categories, nominees and winners for a sourced ceremony |
| `GET` | `/api/specialist/box-office?period=weekend&date=YYYY-MM-DD` | Latest licensed US snapshot not later than the requested date |
| `GET` | `/api/specialist/editorial` | Imported classifications with record-level provenance |
| `POST` | `/api/admin/imports/{awards,box-office,editorial}` | ADMIN-only verified import; supports `dry_run=true` |

Specialist responses use `available`, `partial` or `unavailable` source status and never invent missing facts. Full policy: [`SPECIALIST_DATA_POLICY.md`](SPECIALIST_DATA_POLICY.md).

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

## Discovery contracts

- `GET /api/discovery/home` returns independently named homepage sections and preserves section-level errors.
- `GET /api/discovery/home/:section` supports progressive loading for `now-playing`, `trending-movies`, `trending-tv`, `fan-favourites`, and `airing-today`.
- `GET /api/discovery/search?q=` requires two characters and returns only movie, television and person results.
- `GET /api/discovery/movies` and `/tv` accept validated `genre`, `year`, `decade`, `country`, `language`, `runtime_min`, `runtime_max`, `provider`, `region`, `sort`, and `page` values.
- `GET /api/discovery/filters/:media` returns provider-backed genres, countries, languages and regional streaming providers.
- `GET /api/discovery/charts/:media/:size?page=` supports movie/TV and sizes 100, 250 or 2,000. Every response includes formula, eligibility, source and stable tie-break metadata.
- `GET /api/discovery/franchises/:id` labels the result as an official external franchise.
- `GET /api/discovery/collection-kinds` defines official, editorial and user-owned collection boundaries.

See [`RANKING_METHODOLOGY.md`](RANKING_METHODOLOGY.md) for the exact scoring contract.
