import { Link } from 'react-router-dom';
import DiscoveryCard from './DiscoveryCard';
import { ErrorState } from '../ui/States';

export default function ProgrammeRail({ number, title, copy, items, meta, type, loading, error, destination }) {
  return <section className="programme-reel" aria-labelledby={`reel-${number}`}>
    <header className="programme-reel__header">
      <span className="reel-number" aria-hidden="true">{String(number).padStart(2, '0')}</span>
      <div><p className="eyebrow">Programme reel</p><h2 id={`reel-${number}`}>{title}</h2>{copy && <p>{copy}</p>}</div>
      {destination && <Link className="programme-link" to={destination}>View programme <span>→</span></Link>}
    </header>
    {loading ? <div className="programme-rail" aria-label="Loading titles">{[1,2,3,4,5].map(i => <div className="programme-placeholder" key={i}/>)}</div>
      : error ? <ErrorState message={error} />
      : items?.length ? <div className="programme-rail">{items.slice(0, 12).map(item => <DiscoveryCard key={item.id} item={item} meta={meta} type={type}/>)}</div>
      : <div className="reel-empty"><b>No titles on this reel yet.</b><span>The projectionist is preparing the next programme.</span></div>}
  </section>;
}
