import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { signup } from '../services/api';

const Signup = ({ onAuthSuccess }) => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: '',
    username: '',
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

  const validateForm = () => {
    if (!formData.name.trim()) {
      setError('Name is required');
      return false;
    }
    if (!formData.username.trim()) {
      setError('Username is required');
      return false;
    }
    if (!formData.email.trim()) {
      setError('Email is required');
      return false;
    }
    if (!formData.password.trim()) {
      setError('Password is required');
      return false;
    }
    if (formData.password.length < 6) {
      setError('Password must be at least 6 characters long');
      return false;
    }
    return true;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    try {
      setIsLoading(true);
      setError('');
      
      const trimmedData = {
        name: formData.name.trim(),
        username: formData.username.trim(),
        email: formData.email.trim(),
        password: formData.password.trim(),
      };

      const response = await signup(trimmedData);
      
      if (response.token && response.user) {
        onAuthSuccess(response.user);
        navigate('/', { replace: true });
      } else {
        setError('Unexpected response from server');
      }
    } catch (error) {
      console.error('Signup error:', error);
      if (error.response?.data?.error) {
        setError(error.response.data.error);
      } else {
        setError('Failed to sign up. Please try again.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="auth-form-container">
      <form onSubmit={handleSubmit} className="auth-form">
        <h2>Sign Up</h2>
        {error && <div className="error-message">{error}</div>}
        <div className="form-group">
          <input 
            type="text" 
            name="name" 
            value={formData.name}
            placeholder="Name" 
            onChange={handleChange} 
            disabled={isLoading}
            required 
          />
        </div>
        <div className="form-group">
          <input 
            type="text" 
            name="username" 
            value={formData.username}
            placeholder="Username" 
            onChange={handleChange} 
            disabled={isLoading}
            required 
          />
        </div>
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
            placeholder="Password (min. 6 characters)" 
            onChange={handleChange} 
            disabled={isLoading}
            required 
            minLength={6}
          />
        </div>
        <button type="submit" disabled={isLoading}>
          {isLoading ? 'Signing up...' : 'Sign Up'}
        </button>
      </form>
    </div>
  );
};

export default Signup;
