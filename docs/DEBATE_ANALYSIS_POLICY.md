# Plotlyst Debate analysis policy

Plotlyst Debate summarizes only genuine, visible Plotlyst community reviews. It never substitutes TMDB reviews, demo copy or invented sentiment when the community sample is too small.

## Minimum sample and ratings

Analysis requires at least five visible, top-level reviews for a movie. Reviews are joined server-side to their author’s Plotlyst rating. Rating distributions are calculated with a MongoDB aggregation over the authoritative `ratings` collection. High means 4–5 stars, middle means 3 stars and low means 1–2 stars.

If the minimum is not met, the API returns `status: insufficient`, the actual review count and the required count. `loved` and `disliked` stay empty.

## Grounding and evidence

Every observation contains real review identifiers and excerpts. The deterministic analyzer recognizes a bounded vocabulary for acting, cinematography, pacing, writing, soundtrack and direction, then emits a templated summary with a real mention count. It does not infer facts outside those reviews.

Spoiler-tagged review prose is excluded from theme extraction and generated summaries. Evidence links point back to the original community review. Generated summaries and user-authored excerpts are labelled separately in the interface.

## AI-assisted analysis

`DEBATE_AI_ENABLED=false` is the safe default. When explicitly enabled, LangChainGo/OpenAI receives only non-spoiler reviews belonging to members who enabled AI and activity-based explanations in their privacy settings. The prompt treats review text as untrusted data, output is schema-constrained, and every returned review ID is validated against the submitted set. Provider failure falls back to deterministic analysis.

Enabling OpenAI requires an operator to confirm that the deployment’s user terms, privacy notice, retention settings and provider agreement authorize this processing.

## Caching and refresh

Summaries are cached in `debate_summaries`. One or two new reviews are reported as pending without pretending they were analyzed. Three new reviews, review edits or deletions trigger a deduplicated background refresh; the existing result is returned as `stale` until that refresh completes.

## Polarization

Polarization measures the coexistence of meaningful high- and low-rating camps. It is explicitly described as disagreement, not film quality.
