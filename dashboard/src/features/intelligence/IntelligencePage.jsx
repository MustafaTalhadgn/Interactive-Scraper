import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import useApi from '../../shared/hooks/useApi';
import { intelligenceService } from './services/intelligenceService';
import IntelligenceFilters from './components/IntelligenceFilters';
import IntelligenceList from './components/IntelligenceList';
import Pagination from '../../shared/components/Pagination';
import IntelligenceDetail from './components/IntelligenceDetail'

const IntelligencePage = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [selectedId, setSelectedId] = useState(null);
  const [isDetailOpen, setIsDetailOpen] = useState(false);

  
  const [filters, setFilters] = useState({
    page: parseInt(searchParams.get('page') || '1'),
    limit: 20,
    search: searchParams.get('search') || '',
    criticality: searchParams.get('criticality') || '',
    category: searchParams.get('category') || ''
  });

  const { data, loading, execute } = useApi(intelligenceService.getFeed);

 
  useEffect(() => {
    execute(filters);
    
  
    const params = {};
    if (filters.page > 1) params.page = filters.page;
    if (filters.search) params.search = filters.search;
    if (filters.criticality) params.criticality = filters.criticality;
    if (filters.category) params.category = filters.category;
    setSearchParams(params);
    
  }, [filters, setSearchParams]); 

  const handleFilterChange = (newFilters) => {
    setFilters(newFilters);
  };

  const handlePageChange = (newPage) => {
    setFilters(prev => ({ ...prev, page: newPage }));
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleClearFilters = () => {
    setFilters({
      page: 1,
      limit: 20,
      search: '',
      criticality: '',
      category: ''
    });
  };

  const handleViewDetail = (id) => {
      setSelectedId(id);
      setIsDetailOpen(true);
    };
  
  const handleCloseDetail = () => {
    setIsDetailOpen(false);
    setSelectedId(null);
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">İstihbarat Akışı</h1>
        <p className="text-gray-400 mt-1">Dark Web kaynaklarından toplanan ve analiz edilen veriler.</p>
      </div>

      <IntelligenceFilters 
        filters={filters} 
        onChange={handleFilterChange}
        onClear={handleClearFilters}
      />

      <IntelligenceList 
        items={data?.data?.intelligence} 
        loading={loading} 
        onViewDetail={handleViewDetail}
      />

      {data?.data?.pagination && (
        <Pagination
          currentPage={data.data.pagination.page}
          totalPages={data.data.pagination.total_pages}
          onPageChange={handlePageChange}
        />
      )}
      <IntelligenceDetail 
        id={selectedId}
        isOpen={isDetailOpen}
        onClose={handleCloseDetail}
      />
    </div>
  );
};

export default IntelligencePage;