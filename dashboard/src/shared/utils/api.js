import axios from 'axios';
import { API_BASE_URL } from './constants';


const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  },
  timeout: 30000
});


api.interceptors.request.use(
  (config) => {

    const token = localStorage.getItem('cti_auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);


api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {

    if (error.response) {
      const { status, data } = error.response;
      
      switch (status) {
        case 401:

          console.error('Unauthorized access');
          break;
        case 403:
          console.error('Forbidden access');
          break;
        case 404:
          console.error('Resource not found');
          break;
        case 500:
          console.error('Server error');
          break;
        default:
          console.error('API Error:', data?.error?.message || 'Unknown error');
      }
    } else if (error.request) {
      console.error('Network error - no response received');
    } else {
      console.error('Error:', error.message);
    }
    
    return Promise.reject(error);
  }
);



export const getIntelligenceFeed = async (params = {}) => {
  const response = await api.get('/intelligence', { params });
  return response.data;
};

export const getIntelligenceDetail = async (id) => {
  const response = await api.get(`/intelligence/${id}`);
  return response.data;
};


export const getStatsOverview = async () => {
  const response = await api.get('/stats/overview');
  return response.data;
};


export const getStatsTimeline = async (days = 7) => {
  const response = await api.get('/stats/timeline', { params: { days } });
  return response.data;
};


export const getSources = async () => {
  const response = await api.get('/sources');
  return response.data;
};

export const createSource = async (sourceData) => {
  const response = await api.post('/sources', sourceData);
  return response.data;
};


export const updateSource = async (id, updateData) => {
  const response = await api.patch(`/sources/${id}`, updateData);
  return response.data;
};


export const deleteSource = async (id) => {
  const response = await api.delete(`/sources/${id}`);
  return response.data;
};

export const triggerScrape = async (id) => {
  const response = await api.post(`/sources/${id}/scrape`);
  return response.data;
};


export const healthCheck = async () => {
  const response = await api.get('/health');
  return response.data;
};

export default api;