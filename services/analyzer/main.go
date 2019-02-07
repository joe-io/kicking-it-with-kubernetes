package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	Port string `default:"8088"`
}

// https://obrien.com/shop/life-jackets-vests/ - Life vest company
// https://qvxfxxo9ak-flywheel.netdna-ssl.com/wp-content/uploads/2018/03/Jasper-canoe-tour-at-Pyramid-Lake.jpg
// https://boygeniusreport.files.wordpress.com/2016/11/puppy-dog.jpg?quality=98&strip=all&w=782
// https://www.parksmarina.com/webres/Image/obw/page-top-images/rentals-boat-slips.jpg

func main() {
	err := loadModel()
	if err != nil {
		log.Fatal(err)
	}
	err, result := scoreImage("https://qvxfxxo9ak-flywheel.netdna-ssl.com/wp-content/uploads/2018/03/Jasper-canoe-tour-at-Pyramid-Lake.jpg")
	if err != nil {
		log.Fatal(err)
	}
	for _, lr := range result.Labels {
		fmt.Printf("label: %s, prop: %f\n", lr.Label, lr.Probability)
	}
}

func _main() {
	config := loadConfig()

	r := gin.Default()

	r.GET("/brand-score", func(c *gin.Context) {
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
	err := envconfig.Process("analyzer", &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config
}
