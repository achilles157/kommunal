package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"unleashed-space/models"
)

type ProfileHandler struct {
	db *mongo.Database
}

func NewProfileHandler(db *mongo.Database) *ProfileHandler {
	return &ProfileHandler{db: db}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	// Find user by ID
	var user models.User
	err := h.db.Collection("users").FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	var input models.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create update document
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	// Only update fields that were provided
	if input.Name != "" {
		update["$set"].(bson.M)["name"] = input.Name
	}
	if input.Username != "" {
		// Check if username is already taken
		var existingUser models.User
		err := h.db.Collection("users").FindOne(context.Background(), bson.M{
			"_id":      bson.M{"$ne": userID},
			"username": input.Username,
		}).Decode(&existingUser)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
		update["$set"].(bson.M)["username"] = input.Username
	}
	if input.Email != "" {
		// Check if email is already taken
		var existingUser models.User
		err := h.db.Collection("users").FindOne(context.Background(), bson.M{
			"_id":   bson.M{"$ne": userID},
			"email": input.Email,
		}).Decode(&existingUser)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}
		update["$set"].(bson.M)["email"] = input.Email
	}

	// Update user
	result := h.db.Collection("users").FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": userID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedUser models.User
	if err := result.Decode(&updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       updatedUser.ID,
			"name":     updatedUser.Name,
			"username": updatedUser.Username,
			"email":    updatedUser.Email,
		},
	})
}
