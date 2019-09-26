package main

import (
	"fmt"
	"math"

	"github.com/fbrubbo/go-basics/02_vars_and_packages/util"
)

const isCool = true

var global = "test"

func main() {
	var name string = "fernando"
	var lastname = "rubbo"

	//shorthand does not work outisde of a funciton
	age := 39
	size, weight := 1.87, 83.5

	fmt.Println(name, lastname, age, isCool, global, size, weight)
	fmt.Printf("name type %T\n", name)
	fmt.Printf("weitht %T\n", weight)

	fmt.Println(math.Sqrt(16))

	fmt.Println(util.Reverse("rubbo"))
	fmt.Println(util.Other("aaa"))
}
