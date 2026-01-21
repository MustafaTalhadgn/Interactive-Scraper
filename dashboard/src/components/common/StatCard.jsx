import React from 'react';

const StatCard = ({ title, value, icon, colorClass }) => (
  <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 flex items-center justify-between">
    <div>
      <p className="text-sm text-gray-500 font-medium mb-1 uppercase tracking-wider">{title}</p>
      <h3 className="text-3xl font-bold">{value}</h3>
    </div>
    <div className={`p-4 rounded-lg ${colorClass}`}>
      {icon}
    </div>
  </div>
);

export default StatCard;