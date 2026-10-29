# Production audit record

## Security

- Secrets remain server-only and `.env` files are ignored.
- Authentication uses HTTP-only cookies; CORS uses an explicit origin allowlist.
- Role and ownership middleware protect community and moderation writes.
- AI input is bounded and schema-validated; catalogue payloads are not sent to AI.
- Specialist demo data is visually and structurally marked as synthetic.
- API errors include a request ID without exposing internal errors.
- `npm audit --omit=dev` reports zero known production dependency vulnerabilities as of the audit date.
- Remaining launch work: configure HTTPS-only cookie policy for the final domain, run an external penetration test and rotate all development credentials.

## Accessibility

- Skip navigation, semantic landmarks, labelled forms, visible focus, autocomplete semantics and reduced-motion fallbacks are present.
- Loading messages use live status semantics while decorative film mechanics are hidden from assistive technology.
- Responsive navigation and search suggestion behavior have component tests.
- Remaining launch work: manual screen-reader coverage with VoiceOver/NVDA and user testing at 200%/400% zoom.

## Performance

- Route-level code splitting keeps specialist, analytics, profile and catalogue features out of the initial bundle.
- TMDB calls use bounded retries, TTL caching, a 500-entry cap and concurrent-request coalescing.
- Homepage rails fetch concurrently and fail independently.
- Motion relies primarily on transforms/opacity, avoids fixed canvases and honors reduced motion.
- Remaining launch work: run Lighthouse/WebPageTest against the deployed URLs with production network conditions and set budgets from observed field data.

## Responsive behavior

- Cinema rails preserve native horizontal scrolling with distinct thin scrollbars.
- Global overflow guards and 390px layout rules protect narrow screens.
- Desktop and mobile visual evaluation is required after every major design change.
- Remaining launch work: real-device checks on iOS Safari and Android Chrome after deployment.
