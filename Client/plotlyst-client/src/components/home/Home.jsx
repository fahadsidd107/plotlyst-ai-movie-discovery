import { useEffect, useRef, useState } from 'react';
import { Link } from 'react-router-dom';
import { getDiscoveryHome, getHomeSection, imageUrl } from '../../api/catalog';
import ProgrammeRail from '../discovery/ProgrammeRail';

const secondary = [
  ['trending-movies', 'Movies moving culture', 'The films turning every group chat into a spoiler minefield.', 'movie', '/discover/movie?sort=popularity.desc'],
  ['trending-tv', 'Series everyone is discussing', 'Make room in the diary. These stories refuse to stay on one screen.', 'tv', '/discover/tv?sort=popularity.desc'],
  ['fan-favourites', 'Applause that lasts', 'Beloved films with the votes to prove it.', 'movie', '/charts/movie/250'],
  ['airing-today', 'On television tonight', 'A fresh episode is already waiting.', 'tv', '/discover/tv?sort=first_air_date.desc'],
];

export default function Home() {
  const heroRef = useRef(null);
  const [primary, setPrimary] = useState({ loading: true });
  const [reels, setReels] = useState(() => Object.fromEntries(secondary.map(([key]) => [key, { loading: true }])));
  useEffect(() => {
    const controller = new AbortController();
    getDiscoveryHome(controller.signal).then(({data}) => {
      const payload = data.data || data;
      const section = payload.sections?.find(item => item.key === 'now_playing');
      setPrimary({ ...section, meta: data.meta, loading: false });
    }).catch(error => { if (error.name !== 'CanceledError') setPrimary({ loading: false, error: 'The opening reel missed its cue.' }); });
    secondary.forEach(([key]) => getHomeSection(key, controller.signal).then(({data}) => {
      const payload = data.data || data;
      setReels(current => ({ ...current, [key]: { ...payload, meta: data.meta, loading: false } }));
    }).catch(error => { if (error.name !== 'CanceledError') setReels(current => ({ ...current, [key]: { loading: false, error: 'This reel is temporarily unavailable.' } })); }));
    return () => controller.abort();
  }, []);
  useEffect(() => {
    const hero=heroRef.current; if(!hero)return undefined;
    const reduced=window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    const fine=window.matchMedia('(pointer: fine)').matches;
    const move=event=>{if(reduced||!fine)return;const box=hero.getBoundingClientRect();const x=(event.clientX-box.left)/box.width-.5;const y=(event.clientY-box.top)/box.height-.5;hero.style.setProperty('--pointer-x',`${x*2}`);hero.style.setProperty('--pointer-y',`${y*2}`)};
    const leave=()=>{hero.style.setProperty('--pointer-x','0');hero.style.setProperty('--pointer-y','0')};
    const scroll=()=>{if(reduced)return;const box=hero.getBoundingClientRect();hero.style.setProperty('--film-progress',String(Math.max(-1,Math.min(1,-box.top/Math.max(box.height,1)))))};
    hero.addEventListener('pointermove',move);hero.addEventListener('pointerleave',leave);window.addEventListener('scroll',scroll,{passive:true});scroll();
    return()=>{hero.removeEventListener('pointermove',move);hero.removeEventListener('pointerleave',leave);window.removeEventListener('scroll',scroll)};
  }, []);
  const feature = primary.results?.[0];
  return <div className="discovery-home">
    <section ref={heroRef} className="premiere-hero premiere-hero--3d" style={feature?.backdrop_path ? { '--premiere-art': `url(${imageUrl(primary.meta, feature.backdrop_path, 'original')})` } : undefined}>
      <div className="film-perforations" aria-hidden="true" />
      <div className="projection-set" aria-hidden="true"><div className="projection-set__beam"/><div className="projection-set__frame projection-set__frame--rear"/><div className="projection-set__frame projection-set__frame--mid"/><div className="projection-set__frame projection-set__frame--front"/><div className="projection-set__reel"><i/><i/><i/></div><div className="projection-set__film"><span/><span/><span/><span/><span/></div></div>
      <div className="premiere-hero__content page-shell">
        <p className="eyebrow">Tonight at Plotlyst · Programme 001</p>
        <h1>{feature?.title || 'The curtain is about to rise.'}</h1>
        <p>{feature?.overview || 'Fresh from the big screen, hand-programmed for your next night at the movies.'}</p>
        <div className="premiere-actions">
          {feature && <Link className="ticket-button" to={`/movie/${feature.id}`}>Enter the screening <span>→</span></Link>}
          <Link className="programme-link" to="/discover/movie">Browse all films</Link>
        </div>
        <div className="screening-slip" tabIndex="0"><span>NOW PLAYING</span><b>{feature?.release_date?.slice(0,4) || '2026'}</b><small>{feature?.vote_average ? `AUDIENCE ${feature.vote_average.toFixed(1)}` : 'CURTAIN UP'}</small></div>
      </div>
    </section>
    <div className="page-shell programme-list">
      <ProgrammeRail number={1} title="Now playing in cinemas" copy="The current big-screen programme, from opening credits onward." items={primary.results?.slice(1)} meta={primary.meta} type="movie" loading={primary.loading} error={primary.error} destination="/discover/movie?sort=primary_release_date.desc" />
      {secondary.map(([key,title,copy,type,destination], index) => <ProgrammeRail key={key} number={index+2} title={title} copy={copy} items={reels[key].results} meta={reels[key].meta} type={type} loading={reels[key].loading} error={reels[key].error} destination={destination}/>)}
      <section className="chart-gateway"><span className="leader-count">3</span><div><p className="eyebrow">The permanent programme</p><h2>Every great film needs a place in history.</h2><p>Enter reproducible charts, decade rooms and genre shelves built for serious browsing.</p></div><Link className="ticket-button" to="/charts">Open the charts <span>→</span></Link></section>
    </div>
  </div>;
}
