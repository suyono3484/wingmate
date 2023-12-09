package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		log.Println("using log.Println")
		fmt.Println("using fmt.Println")
		time.Sleep(time.Millisecond * 500)
	}

	log.Println("log: finishing up")
	fmt.Println("fmt: finishing up")
}
