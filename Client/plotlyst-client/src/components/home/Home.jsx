import {useState, useEffect} from 'react';
import axiosClient from '../../api/axiosConfig'
import Movies from '../movies/Movies';
import PageIntro from '../ui/PageIntro';
import { MovieGridSkeleton, ErrorState } from '../ui/States';

const Home =({updateMovieReview}) => {
    const [movies, setMovies] = useState([]);
    const [loading, setLoading] = useState(false)
    const [message, setMessage] = useState();
    const [reloadKey, setReloadKey] = useState(0);

    useEffect(() => {
        const fetchMovies = async () => {
            setLoading(true);
            setMessage("");
            try{
                const response = await axiosClient.get('/movies');
                setMovies(response.data);
                if (response.data.length === 0){
                    setMessage('There are currently no movies available')
                }

            }catch(error){
                console.error('Error fetching movies:', error)
                setMessage("The catalogue could not be loaded. Check your connection and try again.")
            }finally{
                setLoading(false)
            }
        }
        fetchMovies();
    }, [reloadKey]);

    return (
        <div className="page-shell">
            <PageIntro eyebrow="Your next story starts here" title="Find a film worth talking about." constellation
                meta={!loading && !message ? `${movies.length} ${movies.length === 1 ? 'title' : 'titles'} in the current catalogue` : undefined}>
                Browse a focused catalogue, see its editorial ranking, and discover what fits your next movie night.
            </PageIntro>
            <section className="content-section" aria-label="Movie catalogue">
                {loading ? <MovieGridSkeleton /> : message && movies.length === 0 && message.startsWith('The')
                    ? <ErrorState message={message} onRetry={() => setReloadKey(key => key + 1)} />
                    : <Movies movies={movies} updateMovieReview={updateMovieReview} message={message}/ >}
            </section>
        </div>

    );

};

export default Home;


