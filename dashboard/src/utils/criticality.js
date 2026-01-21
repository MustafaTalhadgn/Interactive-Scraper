export const getCriticalityStyle = (score) => {
  
  if (score >= 76) return { 
    label: 'CRITICAL', 
    classes: 'bg-red-100 text-red-700 border-red-200',
    dot: 'bg-red-600'
  };
  if (score >= 51) return { 
    label: 'HIGH', 
    classes: 'bg-orange-100 text-orange-700 border-orange-200',
    dot: 'bg-orange-600'
  };
  if (score >= 26) return { 
    label: 'MEDIUM', 
    classes: 'bg-yellow-100 text-yellow-700 border-yellow-200',
    dot: 'bg-yellow-600'
  };
  return { 
    label: 'LOW', 
    classes: 'bg-emerald-100 text-emerald-700 border-emerald-200',
    dot: 'bg-emerald-600'
  };
};