# Plotlyst 🎬

> **Discover beyond the screen.**

Plotlyst is an evolving, full-stack movie discovery and recommendation platform. It combines a React experience with a Go/Gin API, MongoDB persistence, secure cookie-based authentication, preference-driven recommendations, and optional AI-assisted review analysis through LangChainGo and OpenAI.

The current release combines the original MongoDB-backed editorial catalogue with a server-side TMDB adapter for real movies, television, seasons, episodes and people. Users can browse media, inspect credits, regional release/provider information and recommendations, create accounts, receive genre-based recommendations, and—when authorized—write persisted editorial reviews.

## Why this project exists

Plotlyst is being engineered as a portfolio-grade product that demonstrates how a modern entertainment platform can grow from a focused modular monolith into a data-rich discovery system without discarding working functionality.

Its long-term direction includes real movie and television metadata, rich search, ranked charts, ratings and reviews, personal collections, streaming-provider discovery, explainable AI recommendations, and collaborative movie-night experiences. Planned work is documented separately and is never presented as completed functionality.

## Current capabilities

- Movie catalogue backed by MongoDB
- Registration with bcrypt password hashing
- Login, logout, access tokens, refresh tokens, and HTTP-only cookies
- Role-aware protected routes for users and administrators
- Favourite-genre selection during registration
- Personalized recommendations derived from saved genre preferences
- Administrator review workflow with persisted ranking updates
- Optional LangChainGo/OpenAI review classification
- Offline sentiment fallback when no AI credentials are configured
- CORS configuration through environment variables
- Startup-time environment validation with actionable failures
- Structured, machine-readable API errors
- Separate liveness and database-readiness health checks
- Graceful HTTP shutdown and MongoDB disconnect handling
- Idempotent local seed-data bootstrap
- Single-command local development for the database, API, and frontend
- Owner-controlled profiles, privacy preferences and regional settings
- Five-star audience ratings with aggregate distributions and programme-note review threads
- Private watchlists, custom playlists, watched history and episode progress
- USER, MODERATOR and ADMIN authorization with a dedicated moderation queue
- Traceable specialist-data adapters for regional streaming, verified awards, licensed US box office and editorial classifications
- Responsive React UI with loading and error states
- TMDB-backed movie and television landing/detail experiences
- Seasons, episodes, people profiles, trending people and date-aware Born Today
- Cast, crew, franchises, certifications, recommendations and German provider availability
- Bounded TMDB retries, `Retry-After` handling and short-lived in-memory caching
- Idempotent legacy editorial-ranking to five-star migration
- Dedicated `/test` loading-screen laboratory
- Progressively loaded cinematic homepage with independent programme reels
- Accessible movie, television and people autocomplete plus full search results
- Backend-driven discovery filters for genre, year, decade, country, language, runtime and provider
- Movie and television Top 100, Top 250 and paginated Top 2,000 charts
- Published weighted-ranking formula, eligibility threshold and deterministic tie-breaking
- Trending, fan-favourite, year, decade and genre destinations
- Explicit official-franchise, editorial-collection and user-list ownership boundaries
- Cinema-first Midnight Premiere interface with projector, film-strip, ticket and screening-programme details
- Typed natural-language and movie-memory discovery with AI/local intent provenance
- Server-side retrieval, deterministic ranking and explanations that keep TMDB payloads outside model prompts
- Cinema DNA analytics generated from real MongoDB activity with explicit coverage reporting
- Persistent Group Movie Night rooms with authenticated live updates and deterministic preference intersection
- Licensed OpusData box-office seam with a conspicuous synthetic demo fallback for local development
- Public AI Desk discovery with authenticated personalization kept behind privacy controls
- Structured JSON request logs, correlation IDs, Prometheus metrics and operational health probes
- Hardened cookie policy, security headers, request-size limits and rate-limited write surfaces
- OpenAPI documentation with automated contract validation
- Go race, provider-mock, repository/integration and frontend interaction test suites
- Route-level React lazy loading plus TMDB request coalescing and bounded caching
- Production Docker images, Docker Compose, CI quality gates, Render and Vercel configuration
- Story DNA narrative fingerprints with versioned provenance, bounded confidence and spoiler-safe traits
- Adjustable, deterministic story-similarity discovery with grounded comparison explanations
- Persistent OFF, STANDARD and STRICT spoiler protection enforced by the Go API
- Progress-aware television descriptions and deliberate, accessible one-item spoiler reveals

## Architecture

