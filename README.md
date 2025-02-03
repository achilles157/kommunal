# Komunal

## Description
Komunal is a web application that allows users to create and share posts. It features user authentication, profile management, and a feed of posts. The application is built with a Go backend and a React frontend.

## Installation

### Backend
1. Clone the repository.
2. Navigate to the `backend` directory.
3. Create a `.env` file and set the following environment variables:
   ```
   MONGO_URI=<your_mongo_uri>
   PORT=8080
   GIN_MODE=debug
   ```
4. Install dependencies:
   ```bash
   go mod tidy
   ```
5. Run the server:
   ```bash
   go run main.go
   ```

### Frontend
1. Navigate to the `frontend` directory.
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the development server:
   ```bash
   npm start
   ```

## Usage
- Sign up for a new account or sign in to an existing account.
- Create posts and view them in the feed.
- Manage your profile and view your posts.

## API Endpoints
- **POST /api/auth/signup**: Create a new user account.
- **POST /api/auth/signin**: Sign in to an existing account.
- **GET /api/profile**: Get the authenticated user's profile.
- **PUT /api/profile**: Update the authenticated user's profile.
- **POST /api/posts**: Create a new post.
- **GET /api/posts**: Get all posts.
- **GET /api/posts/user**: Get posts by the authenticated user.
- **GET /api/feed**: Get the public feed of posts.

## Frontend Components
- **Signup**: Component for user registration.
- **Signin**: Component for user login.
- **Profile**: Component for displaying and editing user profiles.
- **CreatePost**: Component for creating new posts.
- **PostList**: Component for displaying a list of posts.

## Technologies Used
- **Backend**: Go, Gin, MongoDB
- **Frontend**: React, React Router

## License
This project is licensed under the MIT License.
