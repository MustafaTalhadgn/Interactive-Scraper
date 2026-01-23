import React, { useEffect, useState } from 'react';
import { Plus, RefreshCw } from 'lucide-react';
import useApi from '../../shared/hooks/useApi';
import { sourceService } from './services/sourceService';
import SourceCard from './components/SourceCard';
import SourceModal from './components/SourceModal';
import Button from '../../shared/components/Button';
import LoadingSpinner from '../../shared/components/LoadingSpinner';

const SourcesPage = () => {
  const { data, loading, execute: fetchSources } = useApi(sourceService.getAll);
  

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingSource, setEditingSource] = useState(null);
  

  const [triggeringId, setTriggeringId] = useState(null);

  useEffect(() => {
    fetchSources();
  }, []);

  
  const handleSave = async (formData) => {
    try {
      if (editingSource) {
        await sourceService.update(editingSource.id, formData);
      } else {
        await sourceService.create(formData);
      }
      setIsModalOpen(false);
      setEditingSource(null);
      fetchSources(); 
    } catch (error) {
      console.error("Save failed:", error);
      alert("İşlem başarısız oldu: " + error.message);
    }
  };


  const handleDelete = async (id) => {
    if (window.confirm('Bu kaynağı silmek istediğinize emin misiniz?')) {
      try {
        await sourceService.delete(id);
        fetchSources();
      } catch (error) {
        console.error("Delete failed:", error);
      }
    }
  };


  const handleTrigger = async (id) => {
    try {
      setTriggeringId(id); 
      await sourceService.triggerScrape(id);
      
   
      alert('Tarama görevi tetiklendi! Scraper arka planda çalışmaya başladı.');
      
    
      setTimeout(() => setTriggeringId(null), 2000);
    } catch (error) {
      console.error("Trigger failed:", error);
      setTriggeringId(null);
      alert('Tetikleme başarısız: ' + error.message);
    }
  };

  const openEditModal = (source) => {
    setEditingSource(source);
    setIsModalOpen(true);
  };

  const openCreateModal = () => {
    setEditingSource(null);
    setIsModalOpen(true);
  };

  return (
    <div className="space-y-6">
 
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white">Kaynak Yönetimi</h1>
          <p className="text-gray-400 mt-1">Dark Web veri kaynaklarını izleyin ve yönetin.</p>
        </div>
        <Button onClick={openCreateModal} icon={Plus}>Yeni Kaynak Ekle</Button>
      </div>


      {loading && !data ? (
        <div className="flex justify-center py-12">
          <LoadingSpinner size="large" />
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {data?.data?.sources?.map((source) => (
            <SourceCard 
              key={source.id} 
              source={source} 
              onEdit={openEditModal}
              onDelete={handleDelete}
              onTrigger={handleTrigger}
              isTriggering={triggeringId === source.id}
            />
          ))}
          
       
          {(!data?.data?.sources || data.data.sources.length === 0) && (
            <div className="col-span-full bg-gray-800/50 border border-gray-700 rounded-xl p-12 text-center border-dashed">
              <p className="text-gray-400 mb-4">Henüz hiç kaynak eklenmemiş.</p>
              <Button variant="secondary" onClick={openCreateModal}>İlk Kaynağı Ekle</Button>
            </div>
          )}
        </div>
      )}

      <SourceModal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSave}
        initialData={editingSource}
      />
    </div>
  );
};

export default SourcesPage;