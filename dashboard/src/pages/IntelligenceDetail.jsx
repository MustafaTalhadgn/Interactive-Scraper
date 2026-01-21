import React from 'react';
import { Shield, Globe, Clock, Tag } from 'lucide-react';

const IntelligenceDetail = ({ item, onBack }) => {
  return (
    <div className="bg-white rounded-xl shadow-lg p-6 border border-gray-100">
      <button onClick={onBack} className="text-blue-600 mb-4 hover:underline text-sm">← Listeye Dön</button>
      <div className="flex justify-between items-start mb-6">
        <h2 className="text-2xl font-bold text-gray-900">{item.title}</h2>
        <span className="px-4 py-1 rounded-full bg-red-100 text-red-700 text-sm font-bold">Skor: {item.criticality_score}</span>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
        <div className="flex items-center gap-2 text-gray-600"><Globe size={18}/> {item.source_name}</div>
        <div className="flex items-center gap-2 text-gray-600"><Tag size={18}/> {item.category}</div>
        <div className="flex items-center gap-2 text-gray-600"><Clock size={18}/> {new Date(item.created_at).toLocaleString()}</div>
      </div>

      <div className="bg-slate-50 p-6 rounded-lg border border-slate-200">
        <h3 className="text-lg font-semibold mb-3">Analiz Özeti</h3>
        <p className="text-gray-700 leading-relaxed">{item.summary || "Bu kayıt için otomatik özet oluşturulmadı."}</p>
      </div>
    </div>
  );
};

