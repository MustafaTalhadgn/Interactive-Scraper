import axios from 'axios';


const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.response.use(
  (response) => response.data, 
  (error) => {
    const message = error.response?.data?.error?.message || 'Bir ağ hatası oluştu';
    console.error('API Hatası:', message);
    return Promise.reject(error);
  }
);

export default api;