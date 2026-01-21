import api from './api';

export const intelligenceService = {

  getFeed: (params) => {
    return api.get('/intelligence', { params });
  },

  getDetail: (id) => {
    return api.get(`/intelligence/${id}`);
  },
};