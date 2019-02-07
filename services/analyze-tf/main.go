package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	Port string `envconfig:"PORT" default:"8088"`
}

func main() {
	config := loadConfig()

	err := loadModel()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/labels", func(c *gin.Context) {
		url := c.Query("url")
		result, err := classifyImage(url)
		if err != nil {
			_ = c.AbortWithError(500, err)
		} else {
			c.JSON(200, result)
		}
	})

	err = r.Run("0.0.0.0:" + config.Port) // listen and serve
	if err != nil {
		log.Fatal(err)
	}
}
func loadConfig() *Config {
	var config Config
	err := envconfig.Process("analyzer", &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config
}
