# Go Webservices

There are a lot of options for building webservices in Golang.  In fact, Go has a great HTTP library built-in.

In this tutorial, we are going to be using one of the more popular Go frameworks for building REST APIs.

We'll first explore a basic example, followed by a more in-depth use-case that includes two services that talk to each other.

## Creating a Basic Web Service

The Go framework we will be using is called Gin: https://github.com/gin-gonic/gin

Gin is a light-weight, high-performance web framework that focuses on simplicty, speed, and ease-of-use.

Unfortunately, we won't have time to touch on most of the features that *Gin* provides, but here are some to be aware of:
- Binding Query, Form or Post Data to Structs and performing Validation
- Grouping Routes
- Path and Query Parameters
- Multipart/Urlencoded Form Handling
- File Uploads
- Middleware
- Easy rendering of XML, JSON, YAML, ProtoBuf, and JSONP
- HTML Rendering
- Static file serving
- Graceful restart/stop
- Cookies, Redirects, etc.
- HTTP 2 server push

Rather than going into too much detail up-front, let's just get going with code.

Create a *main.go* file in this directory (1-helloweb) and copy in the following code:

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

If you haven't already installed gin, run the following: 
> cd ..
> go get github.com/codegangsta/gin
> gin -h
> cd 2-helloweb
 
We want leave the current directory so we make sure this gets installed globally, instead of in the current module.

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
- cd root dir
- go get github.com/codegangsta/gin

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

In *gin* the Context is a convenience class that gives you access to the most commonly used features of a Request / Response in one place (as well as full access to the underlying standard Go Request/Response objects).

Go ahead and stop *gin* (e.g. hit CMD+C on Mac or CTL+C Windows).

## Creating the Brand API Service
Now for the real deal!  Let's create our first real service!

We are going to be creating the main API for our Brand Discovery service.

Let's change directories to services/api.

Let's create a basic service that implements the following REST API:

```
GET /recognizer/identification

Parameters:
url : The URL of an image for which we want to identify the brand.

Response: 
{
  "result": "recognized|unrecognized",
  "brand": "Apple",
  "url": "$url"
}
```

First, we'll get a basic service running and returning a hard-coded result.

Copy the following to the empty *main.go* file in services/api:
```go
package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	r := gin.Default()

	r.GET("/recognizer/identification", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"url":    "someimage.jpg",
			"result": "recognized",
			"brand":  "Apple",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	err := r.Run("0.0.0.0:" + port) // listen and serve
	if err != nil {
		log.Fatal(err)
	}
}
```

Let's run it and make sure it works (in the services/api directory, run *gin*):
> gin

Go to the browser and make sure our service is working: http://localhost:3000/recognizer/identification?url=someimage.jpg

You should see the hard-coded response:
> {"brand":"Apple","result":"recognized","url":"someimage.jpg"}

Let's go ahead and stub-out a method for ingesting images for our training model as well.
 
Add the following after the line "r := gin.Default()":

```go
	r.POST("/recognizer/training-image", func(c *gin.Context) {
		c.JSON(200, gin.H{
    		"result": "ingested",
    	})
	})
```

## Refactoring Time

Once again our *main* method is getting a bit messy, let's clean things up.

As we have seen previously, functions in Go are first-class citizens.  

Our web-handlers are currently anonymous functions passed in directly to the r.POST and r.GET methods.

Let's move those handlers into their own file.  

Let's create a file called *handlers.go*

Let's move both of our handlers into that file by copying them over and given the functions an actual name.

Now replace the inline functions in *main.go* with the functions you named.  No need to import anything because the functions are defined in the same package.

At this point why don't you give it a try and see if you can do it without any help!

If you get stuck, you can check below for what it might look like:
*handlers.go*
```go
package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

// Placeholder for calling a service that will use the image to train the model for a specific brand
func trainImage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ingested",
	})
}

// Identify the image the user passes in
func identifyImage(c *gin.Context) {
    c.JSON(200, gin.H{
        "url":    "someimage.jpg",
        "result": "recognized",
        "brand":  "Apple",
    })
}
```

