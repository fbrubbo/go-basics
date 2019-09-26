package main

import (
	"fmt"
	"strconv"
)

// Person doc here...
type Person struct {
	firstName string
	lastName  string
	city      string
	gender    string
	age       int
	//or
	// firstname, lastname, city, gender string
	// age                               int
}

// Person method - no value change in struct
func (p Person) greet() string {
	return "hello " + p.firstName + " | age: " + strconv.Itoa(p.age)
}

// Person method - change struct value
func (p *Person) hasBirthday() {
	p.age++
}

func main() {
	p1 := Person{firstName: "aa", lastName: "bb", city: "cc", gender: "M", age: 1}
	p2 := Person{"aa2", "bb2", "cc2", "M2", 12}
	fmt.Println(p1, p2)
	fmt.Println(p1.firstName)

	fmt.Println(p1.age)
	p1.age++
	fmt.Println(p1.age)

	fmt.Println(p1.greet())
	p1.hasBirthday()
	fmt.Println(p1.age)
}
