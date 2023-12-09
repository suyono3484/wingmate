package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	log.Println("using log.Println")
	fmt.Println("using fmt.Println")

	time.Sleep(time.Second * 5)

	log.Println("log: finishing up")
	fmt.Println("fmt: finishing up")
}
