package fsm

import (
	. "../definitions"
	"../driver"
	. "../network"
	"fmt"
	"time"
)

//Initializing FSM
func InitFsm() {
	fmt.Println("Currently initializin fsm going to state: IDLE")
	Elevators[GetLocalIP()] = &Elevator{true, 0, 0, IDLE, GetLocalIP(), [4]bool{false, false, false, false}, [4]bool{false, false, false, false}, [4]bool{false, false, false, false}}
}

func GoToIDLE() {
	driver.StopElevate()
	Elevators[GetLocalIP()].FsmState = IDLE
	Elevators[GetLocalIP()].Direction = 0
}

func GoToElevating(direction int) {
	fmt.Println("Getting in to motion by going to state: ELEVATING")
	Elevators[GetLocalIP()].Direction = direction
	fmt.Println("The Direction is:", Elevators[GetLocalIP()].Direction)
	if direction == 1 {
		Elevators[GetLocalIP()].FsmState = ELEVATING
		driver.ElevateUp()
	} else if direction == -1 {
		Elevators[GetLocalIP()].FsmState = ELEVATING
		driver.ElevateDown()
	} else {
		time.Sleep(200 * time.Millisecond)
		GoToIDLE()
	}

}

func GoToDoorOpen() {
	fmt.Println("Going in to state: DOOR_OPEN")
	driver.StopElevate()
	Elevators[GetLocalIP()].FsmState = DOOR_OPEN
	driver.ElevSetDoorOpenLamp(1)
	time.Sleep(2 * time.Second) //The lines under should happend after sleep.
	driver.ElevSetDoorOpenLamp(0)
	GoToIDLE()
}
