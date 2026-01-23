import { useState, useCallback } from 'react';


const useApi = (apiFunc) => {
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  const execute = useCallback(async (...args) => {
    try {
      setLoading(true);
      setError(null);
      const response = await apiFunc(...args);
      setData(response);
      return response;
    } catch (err) {
      const errorMessage = err.response?.data?.error?.message || err.message || 'Bir hata oluştu';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [apiFunc]);

  return { data, error, loading, execute, setData };
};

export default useApi;