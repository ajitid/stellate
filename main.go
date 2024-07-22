package main

import (
	"fmt"
	"log"
)

func main() {
	monitorId, err := getMonitorIdContainingCursor()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(monitorId)
}
