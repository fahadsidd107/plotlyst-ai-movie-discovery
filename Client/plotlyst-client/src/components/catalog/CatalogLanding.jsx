import { useState } from 'react';
import useCatalog from '../../hooks/useCatalog';
import CatalogCard from './CatalogCard';
import { EmptyState, ErrorState, MovieGridSkeleton } from '../ui/States';

const CatalogLanding = ({ kind, eyebrow, title, copy, loader }) => {
  const [page, setPage] = useState(1);
  const [collection, setCollection] = useState('trending');
  const { data, meta, loading, error, retry } = useCatalog(() => loader(page, collection), [page, collection]);
  const items = data?.results || data || [];
  return <div className="catalog-page">
    <header className="catalog-hero page-shell">
      <div><p className="eyebrow">{eyebrow}</p><h1>{title}</h1><p>{copy}</p></div>
      <div className="signal-rail" aria-hidden="true"><i /><i /><i /><i /><i /></div>
    </header>
    <section className="page-shell catalogue-section" aria-labelledby={`${kind}-catalogue`}>
      <div className="section-heading"><div><p className="eyebrow">Now transmitting</p><h2 id={`${kind}-catalogue`}>{collection === 'trending' ? 'Trending now' : 'Popular picks'}</h2></div>
        <div className="catalog-switch" aria-label="Catalogue collection"><button className={collection === 'trending' ? 'active' : ''} onClick={() => { setPage(1); setCollection('trending'); }}>Trending</button><button className={collection === 'popular' ? 'active' : ''} onClick={() => { setPage(1); setCollection('popular'); }}>Popular</button></div>
      </div>
      {loading ? <MovieGridSkeleton count={10} /> : error ? <ErrorState message={error} onRetry={retry} /> : items.length ? <div className="catalog-grid">{items.map(item => <CatalogCard key={item.id} item={item} meta={meta} type={kind} />)}</div> : <EmptyState />}
      {!loading && !error && <nav className="pagination-bar" aria-label="Catalogue pages"><button disabled={page === 1} onClick={() => setPage(p => Math.max(1, p - 1))}>← Previous</button><span>Transmission {page}</span><button disabled={data?.total_pages && page >= data.total_pages} onClick={() => setPage(p => p + 1)}>Next →</button></nav>}
    </section>
  </div>;
};
export default CatalogLanding;
