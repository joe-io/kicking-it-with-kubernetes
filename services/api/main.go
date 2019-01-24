package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
)

type Config struct {
	ModelEndpoint string `default:"http://localhost:8088"`
	Port          string `default:"8082"`
}

var modelApi *ModelApi

func main() {
	config := loadConfig()
	modelApi = NewModelApi(config.ModelEndpoint, &http.Client{})

	r := gin.Default()

	r.POST("/recognizer/trainer-image", trainImage)
	r.GET("/recognizer/identification", identifyImage)

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
