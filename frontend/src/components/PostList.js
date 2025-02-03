import React, { useState, useEffect } from 'react';
import { getPublicFeed } from '../services/api';

function PostList() {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchPosts = async () => {
    try {
      setLoading(true);
      const response = await getPublicFeed();
      setPosts(response?.posts || []);
      setError(null);
    } catch (err) {
      setError('Failed to load posts. Please try again later.');
      console.error('Error fetching posts:', err);
      setPosts([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, []);

  if (loading) {
    return <div className="loading">Loading posts...</div>;
  }

  if (error) {
    return <div className="error">{error}</div>;
  }

  return (
    <div className="post-list">
      <h2>Recent Posts</h2>
      {posts && posts.length > 0 ? (
        <div className="posts-container">
          {posts.map((post) => (
            <div key={post.id} className="post-card">
              <div className="post-header">
                <div className="post-author">
                  <span className="author-name">{post.author?.name || 'Unknown'}</span>
                  <span className="author-username">@{post.author?.username || 'anonymous'}</span>
                </div>
                <span className="post-date">
                  {new Date(post.created_at).toLocaleDateString()}
                </span>
              </div>
              <div className="post-content">{post.content}</div>
            </div>
          ))}
        </div>
      ) : (
        <p className="no-posts">No posts yet</p>
      )}
    </div>
  );
}

export default PostList;
