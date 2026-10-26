import Button from 'react-bootstrap/Button'
import Container from 'react-bootstrap/Container'
import Nav from 'react-bootstrap/Nav'
import Navbar from 'react-bootstrap/Navbar'
import Form from 'react-bootstrap/Form'
import {useNavigate, NavLink, Link} from 'react-router-dom'
import { useState, useEffect, useRef } from 'react';
import useAuth from '../../hooks/useAuth';
import { searchMulti } from '../../api/catalog';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSearch } from '@fortawesome/free-solid-svg-icons';
import logo from '../../assets/PlotlystIcon.png';

const Header = ({handleLogout}) => {
    const navigate = useNavigate();
    const {auth} = useAuth();
    const [searchQuery, setSearchQuery] = useState('');
    const [suggestions, setSuggestions] = useState([]);
    const [showSuggestions, setShowSuggestions] = useState(false);
    const searchRef = useRef(null);

    // Debounce search - only send request after user stops typing for 300ms and after 2 characters
    useEffect(() => {
        const timer = setTimeout(async () => {
            if (searchQuery.length >= 2) {
                try {
                    const response = await searchMulti(searchQuery);
                    setSuggestions(response.data?.data?.results || []);
                    setShowSuggestions(true);
                } catch (error) {
                    console.error('Search error:', error);
                    setSuggestions([]);
                }
            } else {
                setSuggestions([]);
                setShowSuggestions(false);
            }
        }, 300);

        return () => clearTimeout(timer);
    }, [searchQuery]);

    // Close suggestions when clicking outside
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (searchRef.current && !searchRef.current.contains(event.target)) {
                setShowSuggestions(false);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleSuggestionClick = (item) => {
        setShowSuggestions(false);
        setSearchQuery('');
        setSuggestions([]);
        
        // Navigate based on media type
        if (item.media_type === 'movie') {
            navigate(`/movie/${item.id}`);
        } else if (item.media_type === 'tv') {
            navigate(`/tv/${item.id}`);
        } else if (item.media_type === 'person') {
            navigate(`/people/${item.id}`);
        }
    };

    const handleSearchSubmit = (e) => {
        e.preventDefault();
        if (searchQuery.trim()) {
            setShowSuggestions(false);
            // If there are suggestions, navigate to first one or could create a search results page
            if (suggestions.length > 0) {
                handleSuggestionClick(suggestions[0]);
            }
        }
    };


    return (
        <>
        <a className="skip-link" href="#main-content">Skip to content</a>
        <Navbar variant='dark' expand="md" sticky="top" className="site-header">
            <Container className="site-header__inner">
                <Navbar.Brand as={Link} to="/" className="brand-lockup">
                     <img
                        alt=""
                        src={logo}
                        width="38"
                        height="38"
                        className="brand-logo"
                    />
                    Plotlyst
                </Navbar.Brand>

            <Navbar.Toggle aria-controls="main-navbar-nav" aria-label="Open navigation" />
                <Navbar.Collapse id="main-navbar-nav">
                    <Nav className ="main-nav me-auto">
                        <Nav.Link as = {NavLink} to="/">
                            Home
                        </Nav.Link>
                        <Nav.Link as={NavLink} to="/movies">Movies</Nav.Link>
                        <Nav.Link as={NavLink} to="/tv">TV</Nav.Link>
                        <Nav.Link as={NavLink} to="/people">People</Nav.Link>
                        <Nav.Link as = {NavLink} to="/recommended">
                            For You
                        </Nav.Link>
                    </Nav>
    
                    {/* Search Bar */}
                    <div ref={searchRef} className="search-container me-md-3 position-relative">
                        <Form onSubmit={handleSearchSubmit} className="d-flex">
                            <div className="search-input-wrapper position-relative">
                                <FontAwesomeIcon icon={faSearch} className="search-icon" />
                                <Form.Control
                                    type="search"
                                    placeholder="Search movies, TV, people..."
                                    className="search-input"
                                    value={searchQuery}
                                    onChange={(e) => setSearchQuery(e.target.value)}
                                    onFocus={() => searchQuery.length >= 2 && setShowSuggestions(true)}
                                />
                            </div>
                        </Form>
                        {/* Suggestions Dropdown */}
                        {showSuggestions && suggestions.length > 0 && (
                            <div className="search-suggestions position-absolute w-100 mt-1">
                                {suggestions.map((item, index) => {
                                    const imagePath = item.media_type === 'person' ? item.profile_path : item.poster_path;
                                    const imageUrl = imagePath ? `https://image.tmdb.org/t/p/w92${imagePath}` : null;
                                    const mediaTypeLabel = item.media_type === 'movie' ? 'Movie' : item.media_type === 'tv' ? 'Series' : 'Actor';
                                    return (
                                        <div
                                            key={`${item.media_type}-${item.id}-${index}`}
                                            className="suggestion-item p-2 d-flex align-items-center gap-3"
                                            onClick={() => handleSuggestionClick(item)}
                                        >
                                            {imageUrl ? (
                                                <img 
                                                    src={imageUrl} 
                                                    alt={item.title || item.name}
                                                    className="suggestion-image rounded"
                                                    style={{ width: '40px', height: '60px', objectFit: 'cover' }}
                                                />
                                            ) : (
                                                <div 
                                                    className="suggestion-image-placeholder rounded"
                                                    style={{ width: '40px', height: '60px', background: 'var(--surface)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}
                                                >
                                                    <span>{(item.title || item.name)?.charAt(0)}</span>
                                                </div>
                                            )}
                                            <div className="suggestion-content d-flex flex-column">
                                                <span className="suggestion-title">
                                                    {item.title || item.name}
                                                </span>
                                                <small className="suggestion-type" style={{ color: '#b3b3b3', fontSize: '0.75rem' }}>
                                                    {mediaTypeLabel}
                                                </small>
                                            </div>
                                        </div>
                                    );
                                })}
                            </div>
                        )}
                    </div>

                    <Nav className="account-nav ms-auto align-items-md-center">
                        {auth ? (
                        <>
                            <span className="user-greeting">
                                Hello, <strong>{auth.first_name}</strong>
                            </span>
                            <Button className="button button--secondary" onClick={handleLogout}>
                                Logout
                            </Button>
                        </>
                        ) : (
                            <>
                                <Button
                                    className="button button--secondary"
                                    onClick={() => navigate("/login")} 
                                >
                                    Login
                                </Button>
                                <Button
                                    className="button button--primary"
                                    onClick={() => navigate("/register")}  
                                >
                                    Register
                                </Button>                        
                            </>
                        )}
                    </Nav>       
                </Navbar.Collapse>
            </Container>
        </Navbar>
        </>
    )
}
export default Header;