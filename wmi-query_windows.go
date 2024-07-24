package main

import (
	"fmt"
	"strings"

	"github.com/yusufpapurcu/wmi"
)

type WmiMonitorID struct {
	InstanceName string
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
	monitorInstanceName := strings.TrimSuffix(monitors[0].InstanceName, "_0") // monitorian doesn't use this suffix
	return monitorInstanceName, nil
}
