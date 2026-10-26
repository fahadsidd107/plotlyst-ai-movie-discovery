# Plotlyst 🎬

> **Discover beyond the screen.**

Plotlyst is an evolving, full-stack movie discovery and recommendation platform. It combines a React experience with a Go/Gin API, MongoDB persistence, secure cookie-based authentication, preference-driven recommendations, and optional AI-assisted review analysis through LangChainGo and OpenAI.

The current release is a working foundation rather than a collection of static mockups: users can browse seeded films, create accounts, authenticate, receive genre-based recommendations, and—when authorized as an administrator—write and classify reviews that persist in MongoDB.

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
- Responsive React UI with loading and error states

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
  ├── Movie and user controllers
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
├── database/       # MongoDB connection and collection access
├── middleware/     # Authentication and request authorization
├── models/         # BSON/JSON domain models
├── routes/         # Public and protected route registration
├── utils/          # JWT creation, validation, and token persistence
└── main.go         # Configuration, CORS, database, and server bootstrap
```

As the domain grows, business rules will move behind handler → service → repository boundaries while retaining the Go modular-monolith architecture.

## Technology stack

| Area | Current implementation | Planned evolution |
|---|---|---|
| Frontend | React 19, Vite, React Router, Axios | Design system, advanced discovery, accessibility and component tests |
| UI | Bootstrap, React Bootstrap, Font Awesome | Midnight Pop visual language and reusable tokens |
| Backend | Go, Gin, context-aware MongoDB driver | Versioned API, service/repository layers and OpenAPI |
| Database | MongoDB with BSON models and seed automation | Indexed ratings, reviews, lists, activity and cache collections |
| Authentication | bcrypt, JWT access/refresh tokens, HTTP-only cookies | Session hardening, moderation roles and account settings |
| AI | LangChainGo + optional OpenAI, offline fallback | Structured natural-language search and explainable recommendations |
| Live features | Not implemented yet | WebSocket-powered notifications and collaborative group sessions |
| Catalogue data | Curated local fixtures | Server-side TMDB integration with caching and attribution |
| Quality | Go build checks, ESLint/build scripts | Unit, integration, frontend and contract test suites |
| Delivery | One-command local environment | Docker, CI/CD, Vercel frontend and containerized API hosting |

## Current API

| Method | Endpoint | Access | Purpose |
|---|---|---|---|
| `GET` | `/hello` | Public | API health smoke test |
| `GET` | `/health/live` | Public | Confirm the API process is alive |
| `GET` | `/health/ready` | Public | Confirm MongoDB-backed readiness |
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

## Run locally

### Prerequisites

- macOS with Homebrew
- Go 1.24 or newer
- Node.js 20 or newer
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
```

Frontend file: `Client/plotlyst-client/.env`

```dotenv
VITE_API_BASE_URL=http://localhost:8080
```

Generate JWT secrets with `openssl rand -hex 32`. Environment files are ignored by Git; never commit real secrets.

`OPENAI_API_KEY` is optional. Without it, review updates use the deterministic offline classifier. With it, the backend uses LangChainGo to make a real model request and validates the returned ranking before persistence.

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

## Build and validation

```bash
# Frontend
npm --prefix Client/plotlyst-client run lint
npm --prefix Client/plotlyst-client run build

# Backend
cd Server/PlotlystServer
go test ./...
```

The current response conventions and structured error envelope are documented in [`docs/API_CONTRACTS.md`](docs/API_CONTRACTS.md).

Run the idempotent legacy editorial-rating compatibility migration with `npm run migrate:ratings`.

## Status

Plotlyst is under active development. The authenticated MongoDB-backed foundation and TMDB movie, television and people catalogue work today; later roadmap phases will be delivered incrementally.

---

Built as a hands-on full-stack engineering project focused on Go backend architecture, React product development, MongoDB data modelling, secure authentication, and responsible AI integration.
