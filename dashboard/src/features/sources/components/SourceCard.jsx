import React from 'react';
import { Globe, Clock, Play, Edit2, Trash2, AlertTriangle, ShieldCheck } from 'lucide-react';
import Button from '../../../shared/components/Button';
import { formatDate, formatRelativeTime } from '../../../shared/utils/helpers';
import { CRITICALITY_COLORS } from '../../../shared/utils/constants';

const SourceCard = ({ source, onEdit, onDelete, onTrigger, isTriggering }) => {
  const criticalityColor = CRITICALITY_COLORS[source.criticality] || '#gray';
  
  return (
    <div className="bg-gray-800 rounded-xl border border-gray-700 p-5 hover:border-gray-600 transition-all shadow-lg group">
        <div className="flex justify-between items-start mb-4">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-gray-700/50 rounded-lg">
            <Globe className="w-6 h-6 text-blue-400" />
          </div>
          <div>
            <h3 className="font-semibold text-white text-lg">{source.name}</h3>
            <span className="text-xs text-gray-500 bg-gray-700 px-2 py-0.5 rounded-full uppercase">
              {source.category}
            </span>
          </div>
        </div>
        <div 
          className="flex items-center gap-1 px-2 py-1 rounded text-xs font-bold border"
          style={{ 
            borderColor: criticalityColor, 
            color: criticalityColor,
            backgroundColor: `${criticalityColor}10`
          }}
        >
          <AlertTriangle className="w-3 h-3" />
          {source.criticality.toUpperCase()}
        </div>
      </div>

      <div className="bg-gray-900/50 p-2 rounded border border-gray-700/50 mb-4 font-mono text-xs text-gray-400 break-all">
        {source.url}
      </div>

    
      <div className="space-y-2 text-sm text-gray-400 mb-6">
        <div className="flex items-center gap-2">
          <Clock className="w-4 h-4 text-gray-500" />
          <span>Periyot: <span className="text-gray-300">{source.scrape_interval}</span></span>
        </div>
        <div className="flex items-center gap-2">
          <ShieldCheck className="w-4 h-4 text-gray-500" />
          <span>Son Tarama: <span className="text-gray-300">
            {source.last_scraped_at ? formatRelativeTime(source.last_scraped_at) : 'Hiç taranmadı'}
          </span></span>
        </div>
      </div>


      <div className="flex items-center gap-2 pt-4 border-t border-gray-700">
        <Button 
          variant="secondary" 
          size="sm" 
          className="flex-1"
          onClick={() => onTrigger(source.id)}
          loading={isTriggering}
          icon={Play}
        >
          Tetikle
        </Button>
        <Button 
          variant="ghost" 
          size="sm" 
          className="px-2 text-blue-400 hover:text-blue-300"
          onClick={() => onEdit(source)}
        >
          <Edit2 className="w-4 h-4" />
        </Button>
        <Button 
          variant="ghost" 
          size="sm" 
          className="px-2 text-red-400 hover:text-red-300"
          onClick={() => onDelete(source.id)}
        >
          <Trash2 className="w-4 h-4" />
        </Button>
      </div>
    </div>
  );
};

export default SourceCard;