package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	defaultTrustThreshold = 0.80
)

type PostRequest struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	ImageUrl string `json:"imageUrl"`
}

// Identify the image the user passes in
func handlePost(c *gin.Context) {
	var json PostRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := analyzerApi.AnalyzeImage(json.ImageUrl)

	if err != nil {
		log.Println("Error", err)
		_ = c.AbortWithError(500, err)
		return
	}

	keywords := []string{}

	for _, lr := range res.Labels {
		if lr.Probability >= defaultTrustThreshold {
			keywords = append(keywords, lr.Label)
		}
	}

	c.JSON(200, gin.H{
		"id":       "abc-123-def-456",
		"url":      json.ImageUrl,
		"keywords": keywords,
	})
}
