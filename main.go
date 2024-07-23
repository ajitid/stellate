package main

import (
	"math/rand/v2"
)

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func main() {
	// monitorId, err := getMonitorIdContainingCursor()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// monitorId += "\\Monitor0"

	// brightness := randRange(0, 101)

	// out, err := exec.Command("ControlMyMonitor.exe", "/GetValue", monitorId, "10").Output()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("The date is %s\n", out)

	// fmt.Println(monitorId)

}
