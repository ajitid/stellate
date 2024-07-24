package main

import "fmt"

func main() {
	m := getCursorMonitor()
	fmt.Println(m.getBrightness())
	m.setBrightness(88)
}
