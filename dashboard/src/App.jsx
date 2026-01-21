import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import MainLayout from './components/layout/MainLayout';
import Dashboard from './pages/Dashboard';
import IntelligenceFeed from './pages/IntelligenceFeed';
function App() {
  return (
 
      <MainLayout>
      <Routes>
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/intelligence" element={<IntelligenceFeed />} />
        
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
      </MainLayout>
 
  );
}

export default App;