import React, { useState, useEffect } from 'react';
import { Outlet, NavLink, useNavigate } from 'react-router-dom'; 
import {
  LayoutDashboard,
  Database,
  ShieldAlert,
  Menu,
  X,
  LogOut,
  Globe,
} from 'lucide-react';
import { clsx } from 'clsx';

const navigation = [
  { name: 'Overview', href: '/', icon: LayoutDashboard },
  { name: 'Intelligence', href: '/intelligence', icon: ShieldAlert },
  { name: 'Sources', href: '/sources', icon: Database },
];

const Layout = () => {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const navigate = useNavigate();
  

  const [user, setUser] = useState({ username: 'Misafir', role: 'Ziyaretçi' });


  useEffect(() => {
    const storedUser = localStorage.getItem('cti_user');
    if (storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser);
        setUser({
      
          username: parsedUser.username || parsedUser.name || 'Kullanıcı',
          role: parsedUser.role || 'Analyst'
        });
      } catch (error) {
        console.error("Kullanıcı verisi okunamadı:", error);
      }
    }
  }, []);

  const handleLogout = () => {
  
    localStorage.removeItem('cti_auth_token');
    localStorage.removeItem('cti_user');
    
  
    navigate('/login');
   
  };

  return (
    <div className="h-screen flex overflow-hidden bg-gray-900 text-gray-100 font-sans">

      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-gray-900/80 backdrop-blur-sm lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}


      <aside
        className={clsx(
          'z-50 w-72 bg-gray-800 border-r border-gray-700 flex-shrink-0',
          'fixed inset-y-0 left-0 transition-transform duration-300',
          'lg:static lg:translate-x-0',
          'lg:static lg:translate-x-0',
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        )}
      >
        <div className="flex flex-col h-full">

          <div className="flex items-center justify-between h-16 px-6 border-b border-gray-700">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-blue-600 rounded-lg">
                <Globe className="w-6 h-6 text-white" />
              </div>
              <span className="text-xl font-bold">TorAnaliz</span>
            </div>
            <button
              className="lg:hidden text-gray-400 hover:text-white"
              onClick={() => setSidebarOpen(false)}
            >
              <X className="w-6 h-6" />
            </button>
          </div>


          <nav className="flex-1 px-4 py-6 space-y-2 overflow-y-auto">
            {navigation.map((item) => (
              <NavLink
                key={item.name}
                to={item.href}
                onClick={() => setSidebarOpen(false)}
                className={({ isActive }) =>
                  clsx(
                    'flex items-center gap-3 px-4 py-3 text-sm font-medium rounded-lg transition-colors',
                    isActive
                      ? 'bg-blue-600/10 text-blue-400 border border-blue-600/20'
                      : 'text-gray-400 hover:bg-gray-700/50 hover:text-white'
                  )
                }
              >
                <item.icon className="w-5 h-5" />
                {item.name}
              </NavLink>
            ))}
          </nav>

          <div className="p-4 border-t border-gray-700">
            <div className="flex items-center gap-3">
  
              <div className="w-10 h-10 rounded-full bg-gradient-to-tr from-blue-500 to-purple-500 flex items-center justify-center font-bold text-lg uppercase shadow-lg shadow-blue-900/20">
                {user.username.charAt(0)}
              </div>
              
              <div className="flex-1 min-w-0">
                <p className="text-sm font-semibold text-white truncate capitalize">
                  {user.username}
                </p>
                <p className="text-xs text-gray-400 truncate capitalize flex items-center gap-1">
                  <span className="w-1.5 h-1.5 rounded-full bg-green-500 inline-block"></span>
                  {user.role}
                </p>
              </div>
              
              <button
                onClick={handleLogout}
                className="p-2 text-gray-400 hover:text-red-400 hover:bg-red-400/10 rounded-lg transition-all"
                title="Çıkış Yap"
              >
                <LogOut className="w-5 h-5" />
              </button>
            </div>
          </div>
        </div>
      </aside>

      <div className="flex-1 flex flex-col min-w-0">
  
        <header className="h-16 flex items-center gap-4 px-4 border-b border-gray-800 lg:hidden">
          <button
            className="p-2 text-gray-400 hover:text-white"
            onClick={() => setSidebarOpen(true)}
          >
            <Menu className="w-6 h-6" />
          </button>
          <span className="text-lg font-bold">TorAnaliz</span>
        </header>

      
        <main className="flex-1 overflow-y-auto p-4 sm:p-6 lg:px-8 lg:py-6 custom-scrollbar">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default Layout;