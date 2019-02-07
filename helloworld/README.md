# Hello, world! - Intro to Golang

## Overview

The Go programming language is well known for it's simplicity, speed, and novel approach to handling concurrency.

In this part of the lab we will walk through the basic language features and syntax of Go by building a simple command-line program.
  
## Hello, world!

Let's write our first Go program!

First, create a file called main.go (in the helloworld directory).

Next, add the following content: 

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")
}
```

Last, lets run our program:

> go run main.go

You should see the following output:
> Hello, world!

Let's talk about what just happened.

Every go source file belong to a specific package.  If we are making an executable (as opposed to a library), that package should be called 'main'.

The package declaration should always be the first line in the file.

Executable packages should contain a function named main.  This is the entry-point for your program.

Let's make this a bit more interesting.

Let's update the code to ask for a name and then greet the user with their name.

```go
package main

import "fmt"

func main() {
	fmt.Print("Enter a name: ")
	
	var name string
	fmt.Scanln(&name)
	
	fmt.Printf("Hello, %s!\n", name)
}
```

Let's run our program again:
> go run main.go

Now we should be prompted to enter a name, after which we expect to see something like this:
> Hello, Joe!

This new example, introduces some new syntax: *&name* in the call to fmt.Scanln.  In Go, the *&* is an address operator that generates a pointer to a variable.

You can see that we also first declared the variable *name* on this line 
```go 
var name string
```

Great, now allow the user to pass in the name as a command line argument.

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	var name string
	
	if len(os.Args) > 1 {
	    name = os.Args[1]	
	} else {
	    fmt.Print("Enter a name: ")	
	    fmt.Scanln(&name)
    }
	
	fmt.Printf("Hello, %s!\n", name)
}
```

Let's run our program again:
> go run main.go Paul

Which should print:
> Hello, Paul!

One thing to note before we move on is the syntax for the import.  While it is technically valid to have a separate import line for each package you want to import, go supports the more concise form where parenthesis are added and each line contains separate package to import.

You'll see this same pattern used in other places as well, like *const* and *var*, which we will talk about later.

## Refactoring for the Greater Good

Our *main* function feels like it could use some refactoring.

Let's extract the code that gets the user's name into its own method:

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	name := askName()
	fmt.Printf("Hello, %s!\n", name)
}

func askName() string {
	var name string

	if len(os.Args) > 1 {
		name = os.Args[1]
	} else {
		fmt.Print("Enter a name: ")
		fmt.Scanln(&name)
	}
	return name
}
```

Let's run our program again to make sure we haven't broken anything:
> go run main.go

Great, now we are going to take it to the next level!

## Error Handling

Let's ask our user for their favorite number and echo it back.

First we'll create a new function called *askFavoriteNumber*, which we expect to return an integer. 

Add the following to the end of the main.go file:

```go
func askFavoriteNumber() int {
	fmt.Print("Enter your favorite number: ")
	var favNumberText string
	fmt.Scanln(&favNumberText)
	return favNumberText
}
```

Now let's try to run our program again:
> go run main.go

Doh, we got an error.  Looks like we need to convert our string to an integer.

Luckily go has a built-in library for that. Update *askFavoriteNumber* to the following:

```go
func askFavoriteNumber() int {
	fmt.Print("Enter your favorite number: ")
	var favNumberText string
	fmt.Scanln(&favNumberText)
	favNumber, _ := strconv.Atoi(favNumberText)
	return favNumber
}
```

At this point you may still get an error - this time about *strconv*.  This is a built-in library, but you will need to import it.

Simply add it to the import section at the top (if you are using an IDE it might have done that for you).

```go
import (
	"fmt"
	"os"
	"strconv"
)
```  

Now we are converting to an integer.  But the call to *strconv.Atoi* probably looks a little strange.  What is up with the "_"?

We'll talk more about this in just a bit, but the short answer is that Go actually allows functions to return more than one value.

And _ is a placeholder variable name that indicates that you want to ignore that return value.

So in this case we are using the first return value and ignoring the second one (we'll fix that in a minute).

Let's run our program again:
> go run main.go

Everything compiles, but it isn't asking us the question.

I guess we'd better add a call to *askFavoriteNumber* to the *main* function: 

```go
func main() {
	name := askName()
	favNumber := askFavoriteNumber()
	
	fmt.Printf("Hello, %s!\n", name)		
	fmt.Printf("Your favorite number is: %d\n", favNumber)
}
```

Let's run our program again:
> go run main.go

At this point, this is what *main.go* looks like in its entirety: 
```go
package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	name := askName()
	favNumber := askFavoriteNumber()
	
	fmt.Printf("Hello, %s!\n", name)		
	fmt.Printf("Your favorite number is: %d", favNumber)
}

