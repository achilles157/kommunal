import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { signin } from '../services/api';

const Signin = ({ onAuthSuccess }) => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
    if (error) setError('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    try {
      setIsLoading(true);
      setError('');
      
      const response = await signin(formData);
      
      if (response.token && response.user) {
        onAuthSuccess(response.user);
        navigate('/', { replace: true });
      } else {
        setError('Unexpected response from server');
      }
    } catch (error) {
      console.error('Signin error:', error);
      if (error.response?.data?.error) {
        setError(error.response.data.error);
      } else {
        setError('Failed to sign in. Please try again.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="auth-form-container">
      <form onSubmit={handleSubmit} className="auth-form">
        <h2>Sign In</h2>
        {error && <div className="error-message">{error}</div>}
        <div className="form-group">
          <input 
            type="email" 
            name="email" 
            value={formData.email}
            placeholder="Email" 
            onChange={handleChange} 
            disabled={isLoading}
            required 
          />
        </div>
        <div className="form-group">
          <input 
            type="password" 
            name="password" 
            value={formData.password}
            placeholder="Password" 
            onChange={handleChange} 
            disabled={isLoading}
            required 
          />
        </div>
        <button type="submit" disabled={isLoading}>
          {isLoading ? 'Signing in...' : 'Sign In'}
        </button>
      </form>
    </div>
  );
};

export default Signin;
