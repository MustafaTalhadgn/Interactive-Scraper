import React from 'react';
import { PieChart, Pie, Cell, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { capitalizeFirst } from '../../../shared/utils/helpers';


const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899', '#6366f1'];

const CategoryPieChart = ({ data }) => {

  if (!data || Object.keys(data).length === 0) {
    return (
      <div className="h-full flex flex-col items-center justify-center text-gray-500 bg-gray-800 rounded-xl border border-gray-700 border-dashed p-6">
        <span>Veri bulunamadı</span>
      </div>
    );
  }


  const chartData = Object.entries(data).map(([name, value]) => ({
    name: capitalizeFirst(name),
    value: value
  })).sort((a, b) => b.value - a.value); 

  return (
    <div className="bg-gray-800 rounded-xl border border-gray-700 p-6 h-full flex flex-col">
      <h3 className="text-lg font-semibold text-white mb-4">Kategori Dağılımı</h3>
      
      <div className="flex-1 min-h-[250px]">
        <ResponsiveContainer width="100%" height="100%">
          <PieChart>
            <Pie
              data={chartData}
              cx="50%"
              cy="50%"
              innerRadius={60} 
              outerRadius={80}
              paddingAngle={5}
              dataKey="value"
            >
              {chartData.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} stroke="rgba(0,0,0,0.2)" />
              ))}
            </Pie>
            <Tooltip 
              contentStyle={{ backgroundColor: '#1f2937', borderColor: '#374151', color: '#fff', borderRadius: '8px' }}
              itemStyle={{ color: '#fff' }}
              formatter={(value) => [`${value} Kaynak`, 'Adet']}
            />
            <Legend 
              verticalAlign="bottom" 
              height={36}
              iconType="circle"
              formatter={(value) => <span className="text-gray-300 ml-1">{value}</span>}
            />
          </PieChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};

export default CategoryPieChart;