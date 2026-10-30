import axiosClient from './axiosConfig';
const unwrap=response=>response.data?.data??response.data;
export const storyDnaApi={movie:async(id,includeSpoilers=false)=>unwrap(await axiosClient.get(`/api/story-dna/movies/${encodeURIComponent(id)}`,{params:includeSpoilers?{include_spoilers:true}:{}})),discover:async payload=>{const response=await axiosClient.post('/api/story-dna/discover',payload);return {...unwrap(response),provenance:response.data?.provenance,limitations:response.data?.limitations||[]}}};
export const storyDnaMessage=error=>error.response?.data?.error?.message||error.response?.data?.message||error.message||'The story map could not be developed.';
