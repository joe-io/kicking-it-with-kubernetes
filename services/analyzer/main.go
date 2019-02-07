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

	r := gin.Default()

	r.GET("/labels", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"labels": []*LabelResult{
				{
					Label:       "canoe",
					Probability: 0.3231,
				},
				{
					Label:       "lake",
					Probability: 0.2412,
				},
			},
		})
	})

	err := r.Run("0.0.0.0:" + config.Port) // listen and serve
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
