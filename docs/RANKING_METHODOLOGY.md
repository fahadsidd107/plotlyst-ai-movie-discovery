# Plotlyst chart methodology

Plotlyst charts are discovery aids built from TMDB's live `vote_average` and `vote_count` fields. They are not IMDb rankings, awards, editorial verdicts, or a claim of objective quality.

## Eligibility

- Media type is fixed per chart: movie or television.
- Adult titles are excluded.
- A title requires at least 500 TMDB votes.
- Chart sizes are exactly 100, 250, or 2,000 entries, exposed as 5, 13, or 100 pages of at most 20 titles.
- Provider data is cached briefly by the Go service; the API response labels the result as a live cached provider snapshot.

## Score

Version `plotlyst-weighted-v1` uses a Bayesian weighted rating:

```text
score = (v / (v + m)) × R + (m / (v + m)) × C
```

- `R`: TMDB vote average for the title
- `v`: TMDB vote count for the title
- `m`: minimum eligible vote count (`500`)
- `C`: published global prior (`6.5`)

This prevents a title with only a handful of perfect votes from outranking a title supported by a large audience.

## Stable ordering

Each returned page is sorted using these deterministic keys:

1. weighted score, descending
2. vote count, descending
3. TMDB identifier, ascending

Rank numbers are derived from the requested page after this ordering. The response includes the formula, thresholds, tie-break policy, source and snapshot description so clients do not have to infer the methodology.

## Reproducibility boundary

TMDB data changes as votes arrive. Repeating a request against the same cached provider payload is deterministic; requests made after the cache expires can legitimately change. Immutable historical/nightly snapshots require a scheduled persistence pipeline and are intentionally not claimed by this live chart version.

## Other discovery labels

- **Trending** uses TMDB's shorter daily/weekly trending window.
- **Popular** uses TMDB popularity, which is distinct from rating quality.
- **Fan favourites** uses high audience ratings and engagement and remains distinct from the Plotlyst weighted chart.

See [TMDB's popularity and trending explanation](https://developer.themoviedb.org/docs/popularity-and-trending) and [Plotlyst's data policy](TMDB_DATA_POLICY.md).