```text
Browser
  │
  ▼
React + Vite client (:5173)
  │  Axios / JSON / secure cookies
  ▼
Go + Gin API (:8080)
  ├── Authentication middleware
  ├── Request IDs, structured logs, metrics and security middleware
  ├── Movie, television, people and user controllers
  ├── TMDB client → timeout / retry / rate-limit / coalesced cache
  ├── Discovery orchestration → homepage / search / filters / charts
  ├── Intelligence pipeline → intent / validation / retrieval / ranking / explanation
  ├── Specialist providers → streaming / awards / box office / editorial
  ├── Recommendation and review logic
  ├── LangChainGo → OpenAI (optional)
  └── MongoDB driver
          │
          ▼
     MongoDB (:27017)
```

The backend currently follows a compact layered structure:

```text
Server/PlotlystServer/
├── controllers/    # HTTP handlers and current application workflows
├── intelligence/   # Typed AI/local intent, retrieval, ranking and explanations
├── specialist/     # Traceable availability, awards and box-office providers
├── contracts/      # OpenAPI contract validation
├── database/       # MongoDB connection and collection access
├── middleware/     # Authentication and request authorization
├── models/         # BSON/JSON domain models
├── tmdb/           # Typed provider DTOs and resilient HTTP client
├── migration/      # Idempotent rating compatibility rules
├── cmd/             # Runnable maintenance commands
├── routes/         # Public and protected route registration
├── utils/          # JWT creation, validation, and token persistence
└── main.go         # Configuration, CORS, database, and server bootstrap
```

As the domain grows, business rules will move behind handler → service → repository boundaries while retaining the Go modular-monolith architecture.

## Technology stack

| Area | Current implementation | Planned evolution |
|---|---|---|
| Frontend | React 19, Vite, React Router, Axios, lazy route chunks | Ongoing real-device polish |
| UI | Midnight Pop cinema design system, Bootstrap, React Bootstrap | Ongoing accessibility evaluation |
| Backend | Go, Gin, context-aware MongoDB driver, resilient TMDB client | Further domain repository extraction |
| Database | MongoDB with indexed community, activity, live-room and cache collections | Production capacity tuning |
| Authentication | bcrypt, JWT rotation, hardened HTTP-only cookies, USER/MODERATOR/ADMIN | External penetration review |
| AI | Typed OpenAI intent extraction or deterministic local parser; server-side retrieval/ranking | Activity metadata enrichment and model/provider evaluations |
| Live features | Authenticated Movie Night rooms and WebSocket snapshots | Candidate ballots, provider intersection and notifications |
| Catalogue data | TMDB plus compatible curated local fixtures | Search, charts and licensed specialist providers |
| Quality | Go race/unit/integration/provider tests, Vitest interactions, ESLint, OpenAPI contracts | Deployed Lighthouse and physical-device audits |
| Delivery | One-command local environment, Docker Compose, CI, Vercel and Render manifests | Publish to account-owned production hosts |

## Current API

