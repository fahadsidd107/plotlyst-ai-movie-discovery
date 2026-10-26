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

  if (showIntro) return <LoadingScreen />;

  return (
    <div className="app-shell">
      <AmbientStage />
      <Header handleLogout = {handleLogout}/>
      <Routes>
        <Route element={<Layout />}>
          <Route path="/" element={<Home updateMovieReview={updateMovieReview}/>}></Route>
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
            <a href="mailto:siddiqui107@gmail.com">siddiqui107@gmail.com</a>
            <a href="tel:+4917623876395">+49 176 238 76 395</a>
            <a href="https://de.linkedin.com/in/muhammad-fahad-siddiqui-595484176" target="_blank" rel="noreferrer">LinkedIn</a>
          </nav>
          <span className="site-footer__copyright">© {new Date().getFullYear()} Plotlyst</span>
        </div>
      </footer>
    </div>
  )
}

export default App
