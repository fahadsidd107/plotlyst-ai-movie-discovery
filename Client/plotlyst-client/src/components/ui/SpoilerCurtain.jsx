export default function SpoilerCurtain({ label = 'Spoiler-sensitive material', reason, onReveal, children }) {
  if (children) return children;
  return <section className="spoiler-firewall" aria-live="polite">
    <span aria-hidden="true">✦</span>
    <div><strong>{label}</strong><p>{reason || 'Plotlyst concealed this section using your spoiler protection mode.'}</p></div>
    <button type="button" onClick={onReveal} aria-label={`Reveal ${label}`}>Reveal once</button>
  </section>;
}
