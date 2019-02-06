package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		name := c.Query("name")
		if name == "" {
			name = "world"
		}
		c.JSON(200, gin.H{
			"hello": name,
		})
	})

	port := os.Getenv("PORT")
	err := r.Run("0.0.0.0:" + port) // listen and serve
	if err != nil {
		log.Fatal(err)
	}
}
