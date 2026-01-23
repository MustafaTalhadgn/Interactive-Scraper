import React, { useState, useEffect } from 'react';
import Modal from '../../../shared/components/Modal';
import Button from '../../../shared/components/Button';
import { SOURCE_CATEGORIES, SCRAPE_INTERVALS, CRITICALITY_LEVELS } from '../../../shared/utils/constants';

const SourceModal = ({ isOpen, onClose, onSubmit, initialData, loading }) => {
  const [formData, setFormData] = useState({
    name: '',
    url: '',
    category: 'forum',
    criticality: 'medium',
    scrape_interval: '1 hour',
    enabled: true
  });

  useEffect(() => {
    if (initialData) {
      setFormData(initialData);
    } else {
      setFormData({
        name: '',
        url: '',
        category: 'forum',
        criticality: 'medium',
        scrape_interval: '1 hour',
        enabled: true
      });
    }
  }, [initialData, isOpen]);

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <Modal 
      isOpen={isOpen} 
      onClose={onClose} 
      title={initialData ? 'Kaynağı Düzenle' : 'Yeni Kaynak Ekle'}
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1">Kaynak Adı</label>
          <input
            type="text"
            required
            className="w-full bg-gray-900 border border-gray-600 rounded-lg px-3 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:outline-none"
            value={formData.name}
            onChange={e => setFormData({...formData, name: e.target.value})}
            placeholder="Örn: Tor News Forum"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1">Onion URL</label>
          <input
            type="text"
            required
            className="w-full bg-gray-900 border border-gray-600 rounded-lg px-3 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:outline-none font-mono text-sm"
            value={formData.url}
            onChange={e => setFormData({...formData, url: e.target.value})}
            placeholder="http://example.onion"
            disabled={!!initialData} 
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">Kategori</label>
            <select
              className="w-full bg-gray-900 border border-gray-600 rounded-lg px-3 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:outline-none"
              value={formData.category}
              onChange={e => setFormData({...formData, category: e.target.value})}
            >
              {SOURCE_CATEGORIES.map(cat => (
                <option key={cat} value={cat}>{cat.toUpperCase()}</option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">Kritiklik (Önem Derecesi)</label>
            <select
              className="w-full bg-gray-900 border border-gray-600 rounded-lg px-3 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:outline-none"
              value={formData.criticality}
              onChange={e => setFormData({...formData, criticality: e.target.value})}
            >
              {Object.values(CRITICALITY_LEVELS).map(level => (
                <option key={level} value={level}>{level.toUpperCase()}</option>
              ))}
            </select>
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1">Tarama Sıklığı</label>
          <select
            className="w-full bg-gray-900 border border-gray-600 rounded-lg px-3 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:outline-none"
            value={formData.scrape_interval}
            onChange={e => setFormData({...formData, scrape_interval: e.target.value})}
          >
            {SCRAPE_INTERVALS.map(interval => (
              <option key={interval.value} value={interval.value}>{interval.label}</option>
            ))}
          </select>
        </div>

        <div className="flex justify-end gap-3 pt-4 border-t border-gray-700 mt-6">
          <Button variant="secondary" onClick={onClose} type="button">İptal</Button>
          <Button variant="primary" type="submit" loading={loading}>
            {initialData ? 'Güncelle' : 'Kaydet'}
          </Button>
        </div>
      </form>
    </Modal>
  );
};

export default SourceModal;