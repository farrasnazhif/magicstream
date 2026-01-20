package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/farrasnazhif/moviestream-go/database"
	"github.com/farrasnazhif/moviestream-go/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

// global mongo collection reference for "user"
var userCollection *mongo.Collection = database.OpenCollection("users")

func HashPassword(password string) (string, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(HashPassword), nil

}

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// create context with timeout for insert operation
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// user struct that will receive the json request body
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			// return 400 if json is invalid or missing required json fields
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// validator instanse like zod
		var validate = validator.New()

		// validate struct based on validate tags in models.User
		if err := validate.Struct(user); err != nil {
			// return 400 if validation fails
			// gin.H = a shortcut to create JSON objects while bson.M search JSON objects
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}

		// hash password
		hashedPassword, err := HashPassword(user.Password)

		if err != nil {
			// return 500 if insert fails
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to hash password"})
			return

		}

		// count user email to check the email already exist or not
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil {
			// return 500 if query fails
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user."})
			return
		}

		// check the email already exist or not
		if count > 0 {
			// return 409 if user already exist
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists."})
		}

		// define the data
		user.UserID = bson.NewObjectID().Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.Password = hashedPassword

		// insert data as the result
		result, err := userCollection.InsertOne(ctx, user)

		if err != nil {
			// return 500 if failed to create user
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user."})
		}

		// return user as the result
		c.JSON(http.StatusOK, result)

	}
}
