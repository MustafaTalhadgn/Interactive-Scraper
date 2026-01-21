import React from 'react';
import { ExternalLink, ShieldAlert } from 'lucide-react';
import { getCriticalityStyle } from '../../utils/criticality';
import { formatDate } from '../../utils/date';

const IntelligenceTable = ({ data, onDetailClick }) => {
  return (
    <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
      <table className="w-full text-left border-collapse">
        <thead className="bg-gray-50 border-b border-gray-100">
          <tr>
            <th className="px-6 py-4 text-xs font-semibold text-gray-500 uppercase tracking-wider">Tehdit Başlığı</th>
            <th className="px-6 py-4 text-xs font-semibold text-gray-500 uppercase tracking-wider">Kaynak / Kategori</th>
            <th className="px-6 py-4 text-xs font-semibold text-gray-500 uppercase tracking-wider">Kritiklik</th>
            <th className="px-6 py-4 text-xs font-semibold text-gray-500 uppercase tracking-wider">Tarih</th>
            <th className="px-6 py-4 text-xs font-semibold text-gray-500 uppercase tracking-wider text-right">İşlem</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-50">
          {data.map((item) => {
            const style = getCriticalityStyle(item.criticality_score);
            return (
              <tr key={item.id} className="hover:bg-gray-50/50 transition-colors group">
                <td className="px-6 py-4">
                  <div className="flex items-center gap-3">
                    <div className={`p-2 rounded-lg ${style.classes.split(' ')[0]}`}>
                      <ShieldAlert size={18} className={style.classes.split(' ')[1]} />
                    </div>
                    <span className="font-medium text-gray-900 line-clamp-1">{item.title}</span>
                  </div>
                </td>
                <td className="px-6 py-4">
                  <div className="flex flex-col">
                    <span className="text-sm font-medium text-gray-700">{item.source_name}</span>
                    <span className="text-xs text-gray-400">{item.category}</span>
                  </div>
                </td>
                <td className="px-6 py-4">
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${style.classes}`}>
                    <span className={`w-1.5 h-1.5 rounded-full mr-1.5 ${style.dot}`}></span>
                    {style.label}
                  </span>
                </td>
                <td className="px-6 py-4 text-sm text-gray-500 font-mono">
                  {formatDate(item.created_at)}
                </td>
                <td className="px-6 py-4 text-right">
                  <button 
                    onClick={() => onDetailClick(item.id)}
                    className="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all"
                    title="Detayları Gör"
                  >
                    <ExternalLink size={18} />
                  </button>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
      
      {data.length === 0 && (
        <div className="p-12 text-center text-gray-400">
          Henüz analiz edilmiş bir veri bulunmuyor.
        </div>
      )}
    </div>
  );
};

export default IntelligenceTable;