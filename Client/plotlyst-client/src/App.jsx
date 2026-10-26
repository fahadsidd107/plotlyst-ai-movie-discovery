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

import {Route, Routes, useNavigate} from 'react-router-dom'

function App() {

  const navigate = useNavigate();
  const { auth, setAuth } = useAuth();

  
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

  return (
    <div className="app-shell">
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
          <span>Thoughtful movie discovery, shaped around your taste.</span>
          <span>© {new Date().getFullYear()} Plotlyst</span>
        </div>
      </footer>
    </div>
  )
}

export default App
