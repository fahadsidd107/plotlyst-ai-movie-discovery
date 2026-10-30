import { useCallback, useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { debateApi, debateMessage } from '../../api/debate';

const buckets = [
  ['all', 'All ballots'],
  ['high', 'High · 4–5★'],
  ['middle', 'Middle · 3★'],
  ['low', 'Low · 1–2★'],
];

const distributionEntries = distribution => {
  if (Array.isArray(distribution)) return [5, 4, 3, 2, 1].map(star => {
    const entry = distribution.find(item => Number(item.stars ?? item.rating ?? item._id) === star) || {};
    return { star, count: Number(entry.count ?? entry.total ?? entry.value ?? 0) };
  });
  return [5, 4, 3, 2, 1].map(star => ({ star, count: Number(distribution?.[star] ?? distribution?.[String(star)] ?? 0) }));
};
const isAiMethod = method => /ai|openai|langchain/i.test(method || '');

function Evidence({ evidence = [], movieId }) {
  if (!evidence.length) return <p className="debate-no-evidence">No public excerpt is available for this observation.</p>;
  return <ol className="debate-evidence">
    {evidence.map((review, index) => <li key={review.review_id || index}>
      <header><strong>{review.author_name || 'Plotlyst member'}</strong><span>{review.stars ? `${review.stars}★` : 'Unrated'}</span></header>
      {review.spoiler_hidden
        ? <p className="debate-evidence__spoiler"><span aria-hidden="true">◼</span> Excerpt concealed by your spoiler settings.</p>
        : <blockquote>{review.excerpt || 'The author has not made an excerpt available.'}</blockquote>}
      <Link to={`/movie/${movieId}?review=${encodeURIComponent(review.review_id)}#review-${encodeURIComponent(review.review_id)}`}>
        Original review <span className="sr-only">by {review.author_name || 'Plotlyst member'}</span> · #{review.review_id}
      </Link>
    </li>)}
  </ol>;
}

function Observation({ item, tone, method, movieId }) {
  const [open, setOpen] = useState(false);
  const id = `${tone}-${String(item.theme || 'observation').replace(/\W+/g, '-').toLowerCase()}`;
  return <article className={`debate-observation debate-observation--${tone}`}>
    <div className="debate-observation__tally" aria-hidden="true">{Array.from({ length: Math.min(Number(item.mention_count) || 1, 5) }, (_, index) => <i key={index} />)}</div>
    <p className="debate-observation__source">{isAiMethod(method) ? 'AI-assisted synthesis' : 'Deterministic review analysis'}</p>
    <h4>{item.theme || 'Audience observation'}</h4>
    <p>{item.summary}</p>
    <button type="button" aria-expanded={open} aria-controls={id} onClick={() => setOpen(value => !value)}>
      {open ? 'Close transcript' : `Open evidence · ${item.mention_count || item.evidence?.length || 0} mentions`}
    </button>
    <div id={id} hidden={!open}><p className="debate-evidence-label">Original member reviews · not generated</p><Evidence evidence={item.evidence} movieId={movieId} /></div>
  </article>;
}

export default function DebatePanel({ movieId, title }) {
  const [bucket, setBucket] = useState('all');
  const [data, setData] = useState(null);
  const [state, setState] = useState('loading');
  const [error, setError] = useState('');
  const load = useCallback(() => {
    setState('loading'); setError('');
    debateApi.movie(movieId, bucket).then(result => { setData(result); setState('ready'); }).catch(reason => { setError(debateMessage(reason)); setState('error'); });
  }, [bucket, movieId]);
  useEffect(load, [load]);
  const distribution = useMemo(() => distributionEntries(data?.rating_distribution), [data]);
  const total = distribution.reduce((sum, entry) => sum + entry.count, 0) || Number(data?.review_count) || 0;
  const max = Math.max(...distribution.map(entry => entry.count), 1);
  const method = data?.analysis_method || 'deterministic';

  if (state === 'loading') return <section className="debate-panel debate-panel--loading" aria-busy="true"><div className="debate-loading-mark" aria-hidden="true"><i/><i/><i/></div><div><p className="eyebrow">Counting the house</p><h3>Developing the audience ballot…</h3></div></section>;
  if (state === 'error') return <section className="debate-panel debate-state debate-state--error" role="alert"><p className="eyebrow">Ballot interrupted</p><h3>The debate desk lost the transcript.</h3><p>{error}</p><button type="button" onClick={load}>Count again</button></section>;
  if (data?.status === 'insufficient') return <section className="debate-panel debate-state debate-state--insufficient"><span className="debate-stamp">Not enough ballots</span><p className="eyebrow">The honest result</p><h3>{data.review_count || 0} of {data.minimum_reviews || 0} reviews collected.</h3><p>Plotlyst waits for a meaningful sample before describing what viewers loved or disliked. No community sentiment has been generated yet.</p><meter min="0" max={Math.max(data.minimum_reviews || 1, 1)} value={data.review_count || 0}>{data.review_count || 0} reviews</meter></section>;

  return <section className="debate-panel" aria-labelledby="debate-title">
    <header className="debate-masthead"><div><p className="eyebrow">Post-screening roundtable</p><h3 id="debate-title">Why the house argued about <em>{title}</em>.</h3><p>Generated observations are grounded in {data?.analyzed_review_count || data?.review_count || total} genuine Plotlyst reviews. Open any transcript to inspect the supporting words.</p></div><div className="debate-stamp">{isAiMethod(method) ? 'AI-assisted' : 'Rules-based'}<small>community synthesis</small></div></header>
    {(data?.status === 'stale' || data?.new_reviews_pending) && <p className="debate-pending" role="status">The current programme is still valid; {data.new_reviews_pending || 'new'} review{data.new_reviews_pending === 1 ? '' : 's'} await the next recount.</p>}
    <div className="debate-ballot">
      <div className="debate-ballot__heading"><div><p className="eyebrow">Audience ballot</p><h4>Rating distribution</h4></div><strong>{total}<small>counted reviews</small></strong></div>
      <div className="debate-bars" aria-label="Rating distribution">
        {distribution.map(entry => { const percent = total ? Math.round(entry.count / total * 100) : 0; return <div key={entry.star}><b>{entry.star}★</b><span aria-hidden="true"><i style={{ '--bar-height': `${entry.count / max * 100}%` }} /></span><small>{entry.count}<br/>{percent}%</small></div>; })}
      </div>
      <div className="debate-filters" aria-label="Filter debate by rating category">{buckets.map(([value, label]) => <button type="button" aria-pressed={bucket === value} key={value} onClick={() => setBucket(value)}>{label}</button>)}</div>
      <div className="debate-rail" style={{ '--rail-position': bucket === 'low' ? '12%' : bucket === 'middle' ? '50%' : bucket === 'high' ? '88%' : '50%' }} aria-hidden="true"><i /></div>
      {data?.polarization && <p className="debate-polarization"><b>Audience disagreement:</b> {typeof data.polarization === 'string' ? data.polarization : data.polarization.summary || data.polarization.explanation || data.polarization.label || 'Ratings are spread unusually far apart.'}</p>}
    </div>
    <div className="debate-sides">
      <section aria-labelledby="loved-title"><header><span>FOR</span><div><p className="eyebrow">The applause</p><h3 id="loved-title">What viewers loved</h3></div></header>{data?.loved?.length ? data.loved.map((item, index) => <Observation key={`${item.theme}-${index}`} item={item} tone="loved" method={method} movieId={movieId} />) : <p className="debate-column-empty">No supported positive observation appears in this rating group.</p>}</section>
      <section aria-labelledby="disliked-title"><header><span>AGAINST</span><div><p className="eyebrow">The reservations</p><h3 id="disliked-title">What viewers disliked</h3></div></header>{data?.disliked?.length ? data.disliked.map((item, index) => <Observation key={`${item.theme}-${index}`} item={item} tone="disliked" method={method} movieId={movieId} />) : <p className="debate-column-empty">No supported negative observation appears in this rating group.</p>}</section>
    </div>
    <footer>Analysis {data?.generated_at ? `last counted ${new Date(data.generated_at).toLocaleDateString()}` : 'generated from the current review sample'} · Summaries are generated; quotations belong to their authors.</footer>
  </section>;
}
