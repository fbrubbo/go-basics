package main

import "fmt"

func main() {
	// everything in go is pass by value (copy)
	value := 5
	pointer := &value // poiter to value

	fmt.Println(value, pointer)
	fmt.Printf("Type: %T", pointer)
	fmt.Println(*pointer) // read the pointer value

	*pointer = 10 // change the value of the pointer
	fmt.Println(value)
}
