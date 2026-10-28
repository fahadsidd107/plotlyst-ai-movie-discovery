import { useState, useEffect} from 'react';
import Container from 'react-bootstrap/Container';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import axiosClient from '../../api/axiosConfig';
import { useNavigate, Link } from 'react-router-dom';
import logo from '../../assets/PlotlystIcon.png';
import { authApi } from '../../api/step4';
import avatarSheet from '../../assets/avatars/plotlyst-avatar-sheet.png';

const avatars = [
    { id: 'avatar-01', gender: 'female' }, { id: 'avatar-02', gender: 'male' },
    { id: 'avatar-03', gender: 'female' }, { id: 'avatar-04', gender: 'male' },
    { id: 'avatar-05', gender: 'male' }, { id: 'avatar-06', gender: 'female' },
    { id: 'avatar-07', gender: 'male' }, { id: 'avatar-08', gender: 'female' },
    { id: 'avatar-09', gender: 'female' }, { id: 'avatar-10', gender: 'male' },
    { id: 'avatar-11', gender: 'female' }, { id: 'avatar-12', gender: 'male' },
];

const AvatarPortrait = ({ avatar, selected = false }) => {
    const index = Number(avatar.id.slice(-2)) - 1;
    const col = index % 4;
    const row = Math.floor(index / 4);
    return <span className={`avatar-portrait${selected ? ' is-selected' : ''}`} aria-hidden="true">
        <img src={avatarSheet} alt="" style={{ '--avatar-x': `${col * 100}%`, '--avatar-y': `${row * 100}%` }} />
    </span>;
};

