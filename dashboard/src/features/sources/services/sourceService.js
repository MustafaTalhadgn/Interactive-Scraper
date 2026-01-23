import api from '../../../shared/utils/api';

export const sourceService = {
  getAll: async () => {
    const response = await api.get('/sources');
    return response.data;
  },

  create: async (data) => {
    const response = await api.post('/sources', data);
    return response.data;
  },

  update: async (id, data) => {
    const response = await api.patch(`/sources/${id}`, data);
    return response.data;
  },

  delete: async (id) => {
    const response = await api.delete(`/sources/${id}`);
    return response.data;
  },

  triggerScrape: async (id) => {
    const response = await api.post(`/sources/${id}/scrape`);
    return response.data;
  }
};