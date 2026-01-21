import React from 'react';
import { Search, Filter } from 'lucide-react';

const IntelligenceFilters = ({ filters, setFilters }) => {
  return (
    <div className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-wrap gap-4 items-center">
   
      <div className="flex-1 min-w-[200px] relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
        <input
          type="text"
          placeholder="Tehdit başlığı veya içerik ara..."
          className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
          value={filters.search}
          onChange={(e) => setFilters({ ...filters, search: e.target.value })}
        />
      </div>


      <select
        className="px-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white cursor-pointer"
        value={filters.criticality}
        onChange={(e) => setFilters({ ...filters, criticality: e.target.value })}
      >
        <option value="all">Tüm Kritiklik Seviyeleri</option>
        <option value="critical">Critical</option>
        <option value="high">High</option>
        <option value="medium">Medium</option>
        <option value="low">Low</option>
      </select>


      <select
        className="px-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white cursor-pointer"
        value={filters.source}
        onChange={(e) => setFilters({ ...filters, source: e.target.value })}
      >
        <option value="all">Tüm Kaynaklar</option>
        <option value="Shadow Wiki">Shadow Wiki</option>
        <option value="Dread">Dread</option>
        <option value="Tor News">Tor News</option>
      </select>
    </div>
  );
};

export default IntelligenceFilters;