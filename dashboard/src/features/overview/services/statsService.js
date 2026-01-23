import api from '../../../shared/utils/api';

export const statsService = {
  getOverview: async () => {
    const response = await api.get('/stats/overview');
    return response.data;
  },

  getTimeline: async (days = 7) => {
    const response = await api.get('/stats/timeline', { params: { days } });
    return response.data;
  }
};