| Method | Endpoint | Access | Purpose |
|---|---|---|---|
| `GET` | `/hello` | Public | API health smoke test |
| `GET` | `/health/live` | Public | Confirm the API process is alive |
| `GET` | `/health/ready` | Public | Confirm MongoDB-backed readiness |
| `GET` | `/metrics` | Public/operations | Prometheus-compatible request metrics |
| `GET` | `/movies` | Public | List catalogue movies |
| `GET` | `/genres` | Public | List registration/discovery genres |
| `POST` | `/register` | Public | Create a user with a hashed password |
| `POST` | `/login` | Public | Authenticate and issue secure cookies |
| `POST` | `/logout` | Public | Revoke stored tokens and clear cookies |
| `POST` | `/refresh` | Refresh cookie | Rotate access and refresh tokens |
| `GET` | `/movie/:imdb_id` | Authenticated | Retrieve one movie |
| `GET` | `/recommendedmovies` | Authenticated | Recommend films by favourite genres |
| `POST` | `/addmovie` | Authenticated | Add a validated movie document |
| `PATCH` | `/updatereview/:imdb_id` | Administrator | Update and classify an editorial review |
| `GET` | `/api/catalog/movies` | Public | Popular or trending TMDB movies |
| `GET` | `/api/catalog/movies/:id` | Public | Movie, credits, releases, providers and recommendations |
| `GET` | `/api/catalog/tv` | Public | Popular or trending TMDB television |
| `GET` | `/api/catalog/tv/:id` | Public | Series, seasons, credits, providers and recommendations |
| `GET` | `/api/catalog/tv/:id/seasons/:season` | Public | Season and episode summaries |
| `GET` | `/api/catalog/tv/:id/seasons/:season/episodes/:episode` | Public | Complete episode details |
| `GET` | `/api/catalog/people/trending` | Public | Trending people |
| `GET` | `/api/catalog/people/born-today` | Public | Timezone-aware birthday matches from the documented sample |
| `GET` | `/api/catalog/people/:id` | Public | Person biography and combined credits |
| `GET` | `/api/discovery/home` | Public | Concurrent aggregated homepage programme |
| `GET` | `/api/discovery/home/:section` | Public | Independently load one progressive homepage rail |
| `GET` | `/api/discovery/search` | Public | Normalized movie, TV and people search |
| `GET` | `/api/discovery/movies` | Public | Filtered movie discovery |
| `GET` | `/api/discovery/tv` | Public | Filtered television discovery |
| `GET` | `/api/discovery/filters/:media` | Public | Provider-backed filter options |
| `GET` | `/api/discovery/charts/:media/:size` | Public | Top 100, 250 or 2,000 chart page |
| `GET` | `/api/discovery/franchises/:id` | Public | Official TMDB franchise |
| `GET` | `/api/discovery/collection-kinds` | Public | Collection ownership/source definitions |
| `POST` | `/api/ai/discover` | Public; optional session | Typed natural-language discovery with optional personalization |
| `POST` | `/api/ai/memory-search` | Public; optional session | Search from a remembered scene or description |
| `GET` | `/api/story-dna/movies/:id` | Public | Inferred narrative profile with confidence and provenance |
| `POST` | `/api/story-dna/discover` | Public | Weighted, explainable story-similarity discovery |
| `GET/PUT` | `/api/me/spoiler-preferences` | Authenticated | Read or update OFF, STANDARD or STRICT protection |
| `GET` | `/api/community/reviews/:reviewId/reveal` | Public; deliberate action | Reveal one protected review without changing preference |
| `GET` | `/api/analytics/cinema-dna` | Authenticated | Aggregate personal activity analytics |
| `POST` | `/api/movie-night/rooms` | Authenticated | Create a collaborative room |
| `POST` | `/api/movie-night/rooms/:code/join` | Authenticated | Join a room by code |
| `PUT` | `/api/movie-night/rooms/:code/preferences` | Participant | Update the shared preference ballot |
| `GET` | `/api/movie-night/rooms/:code/ws` | Participant | Open the authenticated live room channel |

## Run locally

### Prerequisites

- macOS with Homebrew
- Go 1.24 or newer
- Node.js 22.22.2 or newer
- MongoDB Community 8
- MongoDB Shell and Database Tools

The repository includes a local orchestrator that starts MongoDB, imports collections only when empty, installs missing frontend dependencies, and runs both applications.

```bash
git clone <your-repository-url>
cd <repository-folder>
npm run dev
```

