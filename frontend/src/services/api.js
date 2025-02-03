import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

// Create axios instance with base URL
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor
api.interceptors.request.use(
  (config) => {
    // Add token to request if it exists
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Auth services
export const signup = async (userData) => {
  try {
    const response = await api.post('/auth/signup', userData);
    if (response.data.token) {
      localStorage.setItem('token', response.data.token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
    }
    return response.data;
  } catch (error) {
    console.error('Signup error:', error);
    throw error;
  }
};

export const signin = async (credentials) => {
  try {
    const response = await api.post('/auth/signin', credentials);
    if (response.data.token) {
      localStorage.setItem('token', response.data.token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
    }
    return response.data;
  } catch (error) {
    console.error('Signin error:', error);
    throw error;
  }
};

export const signout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  // Clear any other stored data
  localStorage.clear();
};

// Profile services
export const getProfile = async () => {
  try {
    const response = await api.get('/profile');
    return response.data;
  } catch (error) {
    if (error.response?.status === 401) {
      // Token expired or invalid, sign out
      signout();
    }
    throw error;
  }
};

export const updateProfile = async (profileData) => {
  try {
    const response = await api.put('/profile', profileData);
    // Update stored user data if profile update successful
    if (response.data.user) {
      localStorage.setItem('user', JSON.stringify(response.data.user));
    }
    return response.data;
  } catch (error) {
    if (error.response?.status === 401) {
      signout();
    }
    throw error;
  }
};

// Post services
export const createPost = async (postData) => {
  try {
    const response = await api.post('/posts', postData);
    return response.data;
  } catch (error) {
    if (error.response?.status === 401) {
      signout();
    }
    throw error;
  }
};

export const getPosts = async () => {
  try {
    const response = await api.get('/posts');
    return response.data;
  } catch (error) {
    if (error.response?.status === 401) {
      signout();
    }
    throw error;
  }
};

export const getUserPosts = async () => {
  try {
    const response = await api.get('/posts/user');
    return response.data;
  } catch (error) {
    if (error.response?.status === 401) {
      signout();
    }
    throw error;
  }
};

export const getPublicFeed = async () => {
  try {
    const response = await api.get('/feed');
    return response.data;
  } catch (error) {
    console.error('Error fetching public feed:', error);
    throw error;
  }
};

export default api;