func askName() string {
	var name string

	if len(os.Args) > 1 {
		name = os.Args[1]
	} else {
		fmt.Print("Enter a name: ")
		fmt.Scanln(&name)
	}
	return name
}

func askFavoriteNumber() int {
	fmt.Print("Enter your favorite number: ")
	var favNumberText string
	fmt.Scanln(&favNumberText)
	favNumber, _ := strconv.Atoi(favNumberText)
	return favNumber
}
```

Some of you may be thinking what if I enter something that isn't a number for my favorite number, like "hi world"?

Let's try it out.  Hmm, it doesn't blow up, but we didn't enter "0" as our favorite number.

Let's take a look at *askFavoriteNumber* again.  Remember how we were ignoring that second return value from *strconv.Atoi*?

Well, that value is of type *error*.  Go does not support traditional exception handling like Java or C#.  Instead, Go functions that might produce an error will return *error* as the last return type.

Let's update our *askFavoriteNumber* function to properly handle the error:   

```go
func askFavoriteNumber() (int, error) {
	fmt.Print("Enter your favorite number: ")
	
	var favNumberText string
	fmt.Scanln(&favNumberText)	
	favNumber, err := strconv.Atoi(favNumberText)
	return favNumber, err
}
```

Here we have updated *askFavoriteNumber* to return two types, the number and an error.

Since that is the same set of return-types as *strconv.Atoi* we can further simplify the function by just returning the result of the *strconv.Atoi* call:

```go
func askFavoriteNumber() (int, error) {
	fmt.Print("Enter your favorite number: ")
	
	var favNumberText string
	fmt.Scanln(&favNumberText)	
    return strconv.Atoi(favNumberText)	
}
```

Now let's handle the error in our *main* function:

```go
func main() {
	name := askName()
	
	favNumber, err := askFavoriteNumber()
	if err != nil {
		log.Fatal("Invalid favorite number.")
	}
	
	fmt.Printf("Hello, %s!\n", name)
	fmt.Printf("Your favorite number is: %d\n", favNumber)
}
```

## A Rose by Any Other Name

Unlike many Object Oriented languages, Go does not support classes or inheritance.  It does however have *structs*, which can have methods.

The syntax for a *struct* is pretty simple:

```go
type Person struct {
	FirstName string
	LastName string
	Age int
}
```

Methods can be added to a *struct* like so:

```go
type Person struct {
	FirstName string
	LastName string
	Age int
}

func (p Person) FullName() string {
	return p.FirstName + p.LastName
} 
```

A method is a function that declares a receiver.  The receiver is the *type* that the method is attached to.

## It's All About the User

Let's put what we just learned to practice by adding a User struct to our program.

Unlike Java, Go doesn't care which files your code lives in.  You could have 10 functions and 2 structs in a one file or split across 12 files - Go doesn't care.

Let's add the following the the bottom of *main.go*:

```go
type User struct {
	Name string
	FavoriteNumber int
}
```

Now let's update our *main* function to use this new type:

```go
func main() {
	name := askName()
	
	favNumber, err := askFavoriteNumber()
	if err != nil {
		log.Fatal("Invalid favorite number.")
	}
	
	user := &User{
		Name: name,
		FavoriteNumber: favNumber,
	}
	
	fmt.Printf("Hello, %s!\n", user.Name)
	fmt.Printf("Your favorite number is: %d\n", user.FavoriteNumber)
}
```

Let's further enhance the User type by adding a method to handle the user greeting.

Add the following the end of *main.go*:
```go
func (user *User) Greet() {
}
```

Now let's move the greeting from *main* to this function:
```go
func (user *User) Greet() {
	fmt.Printf("Hello, %s!\n", user.Name)
    fmt.Printf("Your favorite number is: %d\n", user.FavoriteNumber)
}
```

And add a call to user.Greet to *main*:
```go
func main() {
	name := askName()
	
	favNumber, err := askFavoriteNumber()
	if err != nil {
		log.Fatal("Invalid favorite number.")
	}
	
	user := &User{
		Name: name,
		FavoriteNumber: favNumber,
	}
	
	user.Greet()
}
```

Let's run our program:
> go run main.go

## The Prime Directive

We are going to add one more enhancement to our program before moving on to building our first web-service in Go.

Let's tell the user if their favorite number is a prime number or not.

But first, we need to do some housekeeping.

Our *main.go* is getting a bit long in the tooth.  Let's add a new file called *user.go* and move the User struct and it's method into this file.

You'll need to start the file with the package declaration and correct import statements.

Here is the complete *user.go* file:

```go
package main

