package main

import (
	"fmt"
	"log"
)

func main() {
	m, err := cursorOnMonitor()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(m)
}
