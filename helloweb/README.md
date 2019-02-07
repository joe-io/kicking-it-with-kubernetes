# Hello, Web! - REST APIs Using Golang

There are a lot of options for building webservices in Golang.

In this tutorial, we are going to be using one of the more popular Go frameworks for building REST APIs called *Gin*, https://github.com/gin-gonic/gin

We'll first explore a basic example, followed by a more in-depth use-case that includes two services that talk to each other.

## Creating a Basic Web Service

Gin is a light-weight, high-performance web framework that focuses on simplicty, speed, and ease-of-use.

Unfortunately, we won't have time to touch on most of the features that *Gin* provides, but here are some to be aware of:
- Binding Query, Form or Post Data to Structs and performing Validation
- Multipart/Urlencoded Form Handling, File Uploads, etc.
- Grouping Routes
- Easy rendering of XML, JSON, YAML, ProtoBuf, and JSONP
- Graceful restart/stop
- HTTP 2 server push

Rather than going into too much detail up-front, let's just get going with code.

Create a *main.go* file in this directory (helloweb) and copy in the following code:

```go
package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Create a default Gin Engine
	r := gin.Default()

    // Handle the /hello route
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

Let's run our service:
> go run main.go

Open the browser to http://localhost:8282/hello and you should see :
```json
{"hello":"world"}
```

Let's talk a little bit about some new things in the code.  First, you'll notice that the last parameter passed into *r.GET* is a function.

In Go you can create anonymous functions. You can assign them to variables or pass them directly as parameters to other functions.

Having functions be a first-class citizen is quite powerful and useful.

## Handling Query Parameters

Update *main.go* to look for the query parameter 'name'.

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
			name = "world!"
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

Let's restart our service (Ctl+C or Cmd+C) and then:
> go run main.go  

Now let's check it out:
> http://localhost:3000/hello?name=Joe

You probably noticed that we were accessing query parameters through the *Context* variable *c*.

In *gin*, Context is a convenience wrapper that gives you access to the most commonly used features of a Request / Response in one place (as well as full access to the underlying standard Go Request/Response objects).

You can stop the current service (Ctl+C or Cmd+C).

Now that we have our hello-world web service working, let's start building our real services [Building the Services](../services)
