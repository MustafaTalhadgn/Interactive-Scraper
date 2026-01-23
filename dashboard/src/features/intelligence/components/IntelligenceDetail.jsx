import React, { useEffect } from 'react';
import { 
  Calendar, Globe, Shield, ExternalLink, 
  Hash, Mail, Server, Bitcoin, AlertTriangle, Key, Copy, Check
} from 'lucide-react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import Modal from '../../../shared/components/Modal';
import CriticalityBadge from '../../../shared/components/CriticalityBadge';
import LoadingSpinner from '../../../shared/components/LoadingSpinner';
import Button from '../../../shared/components/Button';
import useApi from '../../../shared/hooks/useApi';
import { intelligenceService } from '../services/intelligenceService';
import { formatDate, copyToClipboard } from '../../../shared/utils/helpers';
import { FEATURE_DISPLAY_NAMES } from '../../../shared/utils/constants';


const FeatureSection = ({ title, items, icon: Icon, colorClass }) => {
  const [copiedIndex, setCopiedIndex] = React.useState(null);

  if (!items || items.length === 0) return null;

  const handleCopy = async (text, index) => {
    await copyToClipboard(text);
    setCopiedIndex(index);
    setTimeout(() => setCopiedIndex(null), 2000);
  };

  return (
    <div className="mb-4 bg-gray-900/30 p-3 rounded-lg border border-gray-700/50">
      <h4 className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3 flex items-center gap-2">
        <Icon className="w-3.5 h-3.5" /> {title}
      </h4>
      <div className="flex flex-wrap gap-2">
        {items.map((item, index) => (
          <button
            key={index}
            onClick={() => handleCopy(item, index)}
            className={`group flex items-center gap-2 text-xs px-2.5 py-1.5 rounded border transition-all hover:brightness-110 ${colorClass}`}
            title="Kopyalamak için tıkla"
          >
            <span className="font-mono">{item}</span>
            {copiedIndex === index ? (
              <Check className="w-3 h-3" />
            ) : (
              <Copy className="w-3 h-3 opacity-0 group-hover:opacity-100 transition-opacity" />
            )}
          </button>
        ))}
      </div>
    </div>
  );
};

const IntelligenceDetail = ({ id, isOpen, onClose }) => {
  const { data, loading, execute } = useApi(intelligenceService.getDetail);

  useEffect(() => {
    if (isOpen && id) {
      execute(id);
    }
  }, [isOpen, id, execute]);

  if (loading) {
    return (
      <Modal isOpen={isOpen} onClose={onClose} title="Yükleniyor...">
        <div className="flex flex-col items-center justify-center py-12">
          <LoadingSpinner size="large" />
          <p className="mt-4 text-gray-400 animate-pulse">Tehdit verisi analiz ediliyor...</p>
        </div>
      </Modal>
    );
  }

  const item = data?.data;
  if (!item) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="2xl" title="Tehdit İstihbarat Detayı">
      <div className="flex flex-col h-[80vh] "> 

        <div className="flex-1 overflow-y-auto pr-2 custom-scrollbar space-y-6">
          
 
          <div className="space-y-4">
            <div className="flex items-start justify-between gap-4">
              <h2 className="text-xl font-bold text-white leading-tight">
                {item.title || 'Başlıksız Tehdit'}
              </h2>
              <CriticalityBadge level={item.criticality} score={item.criticality_score} />
            </div>

            <div className="flex flex-wrap gap-y-2 gap-x-6 text-sm text-gray-400 border-b border-gray-700 pb-4">
              <div className="flex items-center gap-2">
                <Globe className="w-4 h-4 text-blue-400" />
                <span className="text-gray-200 font-medium">{item.source_name}</span>
                <span className="text-xs bg-gray-700 px-2 py-0.5 rounded border border-gray-600">
                  {item.category?.toUpperCase()}
                </span>
              </div>
              <div className="flex items-center gap-2">
                <Calendar className="w-4 h-4 text-gray-500" />
                <span>Tespit: {formatDate(item.created_at)}</span>
              </div>
            </div>
          </div>

      
          {item.extracted_features && (
            <div className="space-y-1">
              <FeatureSection 
                title={FEATURE_DISPLAY_NAMES.bitcoin_addresses} 
                items={item.extracted_features.bitcoin_addresses} 
                icon={Bitcoin}
                colorClass="bg-orange-900/20 border-orange-700/50 text-orange-400 hover:bg-orange-900/40"
              />
              <FeatureSection 
                title={FEATURE_DISPLAY_NAMES.onion_urls} 
                items={item.extracted_features.onion_urls} 
                icon={Globe}
                colorClass="bg-green-900/20 border-green-700/50 text-green-400 hover:bg-green-900/40"
              />
              <FeatureSection 
                title={FEATURE_DISPLAY_NAMES.ip_addresses} 
                items={item.extracted_features.ip_addresses} 
                icon={Server}
                colorClass="bg-blue-900/20 border-blue-700/50 text-blue-400 hover:bg-blue-900/40"
              />
              <FeatureSection 
                title={FEATURE_DISPLAY_NAMES.emails} 
                items={item.extracted_features.emails} 
                icon={Mail}
                colorClass="bg-yellow-900/20 border-yellow-700/50 text-yellow-400 hover:bg-yellow-900/40"
              />
            </div>
          )}

        
          <div>
            <h3 className="text-lg font-semibold text-white mb-3 flex items-center gap-2">
              <Shield className="w-5 h-5 text-blue-500" />
              İçerik Analizi
            </h3>
            <div className="bg-gray-900 rounded-lg border border-gray-700 p-5 relative group">

              <div className="text-gray-300 text-sm ...">
              <ReactMarkdown 
                  remarkPlugins={[remarkGfm]} 
                  components={{
                    
                      a: ({node, ...props}) => <a {...props} className="text-blue-400 hover:underline" target="_blank" rel="noopener noreferrer" />
                  }}
              >
                  {item.summary || 'İçerik özeti oluşturulamadı.'}
              </ReactMarkdown>
          </div>
               

               <button 
                 onClick={() => copyToClipboard(item.summary)}
                 className="absolute top-4 right-4 p-2 bg-gray-800 rounded-md text-gray-400 opacity-0 group-hover:opacity-100 transition-all hover:text-white hover:bg-gray-700 border border-gray-600"
                 title="Metni Kopyala"
               >
                 <Copy className="w-4 h-4" />
               </button>
            </div>
          </div>

        </div>


        <div className="pt-4 mt-4 border-t border-gray-700 flex justify-end gap-3">
          <a 
            href={item.source_url} 
            target="_blank" 
            rel="noopener noreferrer"
            className="flex-1 sm:flex-none"
          >
            <Button variant="secondary" className="w-full" icon={ExternalLink}>
              Kaynağa Git (Tor)
            </Button>
          </a>
          <Button variant="primary" onClick={onClose}>
            Kapat
          </Button>
        </div>
      </div>
    </Modal>
  );
};

export default IntelligenceDetail;