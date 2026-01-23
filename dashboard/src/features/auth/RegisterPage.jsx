import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Shield, Lock, User, UserPlus } from 'lucide-react';
import Button from '../../shared/components/Button';
import api from '../../shared/utils/api';

const RegisterPage = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({ username: '', password: '' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleRegister = async (e) => {
    e.preventDefault();
    setError('');
    
    // Basit validasyon
    if (formData.password.length < 6) {
        setError('Şifre en az 6 karakter olmalıdır.');
        return;
    }

    setLoading(true);

    try {
      const response = await api.post('/auth/register', formData);

      if (response.data.success) {
        alert("Kayıt başarılı! Giriş sayfasına yönlendiriliyorsunuz.");
        navigate('/login');
      }
    } catch (err) {
      console.error("Kayıt hatası:", err);
      setError(err.response?.data?.error?.message || 'Kayıt işlemi başarısız oldu.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center p-4">
      <div className="max-w-md w-full bg-gray-800 border border-gray-700 rounded-2xl shadow-2xl p-8">
        
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-green-600/20 rounded-full mb-4">
            <UserPlus className="w-8 h-8 text-green-500" />
          </div>
          <h2 className="text-2xl font-bold text-white">Hesap Oluştur</h2>
          <p className="text-gray-400 mt-2">TorAnaliz Platformuna Katılın</p>
        </div>

        <form onSubmit={handleRegister} className="space-y-6">
          {error && (
            <div className="bg-red-900/20 border border-red-500/50 text-red-200 text-sm p-3 rounded-lg text-center">
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
                required
                className="block w-full pl-10 bg-gray-900 border border-gray-600 rounded-lg py-2.5 text-white placeholder-gray-500 focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all"
                placeholder="analyst_mustafa"
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
                required
                className="block w-full pl-10 bg-gray-900 border border-gray-600 rounded-lg py-2.5 text-white placeholder-gray-500 focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all"
                placeholder="••••••••"
                value={formData.password}
                onChange={(e) => setFormData({...formData, password: e.target.value})}
              />
            </div>
          </div>

          <Button type="submit" variant="primary" className="w-full py-3 bg-green-600 hover:bg-green-700 focus:ring-green-500" loading={loading}>
            Kayıt Ol
          </Button>
        </form>

        <div className="mt-6 text-center text-sm">
          <span className="text-gray-400">Zaten hesabın var mı? </span>
          <Link to="/login" className="text-blue-400 hover:text-blue-300 font-medium transition-colors">
            Giriş Yap
          </Link>
        </div>
      </div>
    </div>
  );
};

export default RegisterPage;