const Register = () => {
    const [firstName, setFirstName] = useState('');
    const [lastName, setLastName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [dateOfBirth, setDateOfBirth] = useState('');
    const [gender, setGender] = useState('');
    const [avatar, setAvatar] = useState('');
    const [favouriteGenres, setFavouriteGenres] = useState([]);
    const [genres, setGenres] = useState([]);

    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    const handleGenreChange = (e) => {
        const options = Array.from(e.target.selectedOptions);
        setFavouriteGenres(options.map(opt => ({
            genre_id: Number(opt.value),
            genre_name: opt.label
        })));
    };
   const handleSubmit = async (e) => {
        e.preventDefault();
        setError(null);
        if (password !== confirmPassword) {
            setError('Passwords do not match.');
            return;
        }

        setLoading(true);

        try {
            const payload = {
                first_name: firstName,
                last_name: lastName,
                email,
                password,
                date_of_birth: dateOfBirth,
                gender,
                avatar,
                favourite_genres: favouriteGenres
            };
            await authApi.register(payload);
            // Registration successful, redirect to login
            navigate('/login', { replace: true });
        } catch {
            setError('Registration failed. Please try again.');
        } finally {
            setLoading(false);
        }
    };


    useEffect(() => {
        const fetchGenres = async () => {
        try {
            const response = await axiosClient.get('/genres');
            setGenres(response.data);
        } catch (error) {
            console.error('Error fetching movie genres:', error);
        }
        };
    
        fetchGenres();
    }, []);


    return (


       <Container className="auth-page">
        <div className="auth-card auth-card--wide">
                <div className="auth-card__intro">
                     <img src={logo} alt="Plotlyst" width={56} />
                    <p className="eyebrow">Start your Plotlyst</p>
                    <h1>Build a better watchlist.</h1>
                    <p>Tell us what you love and we’ll shape your recommendations.</p>
                    {error && <div className="form-error" role="alert">{error}</div>}
                </div>
             <Form onSubmit={handleSubmit} className="register-grid">
                    <Form.Group className="mb-3">
                        <Form.Label>First Name</Form.Label>
                        <Form.Control
                            type="text"
                            placeholder="Enter first name"
                            value={firstName}
                            onChange={e => setFirstName(e.target.value)}
                            required
                        />
                    </Form.Group>
                    <Form.Group className="mb-3">
                        <Form.Label>Date of birth</Form.Label>
                        <Form.Control type="date" value={dateOfBirth} max={new Date().toISOString().slice(0, 10)} onChange={e => setDateOfBirth(e.target.value)} required />
                    </Form.Group>
                    <Form.Group className="mb-3">
                        <Form.Label>Gender</Form.Label>
                        <Form.Select value={gender} onChange={e => { setGender(e.target.value); setAvatar(''); }} required>
                            <option value="">Choose an option</option>
                            <option value="male">Male</option>
                            <option value="female">Female</option>
                            <option value="diverse">Diverse</option>
                        </Form.Select>
                    </Form.Group>
                    {gender && <fieldset className="avatar-picker register-grid__full mb-3">
                        <legend>Choose your Plotlyst avatar</legend>
                        <p>Original AI-created portraits representing different ages, races, styles and body types.</p>
                        <div className="avatar-picker__grid">
                            {avatars.filter(item => gender === 'diverse' || item.gender === gender).map(item => <button type="button" key={item.id} className={avatar === item.id ? 'selected' : ''} aria-pressed={avatar === item.id} aria-label={`Choose avatar ${item.id.slice(-2)}`} onClick={() => setAvatar(item.id)}><AvatarPortrait avatar={item} selected={avatar === item.id} /></button>)}
                        </div>
                        {!avatar && <small>Please choose one avatar to continue.</small>}
                    </fieldset>}
                     <Form.Group className="mb-3">
                        <Form.Label>Last Name</Form.Label>
                        <Form.Control
                            type="text"
                            placeholder="Enter last name"
                            value={lastName}
                            onChange={e => setLastName(e.target.value)}
                            required
                        />
                    </Form.Group>
                     <Form.Group className="mb-3">
                        <Form.Label>Email</Form.Label>
                        <Form.Control
                            type="email"
                            placeholder="Enter email"
                            value={email}
                            onChange={e => setEmail(e.target.value)}
                            required
                        />
                    </Form.Group>
                     <Form.Group className="mb-3">
                        <Form.Label>Password</Form.Label>
                        <Form.Control
                            type="password"
                            placeholder="Password"
                            value={password}
                            onChange={e => setPassword(e.target.value)}
                            required
                        />
                    </Form.Group>
                     <Form.Group className="mb-3">
                        <Form.Label>Confirm Password</Form.Label>
                        <Form.Control
                            type="password"
                            placeholder="Confirm Password"
                            value={confirmPassword}
                            onChange={e => setConfirmPassword(e.target.value)}
                            required
                            isInvalid ={!!confirmPassword && password !== confirmPassword}

                        />
                        <Form.Control.Feedback type="invalid">
                            Passwords do not match.
                        </Form.Control.Feedback>
                    </Form.Group>
                    <Form.Group className="register-grid__full mb-3" controlId="favouriteGenres">
                        <Form.Label>Favourite genres</Form.Label>
                        <Form.Select
                            multiple
                            value={favouriteGenres.map(g => String(g.genre_id))}
                            onChange={handleGenreChange}
                        >
                            {genres.map(genre => (
                                <option key={genre.genre_id} value={genre.genre_id} label={genre.genre_name}>
                                    {genre.genre_name}
                                </option>
                            ))}
                        </Form.Select>
                        <Form.Text>
                            Hold Ctrl (Windows) or Cmd (Mac) to select multiple genres.
                        </Form.Text>
                    </Form.Group>
                     <Button
                        variant="primary"
                        type="submit"
                        className="button button--primary w-100 register-grid__full"
                        disabled={loading || !avatar}
                        style={{fontWeight: 600, letterSpacing: 1}}
                    >
                        {loading ? (
                            <>
                                <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                                Registering...
                            </>
                        ) : 'Register'}
                    </Button>                        
            </Form>
            <div className="auth-card__switch">Already have an account? <Link to="/login">Sign in</Link></div>
            </div>           
       </Container>

    )
}
export default Register;
