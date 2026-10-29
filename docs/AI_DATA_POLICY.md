# Plotlyst AI data policy

Plotlyst separates intent extraction, validated filters, catalogue retrieval, deterministic ranking and grounded explanations.

- Model providers receive only the user's bounded discovery sentence. TMDB catalogue payloads, reviews and private activity are not placed in model prompts.
- `AI_PROVIDER=offline` uses the deterministic parser. The response identifies it as `deterministic`, never as AI.
- `AI_PROVIDER=openai` uses structured JSON output that is revalidated against server allowlists. If extraction fails, Plotlyst reports that the local parser was used.
- Catalogue retrieval happens server-side through the resilient TMDB client. Ranking is stable and deterministic; identical scores use the catalogue ID as a tie-breaker.
- Explanations use fixed grounded templates. They do not invent viewing history, taste or catalogue facts.
- Queries are control-character stripped, length-bounded, delimited as untrusted input and cannot invoke tools or URLs.
- Requests have a timeout and per-user rate budget. Raw prompts are not stored by the Step 6 pipeline.

The current Cinema DNA report aggregates only persisted Plotlyst activity. It returns a coverage notice because legacy activity rows do not yet contain complete genre, language and runtime snapshots.
