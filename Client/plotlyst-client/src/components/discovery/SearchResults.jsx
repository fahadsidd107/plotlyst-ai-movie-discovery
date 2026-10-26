import { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { searchMulti } from '../../api/catalog';
import DiscoveryCard from './DiscoveryCard';
import { ErrorState } from '../ui/States';

export default function SearchResults() {
  const [params, setParams] = useSearchParams();
  const query = params.get('q') || '';
  const page = Number(params.get('page')) || 1;
  const [type, setType] = useState('all');
  const [state, setState] = useState({ loading: true });
  useEffect(() => { const controller = new AbortController(); setState({loading:true}); searchMulti(query,page,controller.signal).then(({data}) => setState({...(data.data||data),meta:data.meta,loading:false})).catch(error => {if(error.name!=='CanceledError') setState({loading:false,error:'Search could not reach the projection booth.'})}); return ()=>controller.abort(); },[query,page]);
  const results = state.results?.filter(item => type === 'all' || item.media_type === type) || [];
  return <div className="page-shell discovery-page"><header className="programme-page-header"><p className="eyebrow">Archive search</p><h1>Results for “{query}”</h1><p>Films, television and the people behind them—one index, no dead ends.</p></header>
    <div className="search-tabs" role="tablist" aria-label="Result type" onKeyDown={event=>{const tabs=['all','movie','tv','person'];const current=tabs.indexOf(type);if(event.key==='ArrowRight'){event.preventDefault();setType(tabs[(current+1)%tabs.length])}if(event.key==='ArrowLeft'){event.preventDefault();setType(tabs[(current+tabs.length-1)%tabs.length])}}}>{['all','movie','tv','person'].map(tab=><button id={`result-tab-${tab}`} role="tab" aria-controls="search-results-panel" tabIndex={type===tab?0:-1} aria-selected={type===tab} className={type===tab?'active':''} key={tab} onClick={()=>setType(tab)}>{tab==='tv'?'TV':tab[0].toUpperCase()+tab.slice(1)}</button>)}</div>
    <div id="search-results-panel" role="tabpanel" aria-labelledby={`result-tab-${type}`}>{state.loading?<div className="discovery-grid">{[1,2,3,4,5].map(i=><div className="programme-placeholder" key={i}/>)}</div>:state.error?<ErrorState message={state.error}/>:results.length?<div className="discovery-grid">{results.map(item=><DiscoveryCard key={`${item.media_type}-${item.id}`} item={item} meta={state.meta} type={item.media_type}/>)}</div>:<div className="empty-programme"><b>Nothing on this reel.</b><p>Try another title, performer or phrase.</p></div>}</div>
    {state.total_pages>1&&<nav className="pagination-bar" aria-label="Search pages"><button disabled={page<=1} onClick={()=>setParams({q:query,page:page-1})}>Previous</button><span>Page {page} of {Math.min(state.total_pages,500)}</span><button disabled={page>=state.total_pages} onClick={()=>setParams({q:query,page:page+1})}>Next</button></nav>}
  </div>;
}
