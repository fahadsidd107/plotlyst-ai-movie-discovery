# Group Movie Night

Movie Night is the only Step 6 feature with a live channel because participants need shared room updates. REST remains the source of truth; the authenticated WebSocket sends room snapshots and preference intersections.

- Rooms persist in MongoDB and expire after 24 hours through a TTL index.
- Only room participants may read the room or open its live channel.
- Preference intersection is deterministic: shared genres survive only when every participant selects them, and the strictest positive runtime limit wins.
- No private reviews or unrelated account activity are broadcast.
- A disconnected browser can continue through REST and reconnect to the live channel.

The first release intentionally keeps voting and provider intersection small. The schema and routes are designed so ballots and candidates can be added without changing room identity or membership rules.
