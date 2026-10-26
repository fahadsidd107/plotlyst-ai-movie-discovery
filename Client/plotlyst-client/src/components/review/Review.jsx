import {Form, Button} from 'react-bootstrap';
import { useRef,useEffect,useState } from 'react';
import {useParams} from 'react-router-dom'
//import axiosPrivate from '../../api/axiosPrivateConfig';
import useAxiosPrivate from '../../hooks/useAxiosPrivate';
import useAuth from '../../hooks/useAuth';
import Movie from '../movie/Movie'
import { MovieGridSkeleton } from '../ui/States';
import PageIntro from '../ui/PageIntro';

const Review = () => {

    const [movie, setMovie] = useState({});
    const [loading, setLoading] = useState(false);
    const revText = useRef();
    const { imdb_id } = useParams();
    const {auth} = useAuth();
    const axiosPrivate = useAxiosPrivate();

    useEffect(() => {
        const fetchMovie = async () => {
            setLoading(true);
            try {
                const response = await axiosPrivate.get(`/movie/${imdb_id}`);
                setMovie(response.data);
                console.log(response.data);
            } catch (error) {
                console.error('Error fetching movie:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchMovie();

    }, [axiosPrivate, imdb_id]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        
        setLoading(true);
        try {
            
            const response = await axiosPrivate.patch(`/updatereview/${imdb_id}`, { admin_review: revText.current.value });
            console.log(response.data);           

            setMovie(currentMovie => ({
                ...currentMovie,
                admin_review: response.data?.admin_review ?? currentMovie.admin_review,
                ranking: {
                    ranking_name: response.data?.ranking_name ?? currentMovie.ranking?.ranking_name
                }
            }));

        } catch (err) {
            console.error(err);
            if (err.response && err.response.status === 401) {
                 console.error('Unauthorized access - redirecting to login');
                 localStorage.removeItem('user');
                // setAuth(null);
            } else {
                console.error('Error updating review:', err);
            }

        } finally {
            setLoading(false);
        }
    }; 

    return (
      <>
        {loading ? (
            <div className="page-shell"><MovieGridSkeleton count={1} /></div>
        ) : (
            <div className="page-shell review-page">
                <PageIntro eyebrow="Editorial desk" title="Review the selection.">Add a concise editorial note; Plotlyst will update the title’s ranking after submission.</PageIntro>
                <div className="review-layout">
                    <div className="review-movie">
                            <Movie movie={movie} />
                    </div>
                    <div className="review-editor">
                            {auth && auth.role === "ADMIN" ? (
                                <Form onSubmit={handleSubmit}>
                                    <Form.Group className="mb-3" controlId="adminReviewTextarea">
                                        <Form.Label>Admin Review</Form.Label>
                                        <Form.Control
                                            ref={revText}
                                            required
                                            as="textarea"
                                            rows={8}
                                            defaultValue={movie?.admin_review}
                                            placeholder="Write your review here..."
                                            style={{ resize: "vertical" }}
                                        />
                                    </Form.Group>
                                    <div className="d-flex justify-content-end">
                                        <Button className="button button--primary" type="submit">
                                            Submit Review
                                        </Button>
                                    </div>
                                </Form> ):(
                                <div className="review-readonly">{movie.admin_review || 'No editorial review has been added yet.'}</div>
                            )}                           
                    </div>
                </div>
            </div>
        )}
    </>      

    );
}

export default Review;
