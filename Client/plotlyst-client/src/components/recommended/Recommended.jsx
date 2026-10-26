import useAxiosPrivate from '../../hooks/useAxiosPrivate';
import {useEffect, useState} from 'react';
import Movies from '../movies/Movies';
import PageIntro from '../ui/PageIntro';
import { MovieGridSkeleton, ErrorState } from '../ui/States';
import useAuth from '../../hooks/useAuth';

const Recommended = () => {
    const [movies, setMovies] = useState([]);
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState();
    const axiosPrivate = useAxiosPrivate();
    const { auth } = useAuth();
    const [error, setError] = useState('');
    const [reloadKey, setReloadKey] = useState(0);

    useEffect(() => {
        const fetchRecommendedMovies = async () => {
            setLoading(true);
            setMessage("");
            setError("");

            try{
                const response = await axiosPrivate.get('/recommendedmovies');
                setMovies(response.data);
            } catch (error){
                console.error("Error fetching recommended movies:", error)
                setError('Your recommendations could not be loaded right now.')
            } finally {
                setLoading(false);
            }

        }
        fetchRecommendedMovies();
    }, [axiosPrivate, reloadKey])

    const genres = auth?.favourite_genres?.map(genre => genre.genre_name).filter(Boolean).join(', ');

    return (
        <div className="page-shell">
            <PageIntro eyebrow="The algorithm knows your vibe" title="Your taste has entered the chat." constellation
              meta={genres ? `Currently vibing with ${genres}` : 'Powered by the genres you chose—not random hype'}>
              No generic “popular right now” pile. This is your private after-hours shelf: fewer choices, stronger energy, and films selected around the stories you actually enjoy.
            </PageIntro>
            <section className="content-section" aria-label="Recommended movies">
              {loading ? <MovieGridSkeleton count={4} /> : error ? <ErrorState message={error} onRetry={() => setReloadKey(key => key + 1)} /> : <Movies movies={movies} message={message} />}
            </section>
        </div>
    )

}
export default Recommended
