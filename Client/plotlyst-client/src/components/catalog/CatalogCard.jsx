import { Link } from 'react-router-dom';
import { imageUrl } from '../../api/catalog';

const CatalogCard = ({ item, meta, type = 'movie', profile = false }) => {
  const title = item.title || item.name;
  const path = profile ? item.profile_path : item.poster_path;
  const year = (item.release_date || item.first_air_date || '').slice(0, 4);
  return <Link className={`catalog-card${profile ? ' catalog-card--person' : ''}`} to={`/${type}/${item.id}`}>
    <div className="catalog-card__image">
      {path ? <img src={imageUrl(meta, path)} alt={`${title} ${profile ? 'portrait' : 'poster'}`} loading="lazy" /> : <span aria-hidden="true">{title?.slice(0, 1)}</span>}
      {item.vote_average > 0 && <b>{item.vote_average.toFixed(1)}</b>}
    </div>
    <h3>{title}</h3>
    <p>{profile ? item.known_for_department : year || 'Coming soon'}</p>
  </Link>;
};
export default CatalogCard;
