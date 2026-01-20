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
)

// global mongo collection reference for "movies"
var movieCollection *mongo.Collection = database.OpenCollection("movies")

// global validator instance to validate request body struct tags, it's like zod validator
var validate = validator.New()

// getmovies returns a gin handler that fetches all movies from mongodb
func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {

		// create context with timeout to avoid hanging db queries
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// movies array variable to store the decoded movies
		var movies []models.Movie

		// find all documents in the movies collection (empty filter {} = find all)
		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			// return 500 if query fails
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies."})
			return
		}

		// always close cursor after finished reading results
		defer cursor.Close(ctx)

		// decode all cursor results into the movies slice
		if err = cursor.All(ctx, &movies); err != nil {
			// return 500 if decoding bson into structs fails
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies."})
			return
		}

		// return movies as json array with 200 status
		c.JSON(http.StatusOK, movies)
	}
}

// getmoviebyid returns a gin handler that fetches one movie by imdb_id
func GetMovieByID() gin.HandlerFunc {
	return func(c *gin.Context) {

		// create context with timeout for single document lookup
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// read route param from url path (example: /api/movies/:imdb_id)
		movieID := c.Param("imdb_id")

		// basic check to ensure id exists
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required."})
			return
		}

		// movie variable to store the decoded movie
		var movie models.Movie

		// find one document where imdb_id equals the given param
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
		if err != nil {
			// return 404 if not found or decode fails
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found."})
			return
		}

		// return the found movie as json with 200 status
		c.JSON(http.StatusOK, movie)
	}
}

// addmovie returns a gin handler that inserts a new movie document
func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {

		// create context with timeout for insert operation
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// movie struct that will receive the json request body
		var movie models.Movie

		// bind and parse json body into the movie struct
		if err := c.ShouldBindJSON(&movie); err != nil {
			// return 400 if json is invalid or missing required json fields
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// validate struct based on validate tags in models.Movie
		if err := validate.Struct(movie); err != nil {
			// return 400 if validation fails
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed.",
				"details": err.Error(),
			})
			return
		}

		// insert movie as a new document in mongodb
		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			// return 500 if insert fails
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return
		}

		// return inserted id and insert info with 201 status
		c.JSON(http.StatusCreated, result)
	}
}