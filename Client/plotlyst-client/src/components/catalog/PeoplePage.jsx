import useCatalog from '../../hooks/useCatalog';
import { getBornToday, getPeople } from '../../api/catalog';
import CatalogCard from './CatalogCard';
import { ErrorState, MovieGridSkeleton } from '../ui/States';

const PeoplePage = () => {
  const born = useCatalog(getBornToday, []); const trending = useCatalog(() => getPeople(1), []);
  const bornItems = born.data?.results || born.data || []; const trendingItems = trending.data?.results || trending.data || [];
  return <div className="catalog-page"><header className="catalog-hero page-shell"><div><p className="eyebrow">People behind the pulse</p><h1>Meet the names lighting up every frame.</h1><p>Performers, makers, scene-stealers and the creative minds currently moving culture.</p></div><div className="signal-rail" aria-hidden="true"><i/><i/><i/><i/><i/></div></header>
    <section className="page-shell catalogue-section"><div className="section-heading"><div><p className="eyebrow">Born today</p><h2>Share a date with greatness</h2></div></div>{born.loading ? <MovieGridSkeleton count={5}/> : born.error ? <ErrorState message={born.error} onRetry={born.retry}/> : <div className="people-rail">{bornItems.map(p => <CatalogCard key={p.id} item={p} meta={born.meta} type="people" profile />)}</div>}
      <div className="section-heading content-section"><div><p className="eyebrow">Culture radar</p><h2>Trending people</h2></div></div>{trending.loading ? <MovieGridSkeleton count={10}/> : trending.error ? <ErrorState message={trending.error} onRetry={trending.retry}/> : <div className="catalog-grid">{trendingItems.map(p => <CatalogCard key={p.id} item={p} meta={trending.meta} type="people" profile />)}</div>}
    </section></div>;
};
export default PeoplePage;
