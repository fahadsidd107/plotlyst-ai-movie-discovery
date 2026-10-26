import './App.css'
import Home from './components/home/Home';
import Recommended from './components/recommended/Recommended';
import Review from './components/review/Review';
import Header from './components/header/Header';
import Register from './components/register/Register';
import Login from './components/login/Login';
import Layout from './components/Layout';
import RequiredAuth from './components/RequiredAuth';
import axiosClient from './api/axiosConfig';
import useAuth from './hooks/useAuth';
import AmbientStage from './components/ui/AmbientStage';
import LoadingScreen from './components/ui/LoadingScreen';
import { useEffect, useState } from 'react';
import CatalogLanding from './components/catalog/CatalogLanding';
import DetailPage from './components/catalog/DetailPage';
import PeoplePage from './components/catalog/PeoplePage';
import PersonDetail from './components/catalog/PersonDetail';
import SeasonPage from './components/catalog/SeasonPage';
import EpisodePage from './components/catalog/EpisodePage';
import LoadingLab from './components/ui/LoadingLab';
import { getMovie, getMovies, getTV, getTVDetail } from './api/catalog';

import {Route, Routes, useNavigate} from 'react-router-dom'

function App() {

  const navigate = useNavigate();
  const { auth, setAuth } = useAuth();
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
            const response = await axiosClient.post("/logout",{user_id: auth.user_id});
            console.log(response.data);
            setAuth(null);
           // localStorage.removeItem('user');
            console.log('User logged out');

        } catch (error) {
            console.error('Error logging out:', error);
        } 

    };

  if (showIntro && window.location.pathname !== '/test') return <LoadingScreen />;

  return (
    <div className="app-shell">
      <AmbientStage />
      <Header handleLogout = {handleLogout}/>
      <Routes>
        <Route path="/test" element={<LoadingLab />}></Route>
        <Route element={<Layout />}>
          <Route path="/" element={<Home updateMovieReview={updateMovieReview}/>}></Route>
          <Route path="/movies" element={<CatalogLanding kind="movie" eyebrow="The big-screen signal" title="Movies worth leaving the group chat for." copy="Fresh phenomena, modern classics and glorious oddities—programmed like the film festival in your head." loader={getMovies}/>}></Route>
          <Route path="/movie/:id" element={<DetailPage type="movie" loader={getMovie}/>}></Route>
          <Route path="/tv" element={<CatalogLanding kind="tv" eyebrow="Your next long-term commitment" title="Television that makes one more episode a threat." copy="The series everyone is talking about, and the quieter obsessions they will discover next." loader={getTV}/>}></Route>
          <Route path="/tv/:id" element={<DetailPage type="tv" loader={getTVDetail}/>}></Route>
          <Route path="/tv/:id/seasons/:season" element={<SeasonPage/>}></Route>
          <Route path="/tv/:id/seasons/:season/episodes/:episode" element={<EpisodePage/>}></Route>
          <Route path="/people" element={<PeoplePage/>}></Route>
          <Route path="/people/:id" element={<PersonDetail/>}></Route>
          <Route path="/register" element={<Register/>}></Route>
          <Route path="/login" element={<Login/>}></Route>
          <Route element = {<RequiredAuth/>}>
              <Route path="/recommended" element={<Recommended/>}></Route>
              <Route path="/review/:imdb_id" element={<Review/>}></Route>
          </Route>
        </Route>
      </Routes>
      <footer className="site-footer">
        <div className="site-footer__inner">
          <strong>Plotlyst</strong>
          <span className="site-footer__tagline">Thoughtful movie discovery, shaped around your taste.</span>
          <nav className="site-footer__contacts" aria-label="Creator contact links">
            <a href="mailto:fsiddiqui107@gmail.com">siddiqui107@gmail.com</a>
            <a href="tel:+4917623876395">+49 176 238 76 395</a>
            <a href="https://www.linkedin.com/in/muhammad-fahad-siddiqui-595484176/" target="_blank" rel="noreferrer">LinkedIn</a>
          </nav>
          <span className="site-footer__copyright">© {new Date().getFullYear()} Plotlyst</span>
          <span className="site-footer__credits">Catalogue data by <a href="https://www.themoviedb.org/" target="_blank" rel="noreferrer">TMDB</a>. Availability by <a href="https://www.justwatch.com/" target="_blank" rel="noreferrer">JustWatch</a>. This product uses the TMDB API but is not endorsed or certified by TMDB.</span>
        </div>
      </footer>
    </div>
  )
}

export default App
