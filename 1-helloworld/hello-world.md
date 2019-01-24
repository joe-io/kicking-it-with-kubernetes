# Hello, world!

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
