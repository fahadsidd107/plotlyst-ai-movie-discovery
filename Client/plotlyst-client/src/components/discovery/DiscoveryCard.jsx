import { Link } from 'react-router-dom';
import { imageUrl } from '../../api/catalog';

const itemRoute = (item, fallback = 'movie') => {
  const type = item.media_type || fallback;
  return type === 'person' ? `/people/${item.id}` : `/${type}/${item.id}`;
};

export default function DiscoveryCard({ item, meta, type = 'movie', rank, wide = false }) {
  const title = item.title || item.name;
  const year = (item.release_date || item.first_air_date || '').slice(0, 4);
  const path = item.media_type === 'person' ? item.profile_path : (wide ? item.backdrop_path || item.poster_path : item.poster_path);
  return <Link className={`programme-card${wide ? ' programme-card--wide' : ''}`} to={itemRoute(item, type)}>
    <div className="programme-card__art">
      {path ? <img src={imageUrl(meta, path, wide ? 'w780' : 'w500')} alt={`${title} artwork`} loading="lazy" /> : <span aria-hidden="true">{title?.[0]}</span>}
      {rank && <b className="programme-card__rank">{String(rank).padStart(3, '0')}</b>}
      <i aria-hidden="true">{item.media_type === 'tv' || type === 'tv' ? 'SERIES' : item.media_type === 'person' ? 'PERSON' : 'FEATURE'}</i>
    </div>
    <h3>{title}</h3>
    <p>{year || item.known_for_department || 'Now programming'}{item.vote_average > 0 ? ` · ★ ${item.vote_average.toFixed(1)}` : ''}</p>
  </Link>;
}
