package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, day 1 learn go!")
	})

	//	err != nil --> error, nill is like null in javascript, so if its not null == error
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server, err")
	}

}
