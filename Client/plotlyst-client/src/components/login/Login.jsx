import {useState} from 'react';
import Container from 'react-bootstrap/Container';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import axiosClient from '../../api/axiosConfig';
import { useNavigate, Link, useLocation } from 'react-router-dom';
import useAuth from '../../hooks/useAuth';
import logo from '../../assets/PlotlystIcon.png';

const Login = () => {
    
    const {setAuth} = useAuth();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(false);
    const location = useLocation();
    const navigate = useNavigate();

    const from = location.state?.from?.pathname || "/";
    
    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError(null);       

        try {
            const response = await axiosClient.post('/login', { email, password });
            console.log(response.data);
            if (response.data.error) {
                setError(response.data.error);
                return;
            }
           // console.log(response.data);
            setAuth(response.data);
            
           // localStorage.setItem('user', JSON.stringify(response.data));
            // Handle successful login (e.g., store token, redirect)
           navigate(from, {replace: true});
           //navigate('/');

        } catch (err) {
            console.error(err);
            setError('Invalid email or password');
        } finally {
            setLoading(false);
        }
    }; 
    return (
        <Container className="auth-page">
            <div className="auth-card">
                <div className="auth-card__intro">
                    <img src={logo} alt="Plotlyst" width={56} />
                    <p className="eyebrow">Welcome back</p>
                    <h1>Sign in to your list.</h1>
                    <p>Continue to recommendations shaped around your favourite genres.</p>
                </div>
                {error && <div className="form-error" role="alert">{error}</div>}
                <Form onSubmit={handleSubmit}>
                    <Form.Group controlId="formBasicEmail" className="mb-3">
                        <Form.Label>Email address</Form.Label>
                        <Form.Control
                            type="email"
                            placeholder="Enter email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                            autoFocus
                        />
                    </Form.Group>

                    <Form.Group controlId="formBasicPassword" className="mb-3">
                        <Form.Label>Password</Form.Label>
                        <Form.Control
                            type="password"
                            placeholder="Password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                        />
                    </Form.Group>

                    <Button
                        variant="primary"
                        type="submit"
                        className="button button--primary w-100"
                        disabled={loading}
                        style={{fontWeight: 600, letterSpacing: 1}}
                    >
                        {loading ? (
                            <>
                                <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                                Logging in...
                            </>
                        ) : 'Login'}
                    </Button>
                </Form>
                <div className="auth-card__switch">
                    <span>New to Plotlyst? </span>
                    <Link to="/register">Create an account</Link>
                </div>
            </div>
        </Container>
    )
}
export default Login;
