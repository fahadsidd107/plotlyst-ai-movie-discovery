import Movie from '../movie/Movie'
import { EmptyState } from '../ui/States'

const Movies = ({movies,updateMovieReview, message}) => {

    return (
        <div className="movie-grid">
                {movies && movies.length > 0
                    ? movies.map((movie) => (
                        <Movie key={movie._id} updateMovieReview={updateMovieReview} movie={movie} />
                    ))
                    : <EmptyState message={message || 'No movies match this collection yet.'} />
                }

        </div>
    )
}
export default Movies;
