import React, { useEffect } from 'react';
import { X } from 'lucide-react';


const Modal = ({ isOpen, onClose, title, children, size = 'md' }) => {
  useEffect(() => {
    const handleEsc = (e) => {
      if (e.key === 'Escape') onClose();
    };
    if (isOpen) window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  
  const sizeClasses = {
    sm: 'max-w-md',
    md: 'max-w-lg',      
    lg: 'max-w-2xl',
    xl: 'max-w-4xl',    
    '2xl': 'max-w-6xl',  
    full: 'max-w-[95vw]' 
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
   
      <div 
        className="fixed inset-0 bg-black/70 backdrop-blur-sm transition-opacity"
        onClick={onClose}
      />


      <div className={`relative w-full ${sizeClasses[size] || sizeClasses.md} bg-gray-800 border border-gray-700 rounded-xl shadow-2xl transform transition-all flex flex-col max-h-[90vh]`}>
        

        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-700 shrink-0">
          <h3 className="text-lg font-semibold text-white">{title}</h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-white transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>


        <div className="px-6 py-4 overflow-y-auto custom-scrollbar">
          {children}
        </div>
      </div>
    </div>
  );
};

export default Modal;