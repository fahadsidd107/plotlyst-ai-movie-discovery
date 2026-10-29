import { Link, useParams } from 'react-router-dom';
import { imageUrl } from '../../api/catalog';
import useCatalog from '../../hooks/useCatalog';
import CatalogCard from './CatalogCard';
import LoadingScreen from '../ui/LoadingScreen';
import { ErrorState } from '../ui/States';
import CommunityPanel from '../community/CommunityPanel';
import StoryDnaPanel from '../storyDna/StoryDnaPanel';

const DetailPage = ({ type, loader }) => {
  const { id } = useParams();
  const { data: item, meta, loading, error, retry } = useCatalog(() => loader(id), [id]);
  if (loading) return <LoadingScreen compact />;
  if (error) return <div className="page-shell content-section"><ErrorState message={error} onRetry={retry} /></div>;
  const title = item.title || item.name;
  const cast = item.credits?.cast || item.cast || [];
  const crew = item.credits?.crew || item.crew || [];
  const recommendations = item.recommendations?.results || item.recommendations || [];
  const providers = item['watch/providers']?.results?.DE || item.watch_providers?.results?.DE || item.providers?.DE;
  const regionalRelease = type === 'movie'
    ? ['DE', 'US'].map(country => item.release_dates?.results?.find(entry => entry.iso_3166_1 === country)).find(Boolean)
    : null;
  const releaseEntries = regionalRelease?.release_dates || [];
  const certification = releaseEntries.find(entry => entry.certification)?.certification;
  const regionalDate = releaseEntries.find(entry => entry.release_date)?.release_date;
  const contentRating = type === 'tv'
    ? ['DE', 'US'].map(country => item.content_ratings?.results?.find(entry => entry.iso_3166_1 === country)).find(Boolean)
    : null;
  return <article className="detail-page">
    <header className="detail-hero" style={item.backdrop_path ? { '--backdrop': `url(${imageUrl(meta, item.backdrop_path, 'w1280')})` } : undefined}>
      <div className="detail-hero__inner page-shell">
        {item.poster_path && <img src={imageUrl(meta, item.poster_path)} alt={`${title} poster`} />}
        <div><p className="eyebrow">{type === 'movie' ? 'Feature presentation' : 'On the small screen'}</p><h1>{title}</h1>
          <p className="detail-facts">{[(item.release_date || item.first_air_date || '').slice(0, 4), item.runtime ? `${item.runtime} min` : item.number_of_seasons ? `${item.number_of_seasons} seasons` : null, item.vote_average ? `★ ${item.vote_average.toFixed(1)}` : null].filter(Boolean).join(' · ')}</p>
          <div className="genre-row">{item.genres?.map(genre => <span key={genre.id}>{genre.name}</span>)}</div><p className="detail-overview">{item.overview}</p></div>
      </div>
    </header>
    <div className="page-shell detail-flow">
      {type === 'movie' && <StoryDnaPanel movieId={id} title={title} />}
      {item.belongs_to_collection && <section className="editorial-band"><p className="eyebrow">Part of a bigger story</p><h2>{item.belongs_to_collection.name}</h2></section>}
      {type === 'tv' && item.seasons?.length > 0 && <section><div className="section-heading"><div><p className="eyebrow">Episode guide</p><h2>Seasons</h2></div></div><div className="season-rail">{item.seasons.filter(s => s.season_number > 0).map(s => <Link key={s.id} to={`/tv/${id}/seasons/${s.season_number}`}><b>{String(s.season_number).padStart(2, '0')}</b><span>{s.name}</span><small>{s.episode_count} episodes</small></Link>)}</div></section>}
      {(regionalRelease || contentRating) && <section className="release-panel"><div><p className="eyebrow">Regional information</p><h2>{type === 'movie' ? 'Release & certification' : 'Content rating'}</h2></div><dl>
        {regionalRelease && <><div><dt>Region</dt><dd>{regionalRelease.iso_3166_1}</dd></div>{regionalDate && <div><dt>Release date</dt><dd>{new Intl.DateTimeFormat('en-GB', { dateStyle: 'long' }).format(new Date(regionalDate))}</dd></div>}{certification && <div><dt>Certification</dt><dd>{certification}</dd></div>}</>}
        {contentRating && <><div><dt>Region</dt><dd>{contentRating.iso_3166_1}</dd></div><div><dt>Rating</dt><dd>{contentRating.rating || 'Not rated'}</dd></div></>}
      </dl></section>}
      {cast.length > 0 && <section><div className="section-heading"><div><p className="eyebrow">Faces in focus</p><h2>Cast</h2></div></div><div className="people-rail">{cast.map(person => <CatalogCard key={`${person.id}-${person.character}`} item={person} meta={meta} type="people" profile />)}</div></section>}
      {(crew.length > 0 || providers) && <section className="detail-columns"><div><p className="eyebrow">Behind the frame</p><h2>Key crew</h2>{crew.slice(0, 8).map(p => <p key={`${p.id}-${p.job}`}><strong>{p.name}</strong><span>{p.job}</span></p>)}</div><div><p className="eyebrow">Watch in Germany</p><h2>Providers</h2><p>{[...(providers?.flatrate || []), ...(providers?.rent || [])].map(p => p.provider_name).filter((v, i, a) => a.indexOf(v) === i).join(' · ') || 'No regional provider data available.'}</p><small>Streaming availability data by JustWatch.</small></div></section>}
      {recommendations.length > 0 && <section><div className="section-heading"><div><p className="eyebrow">Keep the signal going</p><h2>You may also like</h2></div></div><div className="recommendations-rail">{recommendations.slice(0, 10).map(rec => <CatalogCard key={rec.id} item={rec} meta={meta} type={type} />)}</div></section>}
      <CommunityPanel mediaType={type} mediaId={id} title={title} posterPath={item.poster_path} />
    </div>
  </article>;
};
export default DetailPage;
