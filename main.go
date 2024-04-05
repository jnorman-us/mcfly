package main

import (
	"fmt"
	"time"
)

func main() {
	for range time.Tick(time.Second) {
		fmt.Println("Hello world!")
	}
}
