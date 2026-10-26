import axiosClient from './axiosConfig';

export const imageUrl = (meta, path, size = 'w500') => {
  if (!path) return null;
  const base = (meta?.image_base_url || 'https://image.tmdb.org/t/p').replace(/\/$/, '');
  return `${base}/${size}${path}`;
};
export const getMovies = (page = 1, collection = 'trending') => axiosClient.get('/api/catalog/movies', { params: { page, collection } });
export const getMovie = (id) => axiosClient.get(`/api/catalog/movies/${id}`);
export const getTV = (page = 1, collection = 'trending') => axiosClient.get('/api/catalog/tv', { params: { page, collection } });
export const getTVDetail = (id) => axiosClient.get(`/api/catalog/tv/${id}`);
export const getSeason = (id, season) => axiosClient.get(`/api/catalog/tv/${id}/seasons/${season}`);
export const getEpisode = (id, season, episode) => axiosClient.get(`/api/catalog/tv/${id}/seasons/${season}/episodes/${episode}`);
export const getPeople = (page = 1) => axiosClient.get('/api/catalog/people/trending', { params: { page } });
export const getBornToday = () => axiosClient.get('/api/catalog/people/born-today', { params: { timezone: 'Europe/Berlin' } });
export const getPerson = (id) => axiosClient.get(`/api/catalog/people/${id}`);
