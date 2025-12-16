package message

import (
	"fmt"
	"os"
)

func GetLoadAvg() (*LoadAvg, error) {
	// Open the /proc/loadavg file
	f, err := os.Open("/proc/loadavg")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read the load averages from the file
	var loadAvg LoadAvg
	_, err = fmt.Fscanf(f, "%f %f %f %d/%d %d", &loadAvg.Load1, &loadAvg.Load5, &loadAvg.Load15, &loadAvg.Running, &loadAvg.Total, &loadAvg.LastPID)
	if err != nil {
		return nil, err
	}

	return &loadAvg, nil
}

type LoadAvg struct {
	Load1   float64
	Load5   float64
	Load15  float64
	Running int
	Total   int
	LastPID int
}
