import axiosClient from './axiosConfig';

const unwrap = response => response.data?.data ?? response.data;
const request = async (method, url, data, config) => unwrap(await axiosClient({ method, url, data, ...config }));

export const authApi = {
  session: () => request('get', '/api/auth/session'),
  login: credentials => request('post', '/api/auth/login', credentials),
  register: payload => request('post', '/api/auth/register', payload),
  refresh: () => request('post', '/api/auth/refresh'),
  logout: () => request('post', '/api/auth/logout'),
  logoutAll: () => request('post', '/api/auth/logout-all'),
};

export const profileApi = {
  get: username => request('get', `/api/profiles/${encodeURIComponent(username)}`),
  me: () => request('get', '/api/me/profile'),
  update: payload => request('put', '/api/me/profile', payload),
  settings: () => request('get', '/api/me/profile'),
  updateSettings: payload => request('put', '/api/me/settings', payload),
  deleteAccount: confirmation => request('delete', '/api/me/account', { confirmation }),
};

export const communityApi = {
  ratings: (type, id) => request('get', `/api/community/${type}/${id}/ratings`),
  reviews: (type, id) => request('get', `/api/community/${type}/${id}/reviews`),
  rate: (type, id, rating) => request('put', `/api/community/${type}/${id}/rating`, { stars: rating }),
  removeRating: (type, id) => request('delete', `/api/community/${type}/${id}/rating`),
  createReview: (type, id, payload) => request('post', `/api/community/${type}/${id}/reviews`, payload),
  updateReview: (id, payload) => request('patch', `/api/reviews/${id}`, payload),
  deleteReview: id => request('delete', `/api/reviews/${id}`),
  reply: (id, body) => request('post', `/api/reviews/${id}/replies`, { body }),
  helpful: id => request('put', `/api/reviews/${id}/helpful`),
  report: (id, reason) => request('post', `/api/reviews/${id}/reports`, { reason }),
};

export const libraryApi = {
  watchlist: () => request('get', '/api/me/watchlist'),
  watchlistAdd: item => request('put', `/api/me/watchlist/${item.media_type}/${item.media_id}`, item),
  watchlistRemove: (type, id) => request('delete', `/api/me/watchlist/${type}/${id}`),
  playlists: () => request('get', '/api/me/playlists'),
  playlist: id => request('get', `/api/me/playlists/${id}`),
  createPlaylist: payload => request('post', '/api/me/playlists', payload),
  updatePlaylist: (id, payload) => request('patch', `/api/me/playlists/${id}`, payload),
  history: () => request('get', '/api/me/history'),
  progress: showId => request('get', `/api/me/tv/${showId}/progress`),
  saveProgress: (showId, season, episode, payload) => request('put', `/api/me/tv/${showId}/seasons/${season}/episodes/${episode}/progress`, payload),
};

export const moderationApi = {
  reports: status => request('get', '/api/moderation/reports', null, { params: status ? { status } : {} }),
  review: (id, payload) => request('patch', `/api/moderation/reviews/${id}`, payload),
  suspend: (id, payload) => request('post', `/api/moderation/users/${id}/suspend`, payload),
};

export const apiMessage = error => error.response?.data?.error?.message || error.response?.data?.message || error.message || 'The archive desk is temporarily unavailable.';
