package main

import (
	"fmt"
	"time"

	"github.com/pozedorum/or"
)

func main() {
	ch1 := make(chan any)
	ch2 := make(chan any)
	orChan := or.Or(ch1, ch2)

	// В другой горутине
	go func() {
		time.Sleep(time.Second)
		close(ch1)
	}()

	// Будет получено значение, когда закроется ch1
	<-orChan
	fmt.Println("One of the channels closed!")
}
