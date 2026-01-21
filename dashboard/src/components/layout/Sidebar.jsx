import React from 'react';
import { NavLink } from 'react-router-dom';
import { Shield, LayoutDashboard, Search, Database, Settings } from 'lucide-react';

const Sidebar = () => {
  return (
    <aside className="w-64 bg-slate-900 text-white flex flex-col shadow-xl">
      <div className="p-6 flex items-center gap-3 border-b border-slate-800">
        <Shield className="text-blue-400" size={28} />
        <span className="text-xl font-bold tracking-tight uppercase">TorAnaliz</span>
      </div>
      
      <nav className="flex-1 p-4 space-y-2">
        <MenuLink to="/dashboard" icon={<LayoutDashboard size={20}/>} label="Dashboard" />
        <MenuLink to="/intelligence" icon={<Search size={20}/>} label="Intelligence Feed" />
        <MenuLink to="/sources" icon={<Database size={20}/>} label="Sources" />
        <MenuLink to="/settings" icon={<Settings size={20}/>} label="Settings" />
      </nav>

      <div className="p-4 border-t border-slate-800 text-xs text-slate-500 text-center">
        TorAnaliz v1.0 | 2026
      </div>
    </aside>
  );
};

const MenuLink = ({ to, icon, label }) => (
  <NavLink 
    to={to}
    className={({ isActive }) => `
      flex items-center gap-3 p-3 rounded-lg transition-all duration-200
      ${isActive 
        ? 'bg-blue-600 text-white shadow-lg shadow-blue-900/20' 
        : 'text-slate-400 hover:bg-slate-800 hover:text-white'}
    `}
  >
    {icon}
    <span className="font-medium">{label}</span>
  </NavLink>
);

export default Sidebar;