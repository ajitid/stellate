package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/yusufpapurcu/wmi"
)

type WmiMonitorID struct {
	InstanceName string
}

type SimpleWmiMonitorBrightness struct {
	InstanceName string
}

type WmiMonitorBrightness struct {
	InstanceName      string
	CurrentBrightness int
}

func getMonitorInstanceName(deviceInstanceId string) (string, error) {
	var monitors []WmiMonitorID
	query := fmt.Sprintf(`SELECT InstanceName FROM WmiMonitorID WHERE InstanceName LIKE "DISPLAY\\%s\\%%"`, deviceInstanceId) // %% just mean a literal % ref. https://yourbasic.org/golang/fmt-printf-reference-cheat-sheet/
	err := wmi.QueryNamespace(query, &monitors, "root\\wmi")
	if err != nil {
		return "", fmt.Errorf("error querying WMI: %v", err)
	}
	if len(monitors) != 1 {
		return "", fmt.Errorf("cannot determine monitor")
	}
	return monitors[0].InstanceName, nil
}

// While WMIMonitorID may list all types of monitors, only WmiMonitorBrightness and WmiMonitorBrightnessMethods will the monitor it can control brightness via WMI.
// Usually this type of monitor is integrated laptop display
func isTypeWMIMonitor(monitorInstanceName string) (bool, error) {
	var monitors []SimpleWmiMonitorBrightness
	// DISPLAY\SHP1523\5&14db058f&2&UID512_0 needs to be changed to DISPLAY\\SHP1523\\5&14db058f&2&UID512_0
	// pwsh> Get-WmiObject -Query "SELECT * FROM WmiMonitorBrightness" -Namespace root\wmi
	query := fmt.Sprintf(`SELECT InstanceName FROM WmiMonitorBrightness WHERE InstanceName="%s"`, strings.ReplaceAll(monitorInstanceName, "\\", "\\\\"))
	err := wmi.QueryNamespace(query, &monitors, "root\\wmi")
	if err != nil {
		return false, fmt.Errorf("error querying WMI: %v", err)
	}
	if len(monitors) == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

// https://stackoverflow.com/a/62634211/7683365
// usually integrated laptop displays
type WMIMonitor struct {
	name string
	pos  rl.Vector2
}

func (m WMIMonitor) getInstanceName() string {
	return m.name
}

func (m WMIMonitor) getPosition() rl.Vector2 {
	return m.pos
}

func (m WMIMonitor) getBrightness() int {
	var monitors []WmiMonitorBrightness
	// DISPLAY\SHP1523\5&14db058f&2&UID512_0 needs to be changed to DISPLAY\\SHP1523\\5&14db058f&2&UID512_0
	query := fmt.Sprintf(`SELECT CurrentBrightness FROM WmiMonitorBrightness WHERE InstanceName="%s"`, strings.ReplaceAll(m.getInstanceName(), `\`, `\\`))
	err := wmi.QueryNamespace(query, &monitors, "root\\wmi")
	if err != nil {
		log.Fatalf("error querying WMI: %v", err)
	}
	if len(monitors) != 1 {
		log.Fatal("monitor not found")
	}
	return monitors[0].CurrentBrightness
}

// Ref.
// https://superuser.com/a/1781874
// I wanted to use this but I couldn't figure it out
// https://github.com/StackExchange/wmi/pull/45#issuecomment-590396746
func (m WMIMonitor) setBrightness(value int) {
	cmd := exec.Command("wmic", `/NAMESPACE:\\root\wmi`, "PATH", "WmiMonitorBrightnessMethods",
		"WHERE", fmt.Sprintf("Active=TRUE AND InstanceName='%s'", strings.ReplaceAll(m.getInstanceName(), `\`, `\\`)),
		"CALL", "WmiSetBrightness", fmt.Sprintf("Brightness=%d", value), "Timeout=0")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
