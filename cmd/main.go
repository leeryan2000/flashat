package main

import "fmt"
import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	fmt.Println("New line added to test branch")

	r.Run() // default on :8080
}
