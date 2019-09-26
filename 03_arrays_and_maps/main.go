package main

import (
	"fmt"
	"strconv"
)

func main() {
	// Arrays = fixed size
	var arr [2]string
	arr[0] = "asdf"
	arr[0] = "234"
	fmt.Println(len(arr))

	// Slices
	slice := []string{"Apple", "Orange", "asfd"}
	slice[0] = "asdf"
	slice = append(slice, "234", "345", "567", "456")
	fmt.Println(slice)
	fmt.Println(len(slice))
	fmt.Println(contains(slice, "Orange"))
	fmt.Println(contains(slice, "123"))
	printSlice(slice)
	s := slice[1:4]
	fmt.Println(s)
	s = slice[:2]
	fmt.Println(s)
	s = slice[3:]
	fmt.Println(s)

	// Maps
	emails := make(map[string]string) // analog to new
	emails["rubbo"] = "rubbo@xy.com"
	fmt.Println(emails)
	delete(emails, "rubbo")

	emails2 := map[string]string{"ana": "ana@fsd.com",
		"mike": "emal.com"}
	fmt.Println(emails2)
	fmt.Println("When key not found is it empty: " + strconv.FormatBool(emails2["asdf"] == ""))
	for k, v := range emails2 {
		fmt.Printf("%s: %s\n", k, v)
	}
	for k := range emails2 {
		fmt.Printf("key: %s\n", k)
	}
	for _, v := range emails2 {
		fmt.Printf("Value: %s\n", v)
	}

	// while
	i := 0
	for i < 10 {
		fmt.Println(i)
		i++
	}

	// while true
	x := 5
	for {
		fmt.Println("Do stuff.", x)
		x += 3
		if x > 25 {
			break
		}
	}
}
func printSlice(s []string) {
	fmt.Printf("len=%d cap=%d %s\n", len(s), cap(s), s)
}

// // Others
// color := "red"
// switch color {
// case "red":
// 	fmt.Println("red")
// case "blue":
// 	fmt.Println("blue")
// default:
// 	fmt.Println("default")
// }

// There is not bult in exist method.. crap
func contains(list []string, el string) bool {
	// there is no while
	for index, value := range list {
		fmt.Printf("Index %d", index)
		if value == el {
			return true
		} else if false {
			// do something
		} else {
			// do another thing
		}
	}
	return false
}
