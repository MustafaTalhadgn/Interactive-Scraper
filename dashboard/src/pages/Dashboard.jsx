import React, { useEffect, useState } from 'react';

import StatCard from '../components/common/StatCard';
import IntelligenceTable from '../components/intelligence/IntelligenceTable';
import { intelligenceService } from '../services/intelligence';
import { ShieldAlert, Activity, Globe, Zap } from 'lucide-react';

const Dashboard = () => {
  const [feed, setFeed] = useState([]);
  const [stats, setStats] = useState({ total: 0, critical: 0, high: 0, sources: 0 });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
    
        const [feedRes, statsRes] = await Promise.all([
          intelligenceService.getFeed({ limit: 10 }),
          intelligenceService.getStats()
        ]);

        if (feedRes.success) setFeed(feedRes.data.intelligence);
        if (statsRes.success) setStats(statsRes.data);
      } catch (err) {
        console.error("Veri yükleme hatası:", err);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  return (
   
      <div className="space-y-8">
       
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Siber Tehdit İstihbaratı Paneli</h1>
          <p className="text-gray-500 text-sm">Dark Web üzerinden anlık toplanan verilerin analizi.</p>
        </div>

       
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <StatCard title="Toplam İstihbarat" value={stats.total} icon={<Activity size={24}/>} colorClass="bg-blue-50 text-blue-600" />
          <StatCard title="Kritik Tehditler" value={stats.critical} icon={<ShieldAlert size={24}/>} colorClass="bg-red-50 text-red-600" />
          <StatCard title="Aktif Kaynaklar" value={stats.sources} icon={<Globe size={24}/>} colorClass="bg-emerald-50 text-emerald-600" />
          <StatCard title="Son 24 Saat" value={stats.high} icon={<Zap size={24}/>} colorClass="bg-orange-50 text-orange-600" />
        </div>

   
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h2 className="text-lg font-semibold text-gray-800">Son Tehdit Akışı</h2>
            <button className="text-blue-600 text-sm font-medium hover:underline">Tümünü Gör</button>
          </div>
          <IntelligenceTable data={feed} onDetailClick={(id) => console.log("Detay:", id)} />
        </div>
      </div>
   
  );
};

export default Dashboard;