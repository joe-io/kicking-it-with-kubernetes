# Hello, world! - Intro to Golang

## Overview

The Go programming language is well known for it's simplicity, speed, and novel approach to handling concurrency.

    In this lab we will walk through some of the key language features that make Go unique and give you hands on experience working with them.

This lab is divided into 10 sections, each one covering a different area of the language.

In most sections, rather than writing a complete program, we will be exploring the language by writing code in a unit-test harness.

This will allow us to quickly explore different language features in a clean and simple way.
 
## Sections 

 - [Hello, world!](#hello-world)
 - [Functions](#functions)
 - [Structs and Interfaces](#structs-and-interfaces)
 - [Arrays, Slices, and Maps](#arrays-slices-and-maps)
 - [Control Structures](#control-structures)
 - [Concurrency](#concurrency)
 - [Error Handling](#error-handling)
  
## Hello, world!

Let's write our first go program!

First, create a file called main.go

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

Let's update the code to ask for a name and use it in our greeting.

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

Let's talk about one strange thing we saw, *&name* in the call to fmt.Scanln.  In Go, the *&* is an address operator that generates a pointer.

For now, just know that Go supports pointers, but not pointer arithmetic.  We'll get more into that later.

Let's add one more feature.  Let's allow the user to pass in the name as a command line argument.

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

Now that we have created a program that can read from STDIN, write to STDOUT and read command line arguments, we are ready to dig into some of the meatier features that make go so unique!

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

```go
package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
)

func main() {
	name := askName()
	fmt.Printf("Hello, %s!\n", name)
	
	favNumber := askFavoriteNumber()
	bigFavNumber := big.NewInt(int64(favNumber))
	if bigFavNumber.ProbablyPrime(0) {
		fmt.Printf("Your favorite number is prime!")
	} else {
		fmt.Printf("Love your number, but... it is not prime.")
	}
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

```go
package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
)

func main() {
	name, err := askName()
	if err != nil {
		log.Fatal(err)
	}

	favNumber, err := askFavoriteNumber()
	if err != nil {
		log.Fatal(err)
	}
	bigFavNumber := big.NewInt(int64(favNumber))

	fmt.Printf("Hello, %s!\n", name)
	if bigFavNumber.ProbablyPrime(0) {
		fmt.Printf("Your favorite number is prime!")
	} else {
		fmt.Printf("Love your number, but... it is not prime.")
	}
}

func askName() (string, error) {
	var name string

	if len(os.Args) > 1 {
		name = os.Args[1]
	} else {
		fmt.Print("Enter a name: ")
		_, err := fmt.Scanln(&name)
		if err != nil {
			return "", err
		}
	}
	return name, nil
}

func askFavoriteNumber() (int, error) {
	fmt.Print("Enter your favorite number: ")
	var favNumberText string
	_, err := fmt.Scanln(&favNumberText)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(favNumberText)
}
```

```go
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	name, err := askName()
	if err != nil {
		log.Fatal(err)
	}

	favNumber, err := askFavoriteNumber()
	if err != nil {
		log.Fatal(err)
	}

	user := &User{
		Name:           name,
		FavoriteNumber: favNumber,
	}

	user.Greet()
}

func askName() (string, error) {
	var name string

	if len(os.Args) > 1 {
		name = os.Args[1]
	} else {
		fmt.Print("Enter a name: ")
		_, err := fmt.Scanln(&name)
		if err != nil {
			return "", err
		}
	}
	return name, nil
}

func askFavoriteNumber() (int, error) {
	fmt.Print("Enter your favorite number: ")
	var favNumberText string
	_, err := fmt.Scanln(&favNumberText)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(favNumberText)
}
```

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

func (u *User) Greet() {
	fmt.Printf("Hello, %s!\n", u.Name)
	u.checkFavoriteNumber()
}

func (u *User) checkFavoriteNumber() {
	bigFavNumber := big.NewInt(int64(u.FavoriteNumber))
	if bigFavNumber.ProbablyPrime(0) {
		fmt.Printf("Your favorite number is prime!\n")
	} else {
		fmt.Printf("Love your number, but... it is not prime.\n")
	}
}
```


## Functions
Functions in Go are very similar to functions in many other languages, with some key differences.

Let's see some of these differences by writing some code.

First open functions_test.go

- Add
    - Parameter name before type
    - Combined types
- Divide
    - Multiple Return types
    - Named return types
- Anonymous functions
- Functions as parameters

## Structs and Interfaces

While adopting some of ideas that made Object Oriented programming so popular, Go takes the concepts in a unique direction.

For example, in Java or C# you have Classes and Inheritance. These do not exist in Go.

However there are similar constructs, but with some unique twists to them (that arguably improve them).

Lets get started by introducing Structs.

Structs are a lot like classes.  They hold data and can even have methods.

Add the following to the structs_test.go file:

```go
type Employee struct {
	Name string
	HireDate time.Time
}
```


## Arrays, Slices, and Maps

- Arrays
    - Read Only, can't be changed
    - How to create
- Slices
    - Pointer to array and start/end
    - Built-in functions like append
- Maps
    - Typed Key and value
- Literal creation of Structs, Arrays, & Maps    

## Concurrency

- goroutines
- channels

## Error Handling

- returning errors
- checking errors
- defer
- Go 2.0 Draft of check

