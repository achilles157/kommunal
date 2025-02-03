package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"unleashed-space/models"
)

type AuthHandler struct {
	db *mongo.Database
}

func NewAuthHandler(db *mongo.Database) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse and validate input using Gin's binding
	var input models.SignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Processing signup request for email: %s", input.Email)

	// Get users collection
	collection := h.db.Collection("users")

	// Check if email exists
	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&existingUser)
	if err == nil {
		log.Printf("Email already exists: %s", input.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}
	if err != mongo.ErrNoDocuments {
		log.Printf("Error checking email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email availability"})
		return
	}

	// Check if username exists
	err = collection.FindOne(ctx, bson.M{"username": input.Username}).Decode(&existingUser)
	if err == nil {
		log.Printf("Username already exists: %s", input.Username)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}
	if err != mongo.ErrNoDocuments {
		log.Printf("Error checking username: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify username availability"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Create user
	now := time.Now()
	user := models.User{
		ID:        primitive.NewObjectID(),
		Name:      input.Name,
		Username:  input.Username,
		Email:     input.Email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Insert user
	log.Printf("Attempting to insert user with ID: %s", user.ID.Hex())
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		if mongo.IsDuplicateKeyError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	insertedID := result.InsertedID.(primitive.ObjectID)
	log.Printf("Successfully inserted user with ID: %s", insertedID.Hex())

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": insertedID.Hex(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	// Get JWT secret
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "unleashed_space_secret_key_2024_secure_random_string_for_jwt_signing"
	}

	// Sign token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		// Return success without token
		c.JSON(http.StatusCreated, gin.H{
			"user": gin.H{
				"id":       insertedID.Hex(),
				"name":     user.Name,
				"username": user.Username,
				"email":    user.Email,
			},
		})
		return
	}

	// Return success with token
	c.JSON(http.StatusCreated, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       insertedID.Hex(),
			"name":     user.Name,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var input models.SignInInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Processing signin request for email: %s", input.Email)

	// Find user by email
	var user models.User
	err := h.db.Collection("users").FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		log.Printf("Error finding user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process signin"})
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "unleashed_space_secret_key_2024_secure_random_string_for_jwt_signing"
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID.Hex(),
			"name":     user.Name,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
