package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	AnalyzerEndpoint string `envconfig:"API_ANALYZER_ENDPOINT" default:"http://localhost:8088"`
	Port             string `envconfig:"PORT" default:"8082"`
}

func main() {
	config := loadConfig()
	r := gin.Default()

	r.POST("/social-post", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"id":       "abc-123-def-456",
			"imageUrl": "http://somewhere.com/someimage.jpg",
			"keywords": []string{
				"canoe",
				"lake",
			},
		})
	})

	err := r.Run("0.0.0.0:" + config.Port)
	if err != nil {
		log.Fatal(err)
	}
}

func loadConfig() *Config {
	var config Config
	err := envconfig.Process("api", &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config
}