*main.go* 
```go
package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	r := gin.Default()

	r.POST("/recognizer/trainer-image", trainImage)
	r.GET("/recognizer/identification", identifyImage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	err := r.Run("0.0.0.0:" + port) // listen and serve
	if err != nil {
		log.Fatal(err)
	}
}
```

Great, that has cleaned things up quite a bit.

Rather than implementing the brand identification logic directly in this service, we will be relying on another Microservice service to help with this.

Let's look at the *analyzer* service now.

## Creating the Analyzer Service
The *analyizer* service will implement the following REST API:
 
```
GET /brand-score

Parameters:
url : The URL of an image for which we want to look for a brand.

Response: 
{
  "brand": "Apple",
  "probabilty": 0.7
}
```

You'll notice that while this is similar to the main API, the analyzer's job is try and identify a brand and to provide a confidence (probability) score for its identification.

The main *api* service will contain the business logic that determines what confidence level is high enough for the to count as a valid recognition.

We've already created a stub for the analyzer service in services/analyzer/main.go.

Feel free to look it over. It is pretty much what you have seen already.

One thing that you might notice that is different is the bit about configuration.

We'll cover that in moment.

## Calling the Analyzer Service

Now that we have two services, let's see how we go about calling the *analyzer* service from the *api* service

First let's start the *analyzer* service.
* *In a new terminal tab*, go service/analyzer
* run the service (not with gin)
* > go run main.go

Let's double check our service is running, open this URL in your browser: http://localhost:8088/brand-score

Great, now lets add a call to the *analyzer* from the *api* service.

We will be using a *Sling*, a Go HTTP client library specifically designed for making API requests.

https://github.com/dghubble/sling

We can use sling directly, but it our example, let's wrap it in a class that will make the *analyzer* client more reusable.

Let's create a file called analyzer.go in the services/api directory.

Let's first add a *struct* and a constructor for AnalyzerApi:

```go
package main

import (
	"errors"
	"fmt"
	"github.com/dghubble/sling"
	"log"
	"net/http"
)

type AnalyzerApi struct {
	sling *sling.Sling
}

func NewAnalyzerApi(baseUrl string, client *http.Client) *AnalyzerApi {
	return &AnalyzerApi{
		sling: sling.New().Client(client).Base(baseUrl),
	}
}
```

In Go, constructors are simply functions - typically named New[Type] by convention.

In this case, we use the parameters to initialize the default instance of sling for this API.

Now let's add in the call to the "brand-score" endpoint:
```go
func (a *AnalyzerApi) ScoreImage(url string) (*GetScoreResponse, error) {
	req := &GetScoreRequest{Url: url}
	scoreResponse := &GetScoreResponse{}

	res, err := a.sling.New().Get("/brand-score").QueryStruct(req).ReceiveSuccess(scoreResponse)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		log.Print("Error: Status Code: ", res.StatusCode)
		return nil, errors.New(fmt.Sprintf("analyzer-service returned status code: %d ", res.StatusCode))
	}

	return scoreResponse, nil
}
```

The line that does the actual work is here:
```go
a.sling.New().Get("/brand-score").QueryStruct(req).ReceiveSuccess(scoreResponse)
```

The rest of the function is error handling and formatting.


Lastly, we will define the input and response types for the call (you can add them to the end of the file):
```go
type GetScoreRequest struct {
	Url string `json:"url"`
}

type GetScoreResponse struct {
	Brand       string  `json:"brand"`
	Probability float32 `json:"probability"`
}
```

If you look carefully at each of the *structs* we have defined, you will notice something new in the field declarations.

In Go, each field can be followed by a string and by convention, the string typically contains name:"value" pairs that can contain metadata bout the field.

These are often used to define thing like how to map *struct* field names to JSON, as we see here.

Lastly, we need to update *handlers.go* to use our new AnalyizerAPI client:
```go
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

// Placeholder for calling a service that will use the image to train the analyzer for a specific brand
func trainImage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ingested",
	})
}

// Identify the image the user passes in
func identifyImage(c *gin.Context) {
	url := c.Query("url")
	res, err := analyzerApi.ScoreImage(url)

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
``` 

Again, there are a few new things to go over.  We are creating a const named *defaultTrustThreshold* and then defining the Go equivalent of an Enum type.

