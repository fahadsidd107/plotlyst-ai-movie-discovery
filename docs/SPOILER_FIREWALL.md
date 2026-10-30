# Spoiler Firewall policy

Plotlyst applies spoiler protection on the server before spoiler-sensitive prose reaches the browser. CSS blur is presentation, not the security boundary.

## Modes

- **OFF** returns tagged reviews, catalogue synopses and episode descriptions normally.
- **STANDARD** conceals user-tagged reviews and television episode descriptions beyond the member's furthest recorded watched episode. It is also the anonymous default.
- **STRICT** additionally conceals unreviewed community prose, movie/television detail synopses, discovery previews, recommendation summaries and AI Desk overview text.

Preferences are stored by user in `spoiler_preferences`. Invalid or missing values fail safely to STANDARD. Episode visibility is calculated from server-owned `episode_progress`; the client cannot submit an arbitrary progress boundary while browsing.

## Deliberate reveal

Concealed reviews are returned without their title or body. A separate `GET /api/community/reviews/:reviewId/reveal` request is required to reveal one item and is marked `Cache-Control: no-store`. Catalogue and episode reveal requests use an explicit one-time query flag initiated by an accessible button. Revealing one item does not silently change the account preference.

## Authoring and moderation

Review and reply authors can mark content as a spoiler. User tags are stored with source `USER` and confidence `1`. Moderators retain the existing report and removal workflow. Automated AI spoiler detection is not enabled: Plotlyst will not transmit reviews or third-party catalogue text to a model without an authorized data policy, explicit retention rules, conservative confidence thresholds and a moderation-review queue.

## Limits

No automated system can guarantee that every untagged sentence is spoiler-free. Strict mode therefore takes the conservative approach of concealing all community prose that has not been deliberately revealed. Titles, cast names, artwork and release metadata are not currently treated as spoilers.
