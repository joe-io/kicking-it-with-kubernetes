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
