package monit

import (
	"fmt"
	"log"
)

const (
	// Normal status
	Normal = "Normal"
	// Warning status
	Warning = "Warning"
	// Caution status
	Caution = "Caution"
	// Fatal status
	Fatal = "Fatal"
)

//Start the function
func Start() {
	fmt.Println("monitor!")

	cpuDetails := newCPU()
	cpuDetails.collect()
	fmt.Println(cpuDetails.toJSON())

	memDetails := newMemory()
	memDetails.collect()
	fmt.Println(memDetails.toJSON())

	diskDetails := newDisk()
	diskDetails.collect()
	fmt.Println(diskDetails.toJSON())
}

func dealWithError(taskName string, err error) {
	if err != nil {
		log.Fatalf("%sTask failed: %v", taskName, err.Error())
	}
}
