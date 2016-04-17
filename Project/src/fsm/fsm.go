package fsm

import (
	. "../definitions"
	"../driver"
	//. "../network"
	"fmt"
	"time"
)

func GoToIDLE() {
	driver.StopElevate()
	Elevators[LocalIP].Direction = 0
}

func GoToElevating(direction int) {
	Elevators[LocalIP].Direction = direction
	if direction == 1 {
		driver.ElevateUp()
	} else if direction == -1 {
		driver.ElevateDown()
	} else {
		time.Sleep(200 * time.Millisecond)
		GoToIDLE()
	}

}

func GoToDoorOpen() {
	fmt.Println("The doors are opening")
	driver.StopElevate()
	driver.SetDoorLamp(1)
	time.Sleep(2 * time.Second)
	fmt.Println("The doors are closing")
	driver.SetDoorLamp(0)
	GoToIDLE()
}
