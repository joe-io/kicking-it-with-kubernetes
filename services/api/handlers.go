package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

const (
	defaultTrustThreshold = 0.80
)

type IdentificationResult string

const (
	Recognized   IdentificationResult = "recognized"
	UnRecognized IdentificationResult = "unrecognized"
)

// Placeholder for calling a service that will use the image to train the model for a specific brand
func trainImage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ingested",
	})
}

// Identify the image the user passes in
func identifyImage(c *gin.Context) {
	url := c.Query("url")
	res, err := modelApi.ScoreImage(url)

	if err != nil {
		log.Println("Error", err)
		_ = c.AbortWithError(500, err)
		return
	}

	if res.Probability > defaultTrustThreshold {
		c.JSON(200, gin.H{
			"result": Recognized,
			"brand":  res.Brand,
		})
	} else {
		c.JSON(200, gin.H{
			"result": UnRecognized,
		})
	}
}
