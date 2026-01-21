package main

import (
	"fmt"

	controller "github.com/farrasnazhif/moviestream-go/controllers"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, day 1 learn go!")
	})

	router.GET("/api/movies", controller.GetMovies())

	router.POST("/api/auth", controller.RegisterUser())
	router.POST("/api/auth/login", controller.LoginUser())

	//	err != nil --> error, nill is like null in javascript, so if its not null == error
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server, err")
	}

}
