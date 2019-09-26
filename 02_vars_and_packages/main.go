package main

import (
	"fmt"
	"math"

	"github.com/fbrubbo/go-basics/02_vars_and_packages/util"
)

const isCool = true

var global = "test"

func main() {
	//Types:
	// bool
	// string
	// int  int8  int16  int32  int64
	// uint uint8 uint16 uint32 uint64 uintptr
	// byte // alias for uint8
	// rune // alias for int32
	// 	// represents a Unicode code point
	// float32 float64
	// complex64 complex128

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
