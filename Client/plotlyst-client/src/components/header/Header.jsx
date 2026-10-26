import Button from 'react-bootstrap/Button'
import Container from 'react-bootstrap/Container'
import Nav from 'react-bootstrap/Nav'
import Navbar from 'react-bootstrap/Navbar'
import {useNavigate, NavLink, Link} from 'react-router-dom'
import useAuth from '../../hooks/useAuth';
import logo from '../../assets/PlotlystLogo.png';

const Header = ({handleLogout}) => {
    const navigate = useNavigate();
    const {auth} = useAuth();


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
                        <Nav.Link as = {NavLink} to="/recommended">
                            For You
                        </Nav.Link>
                    </Nav>
    
                    <Nav className ="account-nav ms-auto align-items-md-center">
                        {auth ? (
                        <>
                            <span className="user-greeting">
                                Hello, <strong>{auth.first_name}</strong>
                            </span>
                            <Button className="button button--secondary" onClick={handleLogout}>
                                Logout
                            </Button>
                        </>
                        ):(
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
