package fsm

import(
	"fmt"
	"../driver"
	"time"
)

const(
	IDLE = iota
	ELEVATING 
	DOOR_OPEN 
)

type State int

type ElevatorState struct{
	fsmState State
	floor, dir, destination int 
}

func (elev *ElevatorState) IDLE() { // we have to use pointer reciever beacause we want to read and write as oposed to just read. 	
	fmt.Println("Going to state: IDLE")
	driver.StopElevate()
	elev.fsmState = IDLE
}

func (elev *ElevatorState) Elevating(direction int) {
	fmt.Println("Getting in to motion by going to state: ELEVATING")
	if direction == 1 {
		driver.ElevateUp()
	} else {
		driver.ElevateDown()
	}
	elev.fsmState = ELEVATING

}

func (elev *ElevatorState) DoorOpen() {
	fmt.Printline "Going in to state: DOOR_OPEN"
	driver.StopElevate()
	elev.fsmState = DOOR_OPEN
	driver.elevSetDoorOpenLamp(1)
	time.sleep(400*time.Milliseconds) //The lines under should happend after sleep. 
	driver.elevSetDoorOpenLamp(0)
	elev.IDLE()
}



