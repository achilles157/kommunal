import React, { useState, useEffect } from 'react';
import { getProfile, updateProfile, getUserPosts } from '../services/api';

const Profile = () => {
  const [profile, setProfile] = useState(null);
  const [userPosts, setUserPosts] = useState([]);
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    username: '',
    email: '',
  });
  const [error, setError] = useState('');

  useEffect(() => {
    fetchProfile();
    fetchUserPosts();
  }, []);

  const fetchProfile = async () => {
    try {
      const response = await getProfile();
      setProfile(response.user);
      setFormData({
        name: response.user.name,
        username: response.user.username,
        email: response.user.email,
      });
    } catch (error) {
      setError('Failed to fetch profile');
      console.error('Failed to fetch profile:', error);
    }
  };

  const fetchUserPosts = async () => {
    try {
      const response = await getUserPosts();
      setUserPosts(response.posts || []);
    } catch (error) {
      console.error('Failed to fetch user posts:', error);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await updateProfile(formData);
      setProfile(response.user);
      setIsEditing(false);
      setError('');
    } catch (error) {
      setError(error.response?.data?.error || 'Failed to update profile');
    }
  };

  if (!profile && !error) {
    return <div className="loading">Loading profile...</div>;
  }

  return (
    <div className="profile-container">
      <div className="profile-section">
        <h2>Profile</h2>
        {error && <div className="error-message">{error}</div>}
        {isEditing ? (
          <form onSubmit={handleSubmit} className="profile-edit-form">
            <input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Name"
              className="profile-input"
            />
            <input
              type="text"
              name="username"
              value={formData.username}
              onChange={handleChange}
              placeholder="Username"
              className="profile-input"
            />
            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              placeholder="Email"
              className="profile-input"
            />
            <div className="button-group">
              <button type="submit" className="save-button">Save</button>
              <button type="button" onClick={() => setIsEditing(false)} className="cancel-button">
                Cancel
              </button>
            </div>
          </form>
        ) : (
          <div className="profile-info">
            <p><strong>Name:</strong> {profile?.name}</p>
            <p><strong>Username:</strong> {profile?.username}</p>
            <p><strong>Email:</strong> {profile?.email}</p>
            <button onClick={() => setIsEditing(true)} className="edit-button">Edit Profile</button>
          </div>
        )}
      </div>

      <div className="user-posts-section">
        <h3>Your Posts</h3>
        {userPosts.length === 0 ? (
          <p>No posts yet</p>
        ) : (
          <div className="posts-grid">
            {userPosts.map((post) => (
              <div key={post.id} className="post-card">
                <p>{post.content}</p>
                <small>Posted on: {new Date(post.created_at).toLocaleDateString()}</small>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default Profile;
