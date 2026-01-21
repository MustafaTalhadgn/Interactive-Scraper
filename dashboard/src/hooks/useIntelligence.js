import { useState, useEffect, useCallback } from 'react';
import { intelligenceService } from '../services/intelligence';

export const useIntelligence = () => {
  const [data, setData] = useState([]);
  const [stats, setStats] = useState({ total: 0, critical: 0, high: 0, sources: 0 });
  const [loading, setLoading] = useState(true);

  const refreshData = useCallback(async () => {
    try {
      setLoading(true);
      const [feedRes, statsRes] = await Promise.all([
        intelligenceService.getFeed({ limit: 20 }),
        intelligenceService.getStats()
      ]);

      if (feedRes.success) setData(feedRes.data.intelligence);
      if (statsRes.success) setStats(statsRes.data);
    } catch (err) {
      console.error("Veri çekme hatası:", err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refreshData();
  }, [refreshData]);

  return { data, stats, loading, refreshData };
};