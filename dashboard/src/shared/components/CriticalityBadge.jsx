import React from 'react';
import { CRITICALITY_BG_CLASSES } from '../utils/constants';

const CriticalityBadge = ({ level, score }) => {
  const safeLevel = level?.toLowerCase() || 'low';
  const classes = CRITICALITY_BG_CLASSES[safeLevel] || CRITICALITY_BG_CLASSES.low;

  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${classes}`}>
      <span className="w-2 h-2 mr-1.5 rounded-full bg-current opacity-75"></span>
      {safeLevel.toUpperCase()}
      {score !== undefined && <span className="ml-1 opacity-75">({score})</span>}
    </span>
  );
};

export default CriticalityBadge;