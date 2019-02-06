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

## Live Reload
Up until this point we have been using either "go run" or "go build".

For web development, it can be very helpful to have live reload (your service is rebuilt any time you save changes).

This keeps us from having to stop the service, rebuild, and restart the service.

There is a popular tool for live reloading Go programs called *gin* - to be clear, this is separate from the Gin Web Framework we are using.

If you haven't already installed gin, *open a new terminal window* and run the following:
```bash
> cd ~
> go get github.com/codegangsta/gin
> gin -h
``` 
 
If you don't see the help for *gin*, double check that you ran the *go get* command outside the application directory. 
 
Now you can run it in any directory containing Go code and it will build and watch the code for changes.

But we have to do one more thing for it to work: we need to check an environment variable for the port to listen on in our service:

That way *gin* can listen and proxy to our application.

Let's update *main.go* to the following:

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

Now we can simply run *gin* and it will build our app for us (and rebuild it whenever we change it).
> gin

You'll notice *gin* starts up and prints out the port it is listening on (probably 3000).

Let's hit our endpoint: http://localhost:3000/hello

## Handling Query Parameters

Let's try out our new live-reload powers.
 
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

You won't need to stop or start anything.  

Simply save the changes and update the url in the browser to include your name: http://localhost:3000/hello?name=Joe

You probably noticed that we were accessing query parameters through the *Context* variable *c*.

In *gin*, Context is a convenience wrapper that gives you access to the most commonly used features of a Request / Response in one place (as well as full access to the underlying standard Go Request/Response objects).

Now that we have our hello-world web service working, let's start building our real services.

First, let's stop *gin* from running our hello-world web service (e.g. hit CMD+C on Mac or CTL+C Windows).

Now we are ready to apply what we have learned to build the actual services: [Building the Services](../services)
