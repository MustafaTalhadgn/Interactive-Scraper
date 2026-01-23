import React from 'react';
import { ExternalLink, Calendar, Shield, Eye } from 'lucide-react';
import CriticalityBadge from '../../../shared/components/CriticalityBadge';
import { formatDate, truncateText } from '../../../shared/utils/helpers'; // truncateText'i buradan çekiyoruz

const IntelligenceList = ({ items, loading, onViewDetail }) => {
  if (loading) {
    return (
      <div className="bg-gray-800 rounded-xl border border-gray-700 p-6 space-y-4">
        {[1, 2, 3, 4, 5].map((i) => (
          <div key={i} className="h-16 bg-gray-700/50 rounded-lg animate-pulse" />
        ))}
      </div>
    );
  }

  if (!items || items.length === 0) {
    return (
      <div className="bg-gray-800 rounded-xl border border-gray-700 p-12 text-center">
        <Shield className="mx-auto h-12 w-12 text-gray-500 mb-4" />
        <h3 className="text-lg font-medium text-white">Veri Bulunamadı</h3>
        <p className="text-gray-400 mt-2">Aradığınız kriterlere uygun tehdit istihbaratı yok.</p>
      </div>
    );
  }

  return (
    <div className="bg-gray-800 rounded-xl border border-gray-700 overflow-hidden shadow-sm">
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-700">
          <thead className="bg-gray-900/50">
            <tr>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                Tehdit Başlığı
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                Kaynak
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                Kritiklik
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                Tarih
              </th>
              <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-400 uppercase tracking-wider">
                İşlem
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-700 bg-gray-800">
            {items.map((item) => (
              <tr 
                key={item.id} 
                className="hover:bg-gray-700/40 transition-colors group cursor-pointer"
                onClick={() => onViewDetail(item.id)} 
              >
                <td className="px-6 py-4">
                  <div className="flex flex-col">
                   
                    <span className="text-sm font-medium text-white group-hover:text-blue-400 transition-colors">
                      {truncateText(item.title, 60)}
                    </span>
                    <span className="text-xs text-gray-500 mt-1 font-mono">ID: #{item.id}</span>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex flex-col">
                    <span className="text-sm text-gray-300">{item.source_name}</span>
                    <span className="text-xs text-gray-500 bg-gray-700/50 px-2 py-0.5 rounded w-fit mt-1 border border-gray-600">
                      {item.category}
                    </span>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <CriticalityBadge level={item.criticality} score={item.criticality_score} />
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center text-sm text-gray-400">
                    <Calendar className="w-4 h-4 mr-2 opacity-70" />
                    {formatDate(item.created_at)}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      onViewDetail(item.id);
                    }}
                    className="text-blue-400 hover:text-blue-300 inline-flex items-center gap-1 transition-colors bg-blue-500/10 px-3 py-1.5 rounded-lg hover:bg-blue-500/20 border border-blue-500/20"
                  >
                    <Eye className="w-4 h-4" /> İncele
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default IntelligenceList;