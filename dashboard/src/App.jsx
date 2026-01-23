import React from 'react';
import { BrowserRouter, Routes, Route, Navigate, Outlet } from 'react-router-dom';
import Layout from './shared/components/Layout';
import OverviewPage from './features/overview/OverviewPage';
import IntelligencePage from './features/intelligence/IntelligencePage';
import SourcesPage from './features/sources/SourcesPage';
import LoginPage from './features/auth/LoginPage'; 
import RegisterPage from './features/auth/RegisterPage';


const ProtectedRoute = () => {
  const token = localStorage.getItem('cti_auth_token');
  

  if (!token) {
    return <Navigate to="/login" replace />;
  }


  return <Outlet />;
};

function App() {
  return (
  
      <Routes>
    
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route element={<ProtectedRoute />}>
          <Route path="/" element={<Layout />}>
            <Route index element={<OverviewPage />} />
            <Route path="intelligence" element={<IntelligencePage />} />
            <Route path="sources" element={<SourcesPage />} />
          </Route>
        </Route>

     
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
  );
}

export default App;