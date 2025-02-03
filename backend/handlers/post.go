package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"unleashed-space/models"
)

type PostHandler struct {
	db *mongo.Database
}

func NewPostHandler(db *mongo.Database) *PostHandler {
	return &PostHandler{db: db}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var input models.CreatePostInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	// Get user details for post author
	var user models.User
	err := h.db.Collection("users").FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user details"})
		return
	}

	// Create post
	now := time.Now()
	post := models.Post{
		UserID:  userID.(primitive.ObjectID),
		Content: input.Content,
		Author: models.PostAuthor{
			Name:     user.Name,
			Username: user.Username,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := h.db.Collection("posts").InsertOne(context.Background(), post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	post.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	// Set up options for pagination and sorting
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(20)

	cursor, err := h.db.Collection("posts").Find(context.Background(), bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	defer cursor.Close(context.Background())

	var posts []models.Post
	if err := cursor.All(context.Background(), &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func (h *PostHandler) GetUserPosts(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	// Set up options for pagination and sorting
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(20)

	cursor, err := h.db.Collection("posts").Find(context.Background(), bson.M{"user_id": userID}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user posts"})
		return
	}
	defer cursor.Close(context.Background())

	var posts []models.Post
	if err := cursor.All(context.Background(), &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func (h *PostHandler) GetPublicFeed(c *gin.Context) {
	// Set up options for pagination and sorting
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(20)

	cursor, err := h.db.Collection("posts").Find(context.Background(), bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	defer cursor.Close(context.Background())

	var posts []models.Post
	if err := cursor.All(context.Background(), &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}
