import axiosClient from './axiosConfig';

const unwrap = response => response.data?.data ?? response.data;
const request = async (method, url, data) => unwrap(await axiosClient({ method, url, data }));

export const intelligenceApi = {
  discover: payload => request('post', '/api/ai/discover', payload),
  memorySearch: payload => request('post', '/api/ai/memory-search', payload),
  cinemaDna: () => request('get', '/api/analytics/cinema-dna'),
  createRoom: payload => request('post', '/api/movie-night/rooms', payload),
  joinRoom: (code, payload) => request('post', `/api/movie-night/rooms/${encodeURIComponent(code)}/join`, payload),
  room: async code => (await axiosClient.get(`/api/movie-night/rooms/${encodeURIComponent(code)}`)).data,
  preferences: (code, payload) => request('put', `/api/movie-night/rooms/${encodeURIComponent(code)}/preferences`, payload),
};

export const intelligenceMessage = error => error.response?.data?.error?.message || error.response?.data?.message || error.message || 'The projection booth is temporarily unavailable.';

export const movieNightSocketUrl = code => {
  const base = import.meta.env.VITE_API_BASE_URL || window.location.origin;
  const url = new URL(`/api/movie-night/rooms/${encodeURIComponent(code)}/ws`, base);
  url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:';
  return url.toString();
};
