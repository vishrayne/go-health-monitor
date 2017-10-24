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

/*Task is the parent class for all the monitor tasks*/
type Task interface {
	name() string
	run()
	toJSON() string
	toString() string
	handleError(err error)
}

//Start the function
func Start() {
	fmt.Println("monitor!")
	cpuDetails := newCPU()
	cpuDetails.collect()

	fmt.Println(cpuDetails.toJSON())

	memDetails := newMemory()
	memDetails.collect()
	fmt.Println(memDetails.toJSON())
}

func dealWithError(taskName string, err error) {
	if err != nil {
		log.Fatalf("%sTask failed: %v", taskName, err.Error())
	}
}