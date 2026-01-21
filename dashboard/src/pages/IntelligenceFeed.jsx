import React, { useState, useEffect } from 'react';
import IntelligenceFilters from '../components/intelligence/IntelligenceFilters';
import IntelligenceTable from '../components/intelligence/IntelligenceTable';
import { intelligenceService } from '../services/intelligence';

const IntelligenceFeed = () => {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState({
    search: '',
    criticality: 'all',
    source: 'all'
  });

  useEffect(() => {
    const fetchFullFeed = async () => {
      try {
        setLoading(true);
      
        const response = await intelligenceService.getFeed({ limit: 100 });
        if (response.success) {
          setData(response.data.intelligence);
        }
      } catch (err) {
        console.error("Feed yükleme hatası:", err);
      } finally {
        setLoading(false);
      }
    };
    fetchFullFeed();
  }, []);

  
  const filteredData = data.filter(item => {
    const matchesSearch = item.title.toLowerCase().includes(filters.search.toLowerCase());
    const matchesCriticality = filters.criticality === 'all' || 
                               item.criticality_label.toLowerCase() === filters.criticality.toLowerCase();
    const matchesSource = filters.source === 'all' || item.source_name === filters.source;
    
    return matchesSearch && matchesCriticality && matchesSource;
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">İstihbarat Akışı</h1>
        <p className="text-gray-500 text-sm">Toplanan tüm verileri filtreleyin ve derinlemesine analiz edin.</p>
      </div>

      <IntelligenceFilters filters={filters} setFilters={setFilters} />

      {loading ? (
        <div className="flex justify-center py-12 text-blue-600">Veriler yükleniyor...</div>
      ) : (
        <IntelligenceTable 
          data={filteredData} 
          onDetailClick={(id) => console.log("Detay ID:", id)} 
        />
      )}
    </div>
  );
};

export default IntelligenceFeed;