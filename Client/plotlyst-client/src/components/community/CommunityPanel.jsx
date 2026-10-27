import { useCallback, useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import useAuth from '../../hooks/useAuth';
import { apiMessage, communityApi, libraryApi } from '../../api/step4';
import StarRating from './StarRating';

const authorName = review => review.author?.username || review.author_name || review.username || 'Plotlyst member';
const reviewId = review => review.id || review._id;

export default function CommunityPanel({ mediaType, mediaId, title, posterPath }) {
  const { auth } = useAuth();
  const [ratings, setRatings] = useState(null); const [reviews, setReviews] = useState([]);
  const [body, setBody] = useState(''); const [spoiler, setSpoiler] = useState(false); const [revealed, setRevealed] = useState({});
  const [status, setStatus] = useState('loading'); const [notice, setNotice] = useState('');
  const load = useCallback(async () => { setStatus('loading'); try { const [r, v] = await Promise.all([communityApi.ratings(mediaType, mediaId), communityApi.reviews(mediaType, mediaId)]); setRatings(r); setReviews(v.results || v.reviews || v || []); setStatus('ready'); } catch (error) { setNotice(apiMessage(error)); setStatus('error'); } }, [mediaId, mediaType]);
  useEffect(() => { load(); }, [load]);
  const act = async (action, success) => { try { setNotice(''); await action(); setNotice(success); await load(); } catch (error) { setNotice(apiMessage(error)); } };
  const submit = event => { event.preventDefault(); if (!body.trim()) return; act(() => communityApi.createReview(mediaType, mediaId, { body: body.trim(), spoiler }), 'Programme note published.'); setBody(''); setSpoiler(false); };
  const addWatchlist = () => act(() => libraryApi.watchlistAdd({ media_type: mediaType, media_id: Number(mediaId), title, poster_path: posterPath }), 'Saved to your private watchlist.');
  return <section className="community-section" aria-labelledby="community-title">
    <div className="community-heading"><div><p className="eyebrow">After the credits</p><h2 id="community-title">Rate it. Note it. Remember it.</h2></div>{auth ? <button className="archive-action" onClick={addWatchlist}>＋ Private watchlist</button> : <Link className="programme-link" to="/login">Sign in to join →</Link>}</div>
    {notice && <p className="community-notice" role="status">{notice}</p>}
    <div className="rating-ledger">
      <div><span>Audience average</span><strong>{Number(ratings?.average || 0).toFixed(1)}</strong><small>{ratings?.count || 0} ticket punches</small></div>
      <StarRating value={ratings?.viewer_rating || 0} disabled={!auth} onChange={rating => act(() => communityApi.rate(mediaType, mediaId, rating), 'Your rating is locked in.')} />
      <div className="rating-bars">{[5,4,3,2,1].map(n => <span key={n}><b>{n}</b><i style={{'--rating-width': `${ratings?.distribution?.[n] || 0}%`}} /></span>)}</div>
    </div>
    {auth && <form className="review-composer" onSubmit={submit}><label htmlFor="programme-note">Your programme note</label><textarea id="programme-note" value={body} maxLength="4000" onChange={event => setBody(event.target.value)} placeholder="What stayed with you after the lights came up?" required /><div><label className="check-line"><input type="checkbox" checked={spoiler} onChange={event => setSpoiler(event.target.checked)} /> Contains spoilers</label><button className="ticket-button">Publish note</button></div></form>}
    <div className="review-thread">
      <div className="section-heading"><div><p className="eyebrow">Programme notes</p><h3>{reviews.length ? `${reviews.length} member reviews` : 'The first note is yours'}</h3></div></div>
      {status === 'loading' && <p className="thread-state">Fetching notes from the archive…</p>}
      {status === 'error' && <button className="archive-action" onClick={load}>Retry community notes</button>}
      {status === 'ready' && reviews.map(review => { const id = reviewId(review); const hidden = review.spoiler && !revealed[id]; return <article className={`review-note${hidden ? ' review-note--spoiler' : ''}`} key={id}>
        <header><div className="avatar-stub">{authorName(review).slice(0, 1).toUpperCase()}</div><div><Link to={`/profile/${authorName(review)}`}>@{authorName(review)}</Link><small>{review.created_at ? new Date(review.created_at).toLocaleDateString() : 'Recently'} · {review.rating ? `${review.rating}/5` : 'Unrated'}</small></div></header>
        {hidden ? <button className="spoiler-gate" onClick={() => setRevealed(current => ({...current, [id]: true}))}><b>Spoiler curtain</b><span>Reveal this programme note</span></button> : <p>{review.body || review.content}</p>}
        <footer><button onClick={() => act(() => communityApi.helpful(id), 'Marked as helpful.')}>♡ Helpful {review.helpful_count || 0}</button>{auth && <button onClick={() => { const reply = window.prompt('Write a reply'); if (reply?.trim()) act(() => communityApi.reply(id, reply.trim()), 'Reply added.'); }}>Reply</button>}{auth && <button onClick={() => { const reason = window.prompt('Why are you reporting this review?'); if (reason?.trim()) act(() => communityApi.report(id, reason.trim()), 'Report sent privately to moderators.'); }}>Report</button>}{auth?.user_id === (review.user_id || review.author?.id) && <button onClick={() => act(() => communityApi.deleteReview(id), 'Review removed.')}>Delete</button>}</footer>
        {review.replies?.map(reply => <div className="review-reply" key={reply.id || reply._id}><b>@{authorName(reply)}</b><p>{reply.body || reply.content}</p></div>)}
      </article>; })}
    </div>
  </section>;
}
