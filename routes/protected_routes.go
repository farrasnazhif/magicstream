package routes

import (
	controller "github.com/farrasnazhif/moviestream-go/controllers"
	"github.com/farrasnazhif/moviestream-go/middleware"
	"github.com/gin-gonic/gin"
)

func ProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleWare())

	router.GET("/api/movie/:imdb_id", controller.GetMovieByID())
	router.POST("/api/movies", controller.AddMovie())
}
