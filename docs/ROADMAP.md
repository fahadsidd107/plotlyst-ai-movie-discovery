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

**Coverage boundary:** availability is live TMDB/JustWatch data. Awards and editorial records require administrator imports with a source URL and license manifest. Box office uses licensed OpusData only when credentials and the contracted base URL are configured; local development otherwise shows a conspicuous synthetic `demo` snapshot, never a verified claim. See [`SPECIALIST_DATA_POLICY.md`](SPECIALIST_DATA_POLICY.md).

## Step 6 — AI, analytics and live experiences

✅ Split AI into intent extraction, validated filters, retrieval, ranking and explanation services.
✅ Implement natural-language discovery with typed structured output.
✅ Add movie-memory search and mood/runtime/language-aware discovery.
✅ Generate opt-in grounded recommendation explanations from real user activity counts.
✅ Build Cinema DNA analytics with MongoDB aggregations.
✅ Add Group Movie Night and preference intersection.
✅ Introduce WebSocket infrastructure only for collaborative Movie Night rooms.
✅ Add request budgets, timeouts, caching, provenance, privacy controls and prompt-injection defenses.

**Outcome:** real, explainable AI and live features—not relabeled keyword search.

**AI and coverage boundary:** AI receives only bounded user-written discovery text; TMDB payloads never enter model prompts. With `AI_PROVIDER=offline`, Plotlyst visibly uses a deterministic parser. Personal activity explanations require stored opt-in settings plus request consent. Cinema DNA reports the coverage gap for legacy activity rows that lack genre/runtime/language snapshots. Movie Night REST is authoritative and its authenticated live channel is limited to participants. See [`AI_DATA_POLICY.md`](AI_DATA_POLICY.md) and [`GROUP_MOVIE_NIGHT.md`](GROUP_MOVIE_NIGHT.md).

## Step 7 — Production readiness

✅ Add Go unit, integration and repository tests with external-provider mocks.
✅ Add frontend component, interaction and responsive-navigation tests.
✅ Add OpenAPI documentation and contract validation.
✅ Add structured logging, request IDs, metrics and operational health endpoints.
✅ Optimize homepage aggregation, indexes, caching, lazy loading and request deduplication.
✅ Add Docker configuration, CI quality gates and deployment documentation.
⬜ Deploy the client to Vercel and the containerized API to persistent Go hosting.
✅ Complete the local accessibility, security, performance and responsive-device audit baseline.

**Outcome:** a documented, observable and deployable portfolio application. The repository is deployment-ready; publishing it is the only account-dependent Step 7 action still open.

**Production boundary:** CI runs Go race/unit/contract tests, a MongoDB-backed integration suite, frontend interaction tests, lint, production builds and dependency audits. Docker, Render and Vercel configurations plus deployment and production-audit guides are committed. Final public deployment, deployed Lighthouse runs, physical-device/screen-reader passes and an external security review require the production URLs and hosting accounts, so they are not represented as completed deployment work. See [`DEPLOYMENT.md`](DEPLOYMENT.md) and [`PRODUCTION_AUDIT.md`](PRODUCTION_AUDIT.md).

## Step 8 — Story DNA

✅ Add versioned Story DNA profiles for supported movies.
✅ Classify narrative structure, tone, pacing, complexity and storytelling style.
✅ Persist validated structured characteristics with confidence and provenance.
✅ Display an interactive, spoiler-safe narrative fingerprint on movie detail pages.
✅ Add weighted similarity discovery and natural-language seed searches.
✅ Generate deterministic explanations grounded in compared characteristics.
✅ Add adjustable mood, pacing, complexity and structure priorities.
✅ Document the AI/data boundary and expose honest heuristic limitations.

**Outcome:** users can understand a supported movie’s inferred storytelling profile and find explainable narrative matches beyond genre.

**Classification boundary:** the current classifier is a versioned, deterministic metadata heuristic derived from permitted TMDB genre, runtime, release and vote fields. Its confidence is deliberately capped and the interface labels it as an inferred discovery aid—not objective artistic truth. TMDB plot text is never sent to OpenAI. Ending traits remain unknown and spoiler-hidden until an authorized editorial or licensed classifier supplies them. See [`STORY_DNA_POLICY.md`](STORY_DNA_POLICY.md).

## Step 9 — Spoiler Firewall

✅ Add persistent OFF, STANDARD and STRICT protection preferences.
✅ Add accessible account controls explaining each protection level.
✅ Redact protected reviews and comments server-side before delivery.
✅ Preserve author spoiler tags and source/confidence metadata.
✅ Add keyboard-accessible, deliberate, one-item reveal controls.
✅ Protect movie/television synopses, discovery previews and recommendations in Strict mode.
✅ Use server-owned episode progress to conceal future episode descriptions.
✅ Prevent protected AI Desk overview text from leaking in Strict mode.
✅ Document the conservative moderation and authorized-AI boundary.

**Outcome:** members can browse community and catalogue surfaces using a persistent protection level, with episode visibility grounded in their watched progress and every reveal remaining explicit.

**Detection boundary:** review authors supply the current spoiler label. Strict mode conservatively conceals even untagged community prose because automated detection cannot guarantee safety. AI spoiler detection remains intentionally disabled until review text is authorized for model processing and a moderation-review workflow, retention policy and confidence threshold are configured. See [`SPOILER_FIREWALL.md`](SPOILER_FIREWALL.md).

## Step 10 — Plotlyst Debate

✅ Add “What viewers loved” and “What viewers disliked” movie-detail analysis.
✅ Aggregate genuine Plotlyst ratings and reviews with MongoDB pipelines.
✅ Extract grounded acting, cinematography, pacing, writing, soundtrack and direction themes.
✅ Compare high-, middle- and low-rating evidence.
✅ Calculate and explain polarized rating distributions.
✅ Link every observation to supporting original review identifiers and excerpts.
✅ Keep spoiler-tagged review prose out of summaries.
✅ Add a meaningful five-review minimum and honest insufficient-data state.
✅ Cache summaries and refresh asynchronously after meaningful review changes.
✅ Add an authorized, opt-in LangChainGo/OpenAI analyzer with deterministic fallback.
✅ Build the accessible Plotlyst Debate tab, filters, chart and expandable evidence.

**Outcome:** users can inspect why authentic Plotlyst reviewers agreed or disagreed, with every generated observation traceable to the community evidence that supports it.

**Grounding boundary:** the default analyzer is deterministic and review-grounded. AI analysis is disabled unless the operator explicitly enables it, and then receives only non-spoiler reviews from members who opted into AI activity processing. Fewer than five reviews never produces a sentiment summary. See [`DEBATE_ANALYSIS_POLICY.md`](DEBATE_ANALYSIS_POLICY.md).

## Definition of done for every feature

- Reachable through real navigation
- Backed by a working API and permitted real data or genuine user data
- Correctly authorized and validated
- Persistent behavior survives refreshes
- Loading, empty and error states are present
- Usable on mobile and with keyboard navigation
- Important logic is covered by automated tests
- Documentation states its data source and limitations honestly
