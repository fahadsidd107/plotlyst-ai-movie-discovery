export default function StarRating({ value = 0, onChange, disabled = false, label = 'Your rating' }) {
  return <fieldset className="star-rating" disabled={disabled}>
    <legend>{label}</legend>
    <div>{[1, 2, 3, 4, 5].map(star => <button key={star} type="button" aria-label={`${star} star${star > 1 ? 's' : ''}`} aria-pressed={value === star} onClick={() => onChange?.(star)}><span aria-hidden="true">★</span></button>)}</div>
    <output>{value ? `${value}.0 / 5` : 'Not rated'}</output>
  </fieldset>;
}
