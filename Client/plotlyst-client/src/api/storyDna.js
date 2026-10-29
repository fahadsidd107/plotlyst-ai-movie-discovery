import axiosClient from './axiosConfig';
const unwrap=response=>response.data?.data??response.data;
export const storyDnaApi={movie:async id=>unwrap(await axiosClient.get(`/api/story-dna/movies/${encodeURIComponent(id)}`)),discover:async payload=>unwrap(await axiosClient.post('/api/story-dna/discover',payload))};
export const storyDnaMessage=error=>error.response?.data?.error?.message||error.response?.data?.message||error.message||'The story map could not be developed.';
