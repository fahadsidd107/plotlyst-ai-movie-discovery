import axiosClient from './axiosConfig';

const unwrap = response => {
  const payload = response.data;
  if (!payload || !Object.prototype.hasOwnProperty.call(payload, 'data')) return payload;
  if (Array.isArray(payload.data)) return { results: payload.data, source: payload.source };
  if (payload.data && typeof payload.data === 'object') return { ...payload.data, source: payload.source || payload.data.source };
  return { results: [], source: payload.source };
};
const get = async (url, params, signal) => unwrap(await axiosClient.get(url, { params, signal }));

export const specialistApi = {
  providers: signal => get('/api/specialist/providers', undefined, signal),
  availability: async ({ media_type, media_id, region }, signal) => {
    const result = await get(`/api/specialist/availability/${media_type}/${media_id}`, { region }, signal);
    const names = { subscription: 'flatrate', rental: 'rent', purchase: 'buy' };
    const grouped = (result.offers || []).reduce((offers, item) => {
      const key = names[item.type] || item.type;
      return { ...offers, [key]: [...(offers[key] || []), { ...item, link: result.link }] };
    }, {});
    return { ...result, offers: grouped };
  },
  awardOrganizations: signal => get('/api/specialist/awards/organizations', undefined, signal),
  awards: (organization, year, signal) => get(`/api/specialist/awards/${encodeURIComponent(organization)}/${year}`, undefined, signal),
  boxOffice: (period, date, signal) => get('/api/specialist/box-office', { period, date }, signal),
  editorial: signal => get('/api/specialist/editorial', undefined, signal),
};
