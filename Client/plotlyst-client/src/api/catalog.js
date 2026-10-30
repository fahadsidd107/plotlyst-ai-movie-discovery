import axiosClient from './axiosConfig';

export const imageUrl = (meta, path, size = 'w500') => {
  if (!path) return null;
  const base = (meta?.image_base_url || 'https://image.tmdb.org/t/p').replace(/\/$/, '');
  return `${base}/${size}${path}`;
};
export const getMovies = (page = 1, collection = 'trending') => axiosClient.get('/api/catalog/movies', { params: { page, collection } });
export const getMovie = (id, reveal = false) => axiosClient.get(`/api/catalog/movies/${id}`, { params: reveal ? { reveal_spoilers: true } : {} });
export const getTV = (page = 1, collection = 'trending') => axiosClient.get('/api/catalog/tv', { params: { page, collection } });
export const getTVDetail = (id, reveal = false) => axiosClient.get(`/api/catalog/tv/${id}`, { params: reveal ? { reveal_spoilers: true } : {} });
export const getSeason = (id, season, reveal = false) => axiosClient.get(`/api/catalog/tv/${id}/seasons/${season}`, { params: reveal ? { reveal_spoilers: true } : {} });
export const getEpisode = (id, season, episode, reveal = false) => axiosClient.get(`/api/catalog/tv/${id}/seasons/${season}/episodes/${episode}`, { params: reveal ? { reveal_spoilers: true } : {} });
export const getPeople = (page = 1) => axiosClient.get('/api/catalog/people/trending', { params: { page } });
export const getBornToday = () => axiosClient.get('/api/catalog/people/born-today', { params: { timezone: 'Europe/Berlin' } });
export const getPerson = (id) => axiosClient.get(`/api/catalog/people/${id}`);
export const searchMulti = (query, page = 1, signal) => axiosClient.get('/api/discovery/search', { params: { q: query, page }, signal });
export const getDiscoveryHome = (signal) => axiosClient.get('/api/discovery/home', { signal });
export const getHomeSection = (section, signal) => axiosClient.get(`/api/discovery/home/${section}`, { signal });
export const discover = (media, params, signal) => {
  const endpoint = media === 'movie' ? 'movies' : media;
  return axiosClient.get(`/api/discovery/${endpoint}`, { params, signal });
};
export const getDiscoveryFilters = (media, signal) => axiosClient.get(`/api/discovery/filters/${media}`, { signal });
export const getChart = (media, size, page = 1, signal) => axiosClient.get(`/api/discovery/charts/${media}/${size}`, { params: { page }, signal });
export const getFranchise = (id, signal) => axiosClient.get(`/api/discovery/franchises/${id}`, { signal });
export const getCollectionKinds = (signal) => axiosClient.get('/api/discovery/collection-kinds', { signal });
