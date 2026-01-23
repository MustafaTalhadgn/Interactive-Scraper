import React from 'react';
import { Search, Filter, X } from 'lucide-react';
import { SOURCE_CATEGORIES, CRITICALITY_LEVELS } from '../../../shared/utils/constants';

const IntelligenceFilters = ({ filters, onChange, onClear }) => {
  const handleChange = (key, value) => {
    onChange({ ...filters, [key]: value, page: 1 }); 
  };

  return (
    <div className="bg-gray-800 p-4 rounded-xl border border-gray-700 mb-6 space-y-4 md:space-y-0 md:flex md:items-center md:gap-4">
      
   
      <div className="flex-1 relative">
        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <Search className="h-5 w-5 text-gray-400" />
        </div>
        <input
          type="text"
          className="block w-full pl-10 pr-3 py-2 border border-gray-600 rounded-lg leading-5 bg-gray-900 text-gray-300 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm transition duration-150 ease-in-out"
          placeholder="Başlık veya içerikte ara..."
          value={filters.search || ''}
          onChange={(e) => handleChange('search', e.target.value)}
        />
      </div>

      <div className="flex gap-2 overflow-x-auto pb-2 md:pb-0">
        <select
          value={filters.criticality || ''}
          onChange={(e) => handleChange('criticality', e.target.value)}
          className="block w-40 pl-3 pr-10 py-2 text-base border-gray-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-lg bg-gray-900 text-gray-300"
        >
          <option value="">Tüm Kritiklikler</option>
          {Object.entries(CRITICALITY_LEVELS).map(([key, val]) => (
            <option key={key} value={val}>{key}</option>
          ))}
        </select>

        <select
          value={filters.category || ''}
          onChange={(e) => handleChange('category', e.target.value)}
          className="block w-40 pl-3 pr-10 py-2 text-base border-gray-600 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-lg bg-gray-900 text-gray-300"
        >
          <option value="">Tüm Kategoriler</option>
          {SOURCE_CATEGORIES.map((cat) => (
            <option key={cat} value={cat}>{cat.toUpperCase()}</option>
          ))}
        </select>


        {(filters.search || filters.criticality || filters.category) && (
          <button
            onClick={onClear}
            className="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-red-400 bg-red-900/20 hover:bg-red-900/40 focus:outline-none transition"
          >
            <X className="h-4 w-4 mr-1" />
            Temizle
          </button>
        )}
      </div>
    </div>
  );
};

export default IntelligenceFilters;