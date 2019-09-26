package main

import (
	"fmt"
	"math"
)

// Shape docs
type Shape interface {
	area() float64
}

// Circle docs
type Circle struct {
	x, y, radius float64
}

func (c Circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

// Retangle docs
type Retangle struct {
	width, height float64
}

func (r Retangle) area() float64 {
	return r.width * r.height
}

func getArea(s Shape) float64 {
	return s.area()
}

func main() {
	c := Circle{x: 0, y: 0, radius: 5}
	r := Retangle{width: 10, height: 5}
	fmt.Println(getArea(c))
	fmt.Println(getArea(r))
}