import "fmt"

type User struct {
	Name           string
	FavoriteNumber int
}

func (user *User) Greet() {
	fmt.Printf("Hello, %s!\n", user.Name)
	fmt.Printf("Your favorite number is: %d\n", user.FavoriteNumber)
}
```

Now let's try running our program again, to make sure it works:
> go run main.go

You should be seeing an error (if not you may have not removed User from *main.go*):
> ./main.go:18:11: undefined: User

The problem is actually simple to solve.  "go run" is only meant to run a single file.

Now that we have more than one file, we'll want to run "go build" and then run the built executable.

Let's try that now:
> go build && ./helloworld

You'll notice that we didn't create any project file or have to list the files to include in our built output.

Go automatically takes care of all of the resolution for us, including .go files that are in the same directory (package) and ones that are imported from other packages.

All you ever need to do is run "go build" and your package will be built.

Now, let's add the code to see if the user's favorite number is a prime number.

Add the following method to the end of *user.go*:
```go
func (user *User) checkFavoriteNumber() {
	bigFavNumber := big.NewInt(int64(user.FavoriteNumber))
	if bigFavNumber.ProbablyPrime(0) {
		fmt.Printf("Your favorite number is prime!\n")
	} else {
		fmt.Printf("Love your number, but... it is not prime.\n")
	}
}
```

Now let's call the method from *user.Greet*:
```go
func (user *User) Greet() {
	fmt.Printf("Hello, %s!\n", user.Name)
	user.checkFavoriteNumber()
}
```

Let's build and run the program to see it work:
> go build && ./helloworld

Congratulations, you have built a go command line-program that users *structs*, STDIN, STDOUT, and error handling.

We are now ready to build our first web-service in Go!


## TL;DR
If you ran into an issue or just wanted to see the final results, here they are:

Create a file called *main.go* in this directory that contains the following:
```go
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	name := askName()

	favNumber, err := askFavoriteNumber()
	if err != nil {
		log.Fatal("Invalid favorite number.")
	}

	user := &User{
		Name:           name,
		FavoriteNumber: favNumber,
	}

	user.Greet()
}

func askName() string {
	var name string

	if len(os.Args) > 1 {
		name = os.Args[1]
	} else {
		fmt.Print("Enter a name: ")
		fmt.Scanln(&name)
	}
	return name
}

func askFavoriteNumber() (int, error) {
	fmt.Print("Enter your favorite number: ")

	var favNumberText string
	fmt.Scanln(&favNumberText)
	return strconv.Atoi(favNumberText)
}
```

Create another file called *user.go* in this directory that contains the following:
```go
package main

import (
	"fmt"
	"math/big"
)

type User struct {
	Name           string
	FavoriteNumber int
}

func (user *User) Greet() {
	fmt.Printf("Hello, %s!\n", user.Name)
	user.checkFavoriteNumber()
}

func (user *User) checkFavoriteNumber() {
	bigFavNumber := big.NewInt(int64(user.FavoriteNumber))
	if bigFavNumber.ProbablyPrime(0) {
		fmt.Printf("Your favorite number is prime!\n")
	} else {
		fmt.Printf("Love your number, but... it is not prime.\n")
	}
}
```

Build and run the program by running the following in this directory:
> go build && ./helloworld

Now that you have a basic introduction to Golang, let's take a look at web-services in Golang: [Hello, Web! - REST APIs Using Golang](../helloweb)
