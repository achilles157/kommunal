package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"unleashed-space/handlers"
	"unleashed-space/middleware"
)

func initMongoDB() (*mongo.Client, *mongo.Database, error) {
	// Load MongoDB URI from environment
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://127.0.0.1:27017"
		log.Printf("Using default MongoDB URI: %s", mongoURI)
	}

	log.Printf("Connecting to MongoDB at: %s", mongoURI)

	// Set client options
	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetServerSelectionTimeout(5 * time.Second).
		SetConnectTimeout(10 * time.Second)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, err
	}

	// Ping the database
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, err
	}

	// Get database name from URI or use default
	dbName := "unleashed_space"
	db := client.Database(dbName)

	// List existing collections
	collections, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		log.Printf("Error listing collections: %v", err)
	} else {
		log.Printf("Existing collections: %v", collections)
	}

	// Ensure collections exist
	err = ensureCollections(ctx, db)
	if err != nil {
		return nil, nil, err
	}

	return client, db, nil
}

func ensureCollections(ctx context.Context, db *mongo.Database) error {
	// Check if users collection exists
	collections, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return err
	}

	hasUsers := false
	for _, name := range collections {
		if name == "users" {
			hasUsers = true
			break
		}
	}

	// Create users collection if it doesn't exist
	if !hasUsers {
		log.Println("Creating users collection...")
		err = db.CreateCollection(ctx, "users")
		if err != nil {
			// Check if error is "collection already exists"
			cmdErr, ok := err.(mongo.CommandError)
			if !ok || cmdErr.Code != 48 {
				return err
			}
			log.Println("Users collection already exists")
		}
	}

	// Create indexes
	log.Println("Creating indexes for users collection...")
	usersCollection := db.Collection("users")

	// Drop existing indexes
	_, err = usersCollection.Indexes().DropAll(ctx)
	if err != nil {
		log.Printf("Warning: Failed to drop indexes: %v", err)
	}

	// Create new indexes
	_, err = usersCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return err
	}

	// Verify indexes
	indexes, err := usersCollection.Indexes().List(ctx)
	if err != nil {
		log.Printf("Warning: Failed to list indexes: %v", err)
	} else {
		var indexList []bson.M
		if err = indexes.All(ctx, &indexList); err != nil {
			log.Printf("Warning: Failed to decode indexes: %v", err)
		} else {
			log.Printf("Current indexes: %v", indexList)
		}
	}

	return nil
}

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Initialize MongoDB
	client, db, err := initMongoDB()
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	log.Println("Successfully connected to MongoDB")

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	profileHandler := handlers.NewProfileHandler(db)
	postHandler := handlers.NewPostHandler(db)

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Add error logging middleware
	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			log.Printf("Request errors: %v", c.Errors)
		}
	})

	// CORS Configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/signup", authHandler.SignUp)
			auth.POST("/signin", authHandler.SignIn)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Profile routes
			profile := protected.Group("/profile")
			{
				profile.GET("", profileHandler.GetProfile)
				profile.PUT("", profileHandler.UpdateProfile)
			}

			// Posts routes
			posts := protected.Group("/posts")
			{
				posts.POST("", postHandler.CreatePost)
				posts.GET("", postHandler.GetPosts)
				posts.GET("/user", postHandler.GetUserPosts)
			}
		}

		// Public routes
		api.GET("/feed", postHandler.GetPublicFeed)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
