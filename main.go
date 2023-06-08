package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name     string `bson:"name"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

var client *mongo.Client

func connectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(context.Background(), clientOptions)
}

func main() {
	connectDB()
	defer client.Disconnect(context.Background())

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", showSignupForm)
	router.POST("/signup", handleSignup)

	router.Run(":8080")
}

func showSignupForm(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

func handleSignup(c *gin.Context) {
	// Retrieve the form data
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Create a new user document
	user := User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	// Insert the user into the database
	collection := client.Database("test").Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.String(http.StatusOK, "User registered successfully!")
}
