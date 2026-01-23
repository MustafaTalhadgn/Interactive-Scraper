import React from 'react';
import { Database, ShieldAlert, Activity, Globe } from 'lucide-react';
import { clsx } from 'clsx';

const Card = ({ title, value, icon: Icon, color, trend }) => {
  const colorClasses = {
    blue: 'bg-blue-500/10 text-blue-500',
    red: 'bg-red-500/10 text-red-500',
    green: 'bg-green-500/10 text-green-500',
    orange: 'bg-orange-500/10 text-orange-500',
  };

  return (
    <div className="bg-gray-800 rounded-xl border border-gray-700 p-6">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-gray-400 uppercase tracking-wider">{title}</p>
          <p className="mt-2 text-3xl font-bold text-white">{value}</p>
        </div>
        <div className={clsx("p-3 rounded-lg", colorClasses[color])}>
          <Icon className="w-6 h-6" />
        </div>
      </div>
      
      {trend && (
         <div className="mt-4 flex items-center text-sm">
           <span className="text-green-400 font-medium">+{trend}%</span>
           <span className="text-gray-500 ml-2">geçen haftaya göre</span>
         </div>
      )}
    </div>
  );
};

const StatsCards = ({ data, loading }) => {
  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {[1, 2, 3, 4].map((i) => (
          <div key={i} className="h-32 bg-gray-800/50 rounded-xl animate-pulse border border-gray-700/50" />
        ))}
      </div>
    );
  }

  const stats = data || { 
    total_intelligence: 0, 
    critical_count: 0, 
    active_sources: 0, 
    last_24_hours: 0 
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <Card 
        title="Toplam İstihbarat" 
        value={stats.total_intelligence} 
        icon={Database} 
        color="blue" 
      />
      <Card 
        title="Kritik Tehditler" 
        value={stats.critical_count} 
        icon={ShieldAlert} 
        color="red" 
      />
      <Card 
        title="Aktif Kaynaklar" 
        value={stats.active_sources} 
        icon={Globe} 
        color="green" 
      />
      <Card 
        title="Son 24 Saat" 
        value={stats.last_24_hours} 
        icon={Activity} 
        color="orange" 
      />
    </div>
  );
};

export default StatsCards;