package main

import (
	"fmt"
)

// Sum doc..
func Sum(x, y int) int {
	return x + y
}

func main() {
	fmt.Println("test")
}

/*
If we use build contraints/tag (lists the conditions under which a file should be included in the package)

eg.
// +build unit

you must explicitly inform them in the --tags parameter
eg.

go test -v --tags unit_tests
go test -v --tags integration_tests
go test -v --tags unit_tests,integration_tests
*/
