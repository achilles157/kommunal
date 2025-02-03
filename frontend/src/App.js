import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, Navigate } from 'react-router-dom';
import Signup from './components/Signup.js';
import Signin from './components/Signin.js';
import Profile from './components/Profile.js';
import CreatePost from './components/CreatePost.js';
import PostList from './components/PostList.js';
import { signout } from './services/api.js';
import './App.css';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState(null);

  useEffect(() => {
    const token = localStorage.getItem('token');
    const storedUser = localStorage.getItem('user');
    if (token && storedUser) {
      setIsAuthenticated(true);
      const parsedUser = JSON.parse(storedUser);
      setUser(parsedUser);
      document.title = `Komunal - ${parsedUser.name}`;
    }
  }, []);

  const handleSignout = () => {
    signout();
    setIsAuthenticated(false);
    setUser(null);
    document.title = 'Komunal';
  };

  const handleAuthSuccess = (userData) => {
    setIsAuthenticated(true);
    setUser(userData);
    document.title = `Komunal - ${userData.name}`;
  };

  return (
    <Router>
      <div className="app">
        <nav className="navbar">
          <div className="nav-brand">
            <Link to="/">Komunal</Link>
          </div>
          <div className="nav-links">
            {isAuthenticated ? (
              <>
                <Link to="/profile">{user?.name || 'Profile'}</Link>
                <button onClick={handleSignout} className="signout-btn">
                  Sign Out
                </button>
              </>
            ) : (
              <>
                <Link to="/signin">Sign In</Link>
                <Link to="/signup">Sign Up</Link>
              </>
            )}
          </div>
        </nav>

        <main className="main-content">
          <Routes>
            <Route
              path="/"
              element={
                <div className="home-page">
                  {isAuthenticated && <CreatePost />}
                  <PostList />
                </div>
              }
            />
            <Route
              path="/signup"
              element={
                isAuthenticated ? (
                  <Navigate to="/" />
                ) : (
                  <div className="auth-page">
                    <Signup onAuthSuccess={handleAuthSuccess} />
                  </div>
                )
              }
            />
            <Route
              path="/signin"
              element={
                isAuthenticated ? (
                  <Navigate to="/" />
                ) : (
                  <div className="auth-page">
                    <Signin onAuthSuccess={handleAuthSuccess} />
                  </div>
                )
              }
            />
            <Route
              path="/profile"
              element={
                isAuthenticated ? (
                  <Profile user={user} />
                ) : (
                  <Navigate to="/signin" />
                )
              }
            />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;
