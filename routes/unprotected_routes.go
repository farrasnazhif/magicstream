package routes

import (
	controller "github.com/farrasnazhif/moviestream-go/controllers"
	"github.com/gin-gonic/gin"
)

func UnprotectedRoutes(router *gin.Engine) {
	router.GET("/api/movies", controller.GetMovies())
	router.POST("/api/auth", controller.RegisterUser())
	router.POST("/api/auth/login", controller.LoginUser())
}
