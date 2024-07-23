package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os/exec"
	"strconv"
	"strings"
)

// [min, max)
func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func main() {
	monitorDisplayName, err := cursorOnMonitor()
	if err != nil {
		log.Fatal(err)
	}
	// suffixing with `\\Monitor0` is needed otherwise it'd report some other device ID
	// 0 and 1 are options provided by https://github.com/posthumz/DisplayDevices (see readme)
	display := getDisplayFromName(monitorDisplayName+"\\Monitor0", 0)
	if display == nil {
		log.Fatal(fmt.Errorf("couldn't retrieve display's ID"))
	}

	deviceInstanceID := strings.Split(display.DeviceIDString(), "\\")[1] // no idx checking done
	psCommand := fmt.Sprintf(`Get-WmiObject WmiMonitorID -Namespace root\wmi | Where-Object { $_.InstanceName -like "DISPLAY\%s\*" } | Select-Object -ExpandProperty InstanceName`, deviceInstanceID)
	cmd := exec.Command("powershell", "-Command", psCommand)
	monitorInstanceNameInBytes, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	monitorInstanceName := strings.TrimSpace(string(monitorInstanceNameInBytes))
	monitorInstanceName = strings.TrimSuffix(monitorInstanceName, "_0") // monitorian doesn't has this suffix. Run `monitorian /get all` to check

	brightness := randRange(0, 101)
	cmd = exec.Command("monitorian", "/set", monitorInstanceName, strconv.Itoa(brightness))
	if err := cmd.Run(); err != nil {
		log.Println("Error:", err)
	}
}
