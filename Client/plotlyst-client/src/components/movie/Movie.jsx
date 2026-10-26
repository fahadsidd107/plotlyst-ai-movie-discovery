import Button from 'react-bootstrap/Button'
import "./Movie.css";
const Movie = ({movie,updateMovieReview}) => {
    return (
        <article className="movie-card" key={movie._id}>
                <div className="movie-card__poster-wrap">
                    <img src={movie.poster_path} alt={movie.title} 
                        className="movie-card__poster" loading="lazy"
                    />
                    {movie.ranking?.ranking_name && <span className="ranking-badge">{movie.ranking.ranking_name}</span>}
                </div>
                <div className="movie-card__body">
                    <p className="movie-card__id">IMDb · {movie.imdb_id}</p>
                    <h2 className="movie-card__title">{movie.title}</h2>
                    {movie.admin_review && <p className="movie-card__review">“{movie.admin_review}”</p>}
                </div>
                  {updateMovieReview && (
                        <Button
                            className="button button--card"
                            onClick={e => {
                                e.preventDefault();
                                updateMovieReview(movie.imdb_id);
                            }}
                        >
                            Write review <span aria-hidden="true">→</span>
                        </Button>
                    )}
        </article>
    )
}
export default Movie;
