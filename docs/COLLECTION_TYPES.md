# Collection ownership boundaries

Plotlyst deliberately keeps three concepts separate:

| Type | Owner | Source | Current behavior |
|---|---|---|---|
| Official franchise | TMDB | `/collection/{id}` | Read-only canonical film collection |
| Editorial collection | Plotlyst editors | Future dedicated MongoDB collection | Public curated programming; never labelled an official franchise |
| User list | Individual Plotlyst user | Future owner-scoped MongoDB collection | Read-only explanation in Step 3; creation and permissions arrive in Step 4 |

The API exposes these definitions at `GET /api/discovery/collection-kinds`. This prevents external franchise membership, staff curation and user ownership from sharing ambiguous authorization rules.
