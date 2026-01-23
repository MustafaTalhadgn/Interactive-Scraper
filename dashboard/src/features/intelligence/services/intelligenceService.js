import api from '../../../shared/utils/api';

export const intelligenceService = {
  getFeed: async (params) => {
    const response = await api.get('/intelligence', { params });
    return response.data;
  },

  getDetail: async (id) => {
    const response = await api.get(`/intelligence/${id}`);
    return response.data;
  }
};