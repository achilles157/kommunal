import React, { useState } from 'react';
import { createPost } from '../services/api';

const CreatePost = ({ onPostCreated }) => {
  const [content, setContent] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!content.trim()) return;

    setIsSubmitting(true);
    setError('');

    try {
      await createPost({ content: content.trim() });
      setContent('');
      if (onPostCreated) {
        onPostCreated();
      }
      // Reload the page to show the new post
      window.location.reload();
    } catch (error) {
      setError(error.response?.data?.error || 'Failed to create post');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="create-post">
      <h3>Create a Post</h3>
      {error && <div className="error-message">{error}</div>}
      <form onSubmit={handleSubmit} className="create-post-form">
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          placeholder="What's on your mind?"
          rows="4"
          required
          className="create-post-textarea"
        />
        <button type="submit" disabled={isSubmitting || !content.trim()} className="create-post-button">
          {isSubmitting ? 'Posting...' : 'Post'}
        </button>
      </form>
    </div>
  );
};

export default CreatePost;
