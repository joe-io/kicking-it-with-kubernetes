# Hello, web! - Build your first REST API with Go

## Creating a Basic Web Service

```go
package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"hello": "world",
		})
	})

	err := r.Run("0.0.0.0:8282") // listen and serve
	if err != nil {
		log.Fatal(err)
	}
}
```

- build like normal

## Adding Auto-Reload
- global install gin

```go
package main

import (
	"github.com/gin-gonic/gin"
	"log"
    "os"
)

func main() {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"hello": "world",
		})
	})

    port := os.Getenv("PORT")
	err := r.Run("0.0.0.0:" + port) // listen and serve
	if err != nil {
		log.Fatal(err)
	}
}
```
- cd root dir
- go get github.com/codegangsta/gin


## Handling Query Parameters

```go
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
```

## Creating the Analyzer Service

## Creating the Brand API Service

## Calling the Analyzer Service




