package main

import (
	"fmt"
	"sync"
)

// Channels are the pipes that connect concurrent goroutines.
// lots to lear https://gobyexample.com/channels

var wg sync.WaitGroup

func foo(c chan int, someValue int) {
	defer wg.Done()
	c <- someValue * 5
}

func main() {
	fooVal := make(chan int, 10) // must especify the size of the channel to iterate over it
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go foo(fooVal, i)
	}
	wg.Wait()
	close(fooVal)
	for item := range fooVal {
		fmt.Println(item)
	}
}
