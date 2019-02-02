package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func add(a, b int) int {
	return a + b
}

func divide(a, b int) (quotient, remainder int) {
	quotient = a / b
	remainder = a % b
	return
}

func TestAdd(t *testing.T) {
	assert.Equal(t, 4, add(2, 2))
	assert.Equal(t, 5, add(3, 2))
	quotient, remainder := divide(10, 2)
	assert.Equal(t, 5, quotient)
	assert.Equal(t, 0, remainder)
}

func TestDivide(t *testing.T) {
}

func TestAnonymousFunctions(t *testing.T) {
}

func TestFunctionParameters(t *testing.T) {
}
