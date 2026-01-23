import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom'; // <--- Link BURAYA EKLENDİ
import { Shield, Lock, User } from 'lucide-react';
import Button from '../../shared/components/Button';
import api from '../../shared/utils/api';

const LoginPage = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({ username: '', password: '' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response = await api.post('/auth/login', {
        username: formData.username,
        password: formData.password
      });

      if (response.data.success) {
        const { token, user } = response.data.data;
        
        localStorage.setItem('cti_auth_token', token);
        localStorage.setItem('cti_user', JSON.stringify(user));
        
        navigate('/');
      }
    } catch (err) {
      console.error("Login hatası:", err);
      setError(err.response?.data?.error?.message || 'Giriş başarısız oldu. Lütfen bilgileri kontrol edin.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center p-4">
      <div className="max-w-md w-full bg-gray-800 border border-gray-700 rounded-2xl shadow-2xl p-8">
        
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-blue-600/20 rounded-full mb-4">
            <Shield className="w-8 h-8 text-blue-500" />
          </div>
          <h2 className="text-2xl font-bold text-white">TorAnaliz</h2>
          <p className="text-gray-400 mt-2">Cyber Threat Intelligence Platform</p>
        </div>

        <form onSubmit={handleLogin} className="space-y-6">
          {error && (
            <div className="bg-red-900/20 border border-red-500/50 text-red-200 text-sm p-3 rounded-lg text-center animate-pulse">
              {error}
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">Kullanıcı Adı</label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <User className="h-5 w-5 text-gray-500" />
              </div>
              <input
                type="text"
                className="block w-full pl-10 bg-gray-900 border border-gray-600 rounded-lg py-2.5 text-white placeholder-gray-500 focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="admin"
                value={formData.username}
                onChange={(e) => setFormData({...formData, username: e.target.value})}
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">Şifre</label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <Lock className="h-5 w-5 text-gray-500" />
              </div>
              <input
                type="password"
                className="block w-full pl-10 bg-gray-900 border border-gray-600 rounded-lg py-2.5 text-white placeholder-gray-500 focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                placeholder="••••••••"
                value={formData.password}
                onChange={(e) => setFormData({...formData, password: e.target.value})}
              />
            </div>
          </div>

          <Button type="submit" variant="primary" className="w-full py-3 font-semibold text-lg" loading={loading}>
            Giriş Yap
          </Button>
        </form>

        <div className="mt-6 text-center text-sm">
          <span className="text-gray-400">Hesabın yok mu? </span>
          <Link to="/register" className="text-blue-400 hover:text-blue-300 font-medium transition-colors">
            Kayıt Ol
          </Link>
        </div>

        <div className="mt-8 text-center text-xs text-gray-500 border-t border-gray-700 pt-4">
          &copy; 2026 Interactive Scraper Project. Access Restricted.
        </div>
      </div>
    </div>
  );
};

export default LoginPage;