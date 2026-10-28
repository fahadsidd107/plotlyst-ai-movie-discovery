# Plotlyst implementation roadmap

This plan translates the master product brief into manageable vertical slices. A phase is complete only when its user interface, backend, data source, persistence, validation, error handling, responsive behavior, tests, and documentation are working together.

## Step 1 — Foundation and product identity

✅ Rename the product and internal namespaces to Plotlyst.  
✅ Remove the previous playback integration and identifiers.  
✅ Establish a one-command local environment and seed workflow.  
✅ Repair local JWT/cookie behavior and support AI-free review updates.  
✅ Audit API contracts, introduce structured errors, add health checks and graceful shutdown.  
✅ Create Midnight Pop design tokens, typography, reusable cards, skeletons and empty states.  
✅ Add environment validation and `.env.example` files without secrets.

**Outcome:** a stable, branded base that is safe to extend.

## Step 2 — Real movies, television and people

✅ Build a Go integration client for TMDB with timeouts, retries, caching and rate-limit handling.
✅ Add separate summary/detail DTOs for movies, television, seasons, episodes and people.
✅ Implement movie and television landing pages plus complete detail pages.
✅ Add cast, crew, franchises, release data, recommendations and regional providers.
✅ Implement people profiles, trending people and date-aware “Born Today”.
✅ Create a compatibility migration from legacy editorial rankings to five-star semantics.
✅ Verify licensing, attribution and AI-processing permissions before using provider data.

**Outcome:** real, navigable catalogue data supplied through the Go backend.

**Runtime verification:** implementation and automated tests are complete. The TMDB credential currently configured in the local `.env` is rejected by TMDB as invalid, so live catalogue rendering remains blocked until an active v3 API key is supplied.

## Step 3 — Discovery, search and charts

✅ Replace the simple grid with an aggregated, progressively loaded cinematic homepage.
✅ Build global autocomplete for movies, shows and people.
✅ Build backend-driven filters for genre, year, country, language, runtime and provider.
✅ Implement reproducible Top 100, Top 250 and paginated Top 2,000 charts.
✅ Add television charts, fan favourites, trending, year/decade and genre destinations.
✅ Separate official franchises, editorial collections and user-created lists.
✅ Document ranking formulas, thresholds, source provenance and stable tie-breaking.

**Outcome:** a deep discovery product rather than a static catalogue.

**Reproducibility boundary:** chart ordering is deterministic for the same cached TMDB payload. Because provider votes change, immutable historical/nightly chart snapshots remain a production-hardening task rather than being falsely claimed by the live chart API.

## Step 4 — Community and personal library

✅ Harden authentication and introduce USER, MODERATOR and ADMIN authorization.
✅ Add profiles, avatars, preferences, privacy controls and regional settings.
✅ Implement one five-star rating per user/media item with aggregate distributions.
✅ Implement review create/edit/delete, spoilers, helpful votes, replies and reporting.
✅ Add private watchlists and independently owned custom playlists.
✅ Add watched history plus season and episode progress tracking.
✅ Add moderation, sanitization, rate limits, compound indexes and ownership tests.

**Outcome:** authenticated activity persists correctly and remains owner-controlled.

**Security boundary:** registration always creates a USER; elevated roles remain server-managed. Community writes use authoritative cookie identity, active-account checks, owner-scoped filters, unique compound indexes, plain-text sanitization and endpoint rate limits. Avatar fields accept URLs; binary avatar upload requires separately configured object storage.

## Step 5 — Verified specialist data

✅ Define provider interfaces for awards, box office and streaming availability.
✅ Implement regional streaming filters with subscription/rental/purchase distinctions.
✅ Add awards organizations, ceremonies, categories, nominees and winners only from verified sources.
✅ Add reporting-period-aware US box-office snapshots from a licensed source.
✅ Build honest unavailable and partial-coverage states when a source is missing.
✅ Add editorial imports for classifications not supplied by the main catalogue provider.

**Outcome:** specialist pages backed by traceable data rather than invented values.

**Coverage boundary:** availability is live TMDB/JustWatch data. Awards, box-office and editorial records require administrator imports with a source URL and license manifest. Until a licensed source is imported, these APIs deliberately return `unavailable`; they never manufacture records or represent missing revenue as zero. See [`SPECIALIST_DATA_POLICY.md`](SPECIALIST_DATA_POLICY.md).

## Step 6 — AI, analytics and live experiences

⬜ Split AI into intent extraction, validated filters, retrieval, ranking and explanation services.  
⬜ Implement natural-language discovery with typed structured output.  
⬜ Add movie-memory search and mood/runtime/language-aware discovery.  
⬜ Generate grounded recommendation explanations from real user activity.  
⬜ Build Cinema DNA analytics with MongoDB aggregations.  
⬜ Add Group Movie Night and preference intersection.  
⬜ Introduce WebSocket infrastructure only where live notifications or collaborative sessions justify it.  
⬜ Add token/cost budgets, timeouts, caching, provenance, privacy controls and prompt-injection defenses.

**Outcome:** real, explainable AI and live features—not relabeled keyword search.

## Step 7 — Production readiness

⬜ Add Go unit, integration and repository tests with external-provider mocks.  
⬜ Add frontend component, interaction and responsive-navigation tests.  
⬜ Add OpenAPI documentation and contract validation.  
⬜ Add structured logging, request IDs, metrics and operational health endpoints.  
⬜ Optimize homepage aggregation, indexes, caching, lazy loading and request deduplication.  
⬜ Add Docker configuration, CI quality gates and deployment documentation.  
⬜ Deploy the client to Vercel and the containerized API to persistent Go hosting.  
⬜ Complete accessibility, security, performance and responsive-device audits.

**Outcome:** a documented, observable and deployable portfolio application.

## Definition of done for every feature

- Reachable through real navigation
- Backed by a working API and permitted real data or genuine user data
- Correctly authorized and validated
- Persistent behavior survives refreshes
- Loading, empty and error states are present
- Usable on mobile and with keyboard navigation
- Important logic is covered by automated tests
- Documentation states its data source and limitations honestly
