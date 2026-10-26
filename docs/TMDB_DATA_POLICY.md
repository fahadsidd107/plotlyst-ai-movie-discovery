# TMDB data policy

Plotlyst uses the TMDB v3 API from the Go backend. The API key is never shipped to the browser. Responses are cached in memory for ten minutes and are not persisted as a catalogue snapshot.

## Attribution and permitted use

- Plotlyst displays the required TMDB attribution in its footer/credits area and links to [TMDB](https://www.themoviedb.org/).
- Streaming availability is supplied through TMDB's JustWatch partnership; any screen showing provider data must credit JustWatch.
- The developer API is suitable for this non-commercial portfolio project. Commercial use requires a separate written agreement with TMDB.
- TMDB terms limit caching to six months. Plotlyst's ten-minute volatile cache is deliberately far below that limit.

## AI boundary

TMDB's current API terms prohibit using TMDB content in or for training an AI/ML application without written authorization. Plotlyst therefore does **not** send TMDB descriptions, credits, images, provider data, or catalogue payloads to OpenAI or another model. The existing optional review classifier processes only administrator-authored review text and Plotlyst's local ranking labels.

Any future AI feature that uses TMDB-derived content is blocked until the project owner obtains written permission or replaces that input with a provider whose licence explicitly permits the intended AI processing.

## Source behavior

- Images use TMDB image paths and the configured image CDN base URL.
- Provider availability is region-specific and is not a guarantee that a title is currently playable.
- “Born Today” samples the first three pages of TMDB's popular people feed, then checks full person details against the requested IANA timezone. The API response exposes this coverage limitation instead of presenting it as exhaustive.
- TMDB has no service-level agreement; Plotlyst surfaces provider errors, applies bounded retries, honours `Retry-After`, and keeps stale data out of persistent storage.

## Primary references

- [TMDB API terms](https://www.themoviedb.org/api-terms-of-use)
- [TMDB API FAQ and attribution](https://developer.themoviedb.org/docs/faq)
- [TMDB application authentication](https://developer.themoviedb.org/docs/authentication-application)
- [TMDB image basics](https://developer.themoviedb.org/docs/image-basics)
- [TMDB append-to-response](https://developer.themoviedb.org/docs/append-to-response)
- [TMDB watch-provider attribution](https://developer.themoviedb.org/reference/tv-series-watch-providers)
