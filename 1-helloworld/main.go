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
