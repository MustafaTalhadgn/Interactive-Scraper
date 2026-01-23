import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import useApi from '../../shared/hooks/useApi';
import { statsService } from './services/statsService';
import { intelligenceService } from '../intelligence/services/intelligenceService'; 


import StatsCards from './components/StatsCards';
import TimelineChart from './components/TimelineChart';
import CategoryPieChart from './components/CategoryPieChart';
import IntelligenceList from '../intelligence/components/IntelligenceList';     
import IntelligenceDetail from '../intelligence/components/IntelligenceDetail'; 
import Button from '../../shared/components/Button';
import { RefreshCw, ArrowRight } from 'lucide-react';

const OverviewPage = () => {
  const navigate = useNavigate();


  const statsApi = useApi(statsService.getOverview);
  const timelineApi = useApi(statsService.getTimeline);
  const recentFeedApi = useApi(intelligenceService.getFeed); 


  const [selectedId, setSelectedId] = useState(null);
  const [isDetailOpen, setIsDetailOpen] = useState(false);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
   
    await Promise.all([
      statsApi.execute(),
      timelineApi.execute(),
      recentFeedApi.execute({ limit: 5, page: 1 }) 
    ]);
  };


  const handleViewDetail = (id) => {
    setSelectedId(id);
    setIsDetailOpen(true);
  };

  const handleCloseDetail = () => {
    setIsDetailOpen(false);
    setSelectedId(null);
  };

  const isLoading = statsApi.loading || timelineApi.loading || recentFeedApi.loading;

  return (
    <div className="space-y-6">
 
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white">Siber Tehdit İstihbaratı Paneli</h1>
          <p className="text-gray-400 mt-1">Dark Web üzerinden toplanan verilerin anlık analizi.</p>
        </div>
        <Button 
          variant="secondary" 
          icon={RefreshCw} 
          onClick={fetchData} 
          loading={isLoading}
        >
          Yenile
        </Button>
      </div>


      <StatsCards data={statsApi.data?.data} loading={isLoading} />


      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">

        <div className="lg:col-span-2 min-h-[400px]">
          <TimelineChart data={timelineApi.data?.data?.timeline} />
        </div>

        <div className="h-full min-h-[400px]">
          <CategoryPieChart data={statsApi.data?.data?.category_distribution} />
        </div>
      </div>
      

      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold text-white">Son Tehdit Akışı</h3>
          <button 
            onClick={() => navigate('/intelligence')}
            className="text-sm text-blue-400 hover:text-blue-300 flex items-center gap-1 transition-colors"
          >
            Tümünü Gör <ArrowRight className="w-4 h-4" />
          </button>
        </div>

        <div className="bg-gray-800 rounded-xl border border-gray-700 overflow-hidden">
          <IntelligenceList 
            items={recentFeedApi.data?.data?.intelligence} 
            loading={recentFeedApi.loading}
            onViewDetail={handleViewDetail}
          />
        </div>
      </div>

      <IntelligenceDetail 
        id={selectedId}
        isOpen={isDetailOpen}
        onClose={handleCloseDetail}
        size="2xl"
      />
    </div>
  );
};

export default OverviewPage;