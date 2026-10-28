import './App.css'
import Home from './components/home/Home';
import Recommended from './components/recommended/Recommended';
import Review from './components/review/Review';
import Header from './components/header/Header';
import Register from './components/register/Register';
import Login from './components/login/Login';
import Layout from './components/Layout';
import RequiredAuth from './components/RequiredAuth';
import useAuth from './hooks/useAuth';
import LoadingScreen from './components/ui/LoadingScreen';
import { useEffect, useState } from 'react';
import CatalogLanding from './components/catalog/CatalogLanding';
import DetailPage from './components/catalog/DetailPage';
import PeoplePage from './components/catalog/PeoplePage';
import PersonDetail from './components/catalog/PersonDetail';
import SeasonPage from './components/catalog/SeasonPage';
import EpisodePage from './components/catalog/EpisodePage';
import LoadingLab from './components/ui/LoadingLab';
import SearchResults from './components/discovery/SearchResults';
import DiscoverPage from './components/discovery/DiscoverPage';
import { ChartsIndex, ChartPage } from './components/discovery/Charts';
import { Collections, Franchise } from './components/discovery/Collections';
import { getMovie, getMovies, getTV, getTVDetail } from './api/catalog';
import { authApi } from './api/step4';
import RoleGuard from './components/RoleGuard';
import { LibraryDashboard, WatchlistPage, HistoryPage, PlaylistPage } from './components/library/LibraryPages';
import { PublicProfile, ProfileSettings, PrivacySettings, PreferenceSettings } from './components/profile/ProfilePages';
import Moderation from './components/moderation/Moderation';
import { SpecialistIndex, Awards, AwardCeremony, BoxOffice, Streaming, Editorial } from './components/specialist/SpecialistPages';

import {Route, Routes, useNavigate} from 'react-router-dom'

function App() {

  const navigate = useNavigate();
  const { setAuth } = useAuth();
  const [showIntro, setShowIntro] = useState(() => sessionStorage.getItem('plotlyst-intro-seen') !== 'true');

  useEffect(() => {
    if (!showIntro) return undefined;
    const timer = window.setTimeout(() => {
      sessionStorage.setItem('plotlyst-intro-seen', 'true');
      setShowIntro(false);
    }, 1350);
    return () => window.clearTimeout(timer);
  }, [showIntro]);

  
  const updateMovieReview = (imdb_id) => {
      navigate(`/review/${imdb_id}`);
  };
   
  const handleLogout = async () => {

        try {
            await authApi.logout();
            setAuth(null);

        } catch (error) {
            console.error('Error logging out:', error);
        } 

    };

  if (showIntro && window.location.pathname !== '/test') return <LoadingScreen />;

  return (
    <div className="app-shell">
      <Header handleLogout = {handleLogout}/>
      <Routes>
        <Route path="/test" element={<LoadingLab />}></Route>
        <Route element={<Layout />}>
          <Route path="/" element={<Home updateMovieReview={updateMovieReview}/>}></Route>
          <Route path="/search" element={<SearchResults/>}></Route>
          <Route path="/discover/:media" element={<DiscoverPage/>}></Route>
          <Route path="/charts" element={<ChartsIndex/>}></Route>
          <Route path="/charts/:media/:size" element={<ChartPage/>}></Route>
          <Route path="/collections" element={<Collections/>}></Route>
          <Route path="/specialist" element={<SpecialistIndex/>}></Route>
          <Route path="/awards" element={<Awards/>}></Route>
          <Route path="/awards/:organization/:year" element={<AwardCeremony/>}></Route>
          <Route path="/box-office" element={<BoxOffice/>}></Route>
          <Route path="/streaming" element={<Streaming/>}></Route>
          <Route path="/editorial" element={<Editorial/>}></Route>
          <Route path="/franchise/:id" element={<Franchise/>}></Route>
          <Route path="/movies" element={<CatalogLanding kind="movie" eyebrow="The big-screen signal" title="Movies worth leaving the group chat for." copy="Fresh phenomena, modern classics and glorious oddities—programmed like the film festival in your head." loader={getMovies}/>}></Route>
          <Route path="/movie/:id" element={<DetailPage type="movie" loader={getMovie}/>}></Route>
          <Route path="/tv" element={<CatalogLanding kind="tv" eyebrow="Your next long-term commitment" title="Television that makes one more episode a threat." copy="The series everyone is talking about, and the quieter obsessions they will discover next." loader={getTV}/>}></Route>
          <Route path="/tv/:id" element={<DetailPage type="tv" loader={getTVDetail}/>}></Route>
          <Route path="/tv/:id/seasons/:season" element={<SeasonPage/>}></Route>
          <Route path="/tv/:id/seasons/:season/episodes/:episode" element={<EpisodePage/>}></Route>
          <Route path="/people" element={<PeoplePage/>}></Route>
          <Route path="/people/:id" element={<PersonDetail/>}></Route>
          <Route path="/profile/:username" element={<PublicProfile/>}></Route>
          <Route path="/register" element={<Register/>}></Route>
          <Route path="/login" element={<Login/>}></Route>
          <Route element = {<RequiredAuth/>}>
              <Route path="/profile" element={<ProfileSettings/>}></Route>
              <Route path="/recommended" element={<Recommended/>}></Route>
              <Route path="/review/:imdb_id" element={<Review/>}></Route>
              <Route path="/library" element={<LibraryDashboard/>}></Route>
              <Route path="/library/watchlist" element={<WatchlistPage/>}></Route>
              <Route path="/library/history" element={<HistoryPage/>}></Route>
              <Route path="/playlists/:id" element={<PlaylistPage/>}></Route>
              <Route path="/settings/profile" element={<ProfileSettings/>}></Route>
              <Route path="/settings/privacy" element={<PrivacySettings/>}></Route>
              <Route path="/settings/preferences" element={<PreferenceSettings/>}></Route>
          </Route>
          <Route element={<RoleGuard roles={['MODERATOR','ADMIN']}/>}><Route path="/moderation" element={<Moderation/>}/></Route>
        </Route>
      </Routes>
      <footer className="site-footer">
        <div className="site-footer__inner">
          <strong>Plotlyst</strong>
          <span className="site-footer__tagline">Thoughtful movie discovery, shaped around your taste.</span>
          <nav className="site-footer__contacts" aria-label="Creator contact links">
            <a href="/specialist">Specialist data</a>
            <a href="/streaming">Where to watch</a>
            <a href="mailto:siddiqui107@gmail.com">siddiqui107@gmail.com</a>
            <a href="tel:+4917623876395">+49 176 238 76 395</a>
            <a href="https://www.linkedin.com/in/muhammad-fahad-siddiqui-595484176/" target="_blank" rel="noreferrer">LinkedIn</a>
          </nav>
          <span className="site-footer__copyright">© {new Date().getFullYear()} Plotlyst</span>
          <span className="site-footer__credits"><a className="tmdb-credit" href="https://www.themoviedb.org/" target="_blank" rel="noreferrer"><img src="https://www.themoviedb.org/assets/2/v4/logos/v2/blue_square_2-d537fb228cf3ded904ef09b136fe3fec72548ebc1fea3fbbd1ad9e36364db38b.svg" alt="TMDB" /></a> Catalogue data by TMDB. Availability by <a href="https://www.justwatch.com/" target="_blank" rel="noreferrer">JustWatch</a>. This product uses the TMDB API but is not endorsed or certified by TMDB.</span>
        </div>
      </footer>
    </div>
  )
}

export default App
