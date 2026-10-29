# Story DNA data and classification policy

Story DNA is a discovery aid, not an objective judgment about a film. Plotlyst stores a compact, versioned profile containing narrative labels, bounded 0–100 characteristic scores, provenance, confidence and timestamps in `story_dna_profiles`.

## Current classifier

The current `deterministic-metadata-heuristic-v1` classifier uses permitted structured TMDB fields: genre identifiers, runtime, release year and aggregate vote presence. It does not send TMDB descriptions, credits, images or other catalogue payloads to OpenAI. Confidence is deliberately capped below 0.70 because these signals cannot prove narrative form.

The API and interface therefore describe profiles as “inferred” and expose the method, source, policy, schema version and confidence. Consumers must not present these values as verified artistic facts.

## Spoilers

Ending shape and resolution are spoiler-sensitive. They are omitted by default and currently remain unknown because the metadata heuristic does not invent them. A future licensed or administrator-reviewed classifier may supply those fields; clients must still require an explicit spoiler reveal.

## Similarity

Similarity is deterministic. Mood compares warmth, tension, darkness and hope; pacing and complexity compare numeric distance; structure uses shared labels. User-controlled weights are normalized server-side and stable ties use the TMDB movie identifier. Explanations cite only these compared characteristics.

## Future AI seam

LangChainGo/OpenAI may classify only material Plotlyst is authorized to process. Generated JSON must be schema-validated, bounded, provenance-labelled and reviewed or confidence-gated before persistence. Adding an API key alone does not authorize sending TMDB content to a model.
