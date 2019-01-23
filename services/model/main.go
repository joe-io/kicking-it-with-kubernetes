package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	Port string `default:"8088"`
}

func main() {
	config := loadConfig()

	r := gin.Default()

	r.GET("/score-image", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"brand":       "Apple",
			"probability": 0.925,
		})
	})

	err := r.Run("0.0.0.0:" + config.Port) // listen and serve
	if err != nil {
		log.Fatal(err)
	}
}

func loadConfig() *Config {
	var config Config
	err := envconfig.Process("model", &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config
}
