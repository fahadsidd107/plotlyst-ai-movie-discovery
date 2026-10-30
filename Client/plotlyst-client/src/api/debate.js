import axiosClient from './axiosConfig';

const unwrap = response => response.data?.data ?? response.data;

export const debateApi = {
  movie: async (id, ratingBucket = 'all') => unwrap(await axiosClient.get(`/api/debate/movies/${encodeURIComponent(id)}`, {
    params: { rating_bucket: ratingBucket },
  })),
};

export const debateMessage = error => error.response?.data?.error?.message
  || error.response?.data?.message
  || error.message
  || 'The audience ballot could not be counted.';
