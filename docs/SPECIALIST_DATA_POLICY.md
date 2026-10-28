# Specialist data policy

Plotlyst never substitutes generated values for awards, box-office totals, streaming availability or editorial classifications.

## Providers

- Streaming availability is requested from TMDB's watch-provider endpoints. The underlying availability data is supplied by JustWatch and must display JustWatch attribution. Results are country-level and distinguish subscription, rental, purchase and free/ad-supported offers.
- Awards are accepted only through an administrator import carrying the organization, HTTPS source URL, license/permission statement and retrieval timestamp. Wikidata structured data is suitable for a future automated adapter because it is released under CC0, but imported records must still retain their record-level source.
- US box-office snapshots are accepted only from a feed the operator is licensed to use. No public credential is bundled. Without a licensed import, the API returns `unavailable` instead of zero revenue or invented rankings.
- Editorial classifications use versioned administrator imports with the same source-manifest requirements.

## Coverage states

- `available`: a verified provider returned records for the requested scope.
- `partial`: the provider is configured, but the requested country, title, ceremony or period has incomplete coverage.
- `unavailable`: no licensed/verified source has been configured or imported.

Every specialist response includes provider, source URL where applicable, license statement, attribution, retrieval time, reporting period and a human-readable coverage note. Imports support `dry_run=true` and return a deterministic checksum before persistence.
