export const MovieGridSkeleton = ({ count = 8 }) => (
  <div className="movie-grid" aria-label="Loading movies" aria-busy="true">
    {Array.from({ length: count }, (_, index) => (
      <div className="skeleton-card" key={index} aria-hidden="true">
        <div className="skeleton skeleton--poster" />
        <div className="skeleton-card__body">
          <div className="skeleton skeleton--title" />
          <div className="skeleton skeleton--line" />
        </div>
      </div>
    ))}
    <span className="visually-hidden">Loading movies…</span>
  </div>
)

const StatePanel = ({ tone, eyebrow, title, message, actionLabel, onAction }) => (
  <section className={`state-panel state-panel--${tone}`} role={tone === 'error' ? 'alert' : 'status'}>
    <span className="state-panel__mark" aria-hidden="true">{tone === 'error' ? '!' : '○'}</span>
    <div>
      <p className="eyebrow">{eyebrow}</p>
      <h2>{title}</h2>
      <p>{message}</p>
      {actionLabel && onAction && <button className="button button--secondary" onClick={onAction}>{actionLabel}</button>}
    </div>
  </section>
)

export const EmptyState = ({ message = 'There is nothing to show here yet.', actionLabel, onAction }) => (
  <StatePanel tone="empty" eyebrow="Quiet scene" title="No titles found" message={message} actionLabel={actionLabel} onAction={onAction} />
)

export const ErrorState = ({ message = 'We could not load this page.', onRetry }) => (
  <StatePanel tone="error" eyebrow="Something went off-script" title="Unable to load titles" message={message} actionLabel={onRetry ? 'Try again' : undefined} onAction={onRetry} />
)
