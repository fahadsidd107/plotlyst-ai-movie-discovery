# Plotlyst deployment runbook

Plotlyst is deployment-ready but deployment requires credentials owned by the operator. Never place production secrets in Git, Vercel build logs or a client-side `VITE_*` variable.

## 1. Production database

Create a persistent MongoDB Atlas cluster (or another durable MongoDB 8 service), a least-privilege application user and an IP/network rule for the API host. Save the connection string as `MONGODB_URI`. Do not use the Docker Compose database for public production.

## 2. API on persistent Go hosting

`render.yaml` describes a Docker-based Render service with `/health/ready` health checks. Connect the repository through Render Blueprint and provide the unsynced values shown in the file. Equivalent container hosts such as Fly.io, Railway or Google Cloud Run may use `Server/PlotlystServer/Dockerfile`.

Required production values:

- `MONGODB_URI`, `DATABASE_NAME`
- `ALLOWED_ORIGINS=https://your-vercel-domain`
- two independently generated JWT secrets
- `TMDB_API_KEY`

Optional provider values:

- `OPENAI_API_KEY` with `AI_PROVIDER=openai`
- licensed `OPUSDATA_API_KEY` and contracted `OPUSDATA_BASE_URL`

Set `BOX_OFFICE_DEMO_ENABLED=false` in production unless the product intentionally demonstrates the clearly labelled fixture.

Verify after deployment:

```bash
curl -fsS https://your-api.example/health/live
curl -fsS https://your-api.example/health/ready
curl -fsS https://your-api.example/metrics
```

## 3. Client on Vercel

Import the repository in Vercel and choose `Client/plotlyst-client` as the root directory. The committed `vercel.json` supplies Vite and SPA routing configuration. Set:

```dotenv
VITE_API_BASE_URL=https://your-api.example
```

Deploy, copy the final Vercel URL into the API's `ALLOWED_ORIGINS`, and redeploy the API. Confirm login cookies over HTTPS, direct navigation to nested routes and WebSocket connectivity for Movie Night.

## 4. Container smoke test

For a production-like local build:

```bash
cp Server/PlotlystServer/.env.example Server/PlotlystServer/.env
docker compose up --build
```

Open `http://localhost:5173`, then check the API health and metrics endpoints on port 8080.

## 5. Rollback and operations

- Keep at least one previously healthy container image or host deployment.
- Roll back application code before applying destructive database changes. Current startup indexes are additive.
- Monitor 5xx counts and cumulative latency from `/metrics`; correlate failures using `X-Request-ID` and structured JSON logs.
- Rotate JWT/provider credentials if logs or hosting access are compromised.
- Database backups and restore exercises belong to the MongoDB host; verify them before public launch.

## Account-dependent boundary

Repository code cannot create or bill external hosting accounts on its own. The Vercel/Render deployment is complete only after an authorized owner connects the Git repository, supplies secrets and confirms the resulting public URLs.
