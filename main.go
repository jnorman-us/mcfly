package main

import (
	"fmt"
	"time"
)

func main() {
	for range time.Tick(time.Minute) {
		fmt.Println("Hello world!")
	}
}