Open [http://localhost:5173](http://localhost:5173). Press `Ctrl+C` to stop both development servers.

### Local administrator fixture

```text
Email:    fahd@gmail.com
Password: Password1!
```

This account is for local development only. Production deployments must provision credentials securely and must never expose demonstration passwords.

## Environment configuration

Backend file: `Server/PlotlystServer/.env`

Copy the safe templates before adding local values:

```bash
cp Server/PlotlystServer/.env.example Server/PlotlystServer/.env
cp Client/plotlyst-client/.env.example Client/plotlyst-client/.env
```

```dotenv
MONGODB_URI=mongodb://127.0.0.1:27017
DATABASE_NAME=plotlyst
PORT=8080
ALLOWED_ORIGINS=http://localhost:5173
SECRET_KEY=<random-64-character-secret>
SECRET_REFRESH_KEY=<different-random-64-character-secret>
RECOMMENDED_MOVIE_LIMIT=5
OPENAI_API_KEY=
BASE_PROMPT_TEMPLATE="Classify the review using one ranking from: {rankings}."
TMDB_API_KEY=<your-v3-api-key>
TMDB_BASE_URL=https://api.themoviedb.org/3
TMDB_IMAGE_BASE_URL=https://image.tmdb.org/t/p
TMDB_TIMEOUT=8s
TMDB_CACHE_TTL=10m
TMDB_MAX_RETRIES=3
AI_PROVIDER=offline
AI_MODEL=gpt-4.1-mini
AI_TIMEOUT=8s
AI_MAX_INPUT_CHARS=800
AI_MAX_OUTPUT_TOKENS=350
AI_DAILY_USER_BUDGET=20
AI_CACHE_TTL=30m
OPUSDATA_API_KEY=
OPUSDATA_BASE_URL=
OPUSDATA_TIMEOUT=8s
BOX_OFFICE_DEMO_ENABLED=true
```

Frontend file: `Client/plotlyst-client/.env`

```dotenv
VITE_API_BASE_URL=http://localhost:8080
```

Generate JWT secrets with `openssl rand -hex 32`. Environment files are ignored by Git; never commit real secrets.

The browser never receives `TMDB_API_KEY`. If TMDB returns `Invalid API key`, create or activate a v3 key in your TMDB account and replace only the value in `Server/PlotlystServer/.env`.

`OPENAI_API_KEY` is optional. Without it, review updates use the deterministic offline classifier. With it, the backend uses LangChainGo to make a real model request and validates the returned ranking before persistence.

For Step 6 discovery, `AI_PROVIDER=offline` keeps intent extraction local. Set it to `openai` and provide `OPENAI_API_KEY` to enable typed model extraction; retrieval and ranking remain server-side, and the response always identifies which path ran.

Box-office pages use explicitly labelled synthetic demo fixtures until both OpusData settings are supplied. The licensed API endpoint and response mapping must match your OpusData contract; Plotlyst does not guess or scrape The Numbers. A configured provider failure is surfaced instead of silently falling back to demo values.

## Data model today

- `users`: identity, role, hashed password, favourite genres, stored token state and timestamps
- `movies`: IMDb identifier, title, poster, genres, editorial review and ranking
- `genres`: normalized genre identifiers and names
- `rankings`: legacy editorial ranking values (`Excellent` through `Terrible`)

Seed fixtures live in `plotlyst-seed-data/`. The development script imports a collection only when that collection is empty, preventing duplicate fixtures on ordinary restarts.

## Engineering principles

- Real backend integration over hardcoded interface demos
- Persistent user actions and database-enforced ownership
- Explicit separation between source ratings and future Plotlyst community ratings
- Server-side ownership of third-party credentials
- Honest empty/error states when a provider is unavailable
- No fabricated awards, reviews, rankings, availability, or box-office figures
- Structured AI output, validated filters, cost controls and non-AI fallbacks
- Respect for third-party licences, attribution requirements and data-use limitations

## Product roadmap

The full implementation is divided into auditable vertical slices so every feature includes its backend, frontend, persistence, loading/error states, authorization, and tests. See [docs/ROADMAP.md](docs/ROADMAP.md) for the detailed plan.

High-level sequence:

1. Foundation, design system, configuration and API contracts
2. TMDB-backed movies, television, people and media details
3. Discovery homepage, search, charts, genres and collections
4. Ratings, reviews, watchlists, playlists, profiles and episode progress
5. Awards, regional availability, box office and editorial data adapters
6. AI discovery, explainable recommendations, analytics and live collaboration
7. Security, testing, performance, observability and deployment

## Responsible external integrations

The TMDB catalogue adapter keeps provider credentials on the Go server and uses timeouts, bounded retries, rate-limit handling and a short in-memory cache. TMDB content is not sent to AI providers. See [`docs/TMDB_DATA_POLICY.md`](docs/TMDB_DATA_POLICY.md) for attribution, provider and licence details.

Chart eligibility, formula, stable tie-breaking and live-snapshot limits are documented in [`docs/RANKING_METHODOLOGY.md`](docs/RANKING_METHODOLOGY.md). Franchise, editorial and user-list ownership boundaries are documented in [`docs/COLLECTION_TYPES.md`](docs/COLLECTION_TYPES.md).

Story DNA classification, confidence, spoilers, deterministic weighting and the future authorized-AI seam are documented in [`docs/STORY_DNA_POLICY.md`](docs/STORY_DNA_POLICY.md).

Server-side review redaction, progress-aware episode visibility, protection modes and the conservative AI-detection boundary are documented in [`docs/SPOILER_FIREWALL.md`](docs/SPOILER_FIREWALL.md).

## Build and validation

```bash
# Frontend
npm --prefix Client/plotlyst-client run lint
npm --prefix Client/plotlyst-client test
npm --prefix Client/plotlyst-client run build

# Backend
cd Server/PlotlystServer
go test ./...
go test -race ./...
go vet ./...
```

The HTTP contract is published as [`docs/openapi.yaml`](docs/openapi.yaml), and response conventions are documented in [`docs/API_CONTRACTS.md`](docs/API_CONTRACTS.md). Production setup and the honest remaining audit boundary are described in [`docs/DEPLOYMENT.md`](docs/DEPLOYMENT.md) and [`docs/PRODUCTION_AUDIT.md`](docs/PRODUCTION_AUDIT.md).

Run the idempotent legacy editorial-rating compatibility migration with `npm run migrate:ratings`.

## Status

Plotlyst is feature-complete through the local production-readiness baseline. The application, tests, containers, CI and deployment manifests are ready; publishing the frontend and persistent API remains intentionally open until account-owned Vercel and Go-hosting credentials are supplied.

---

Built as a hands-on full-stack engineering project focused on Go backend architecture, React product development, MongoDB data modelling, secure authentication, and responsible AI integration.
