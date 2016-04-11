package fsm

import (
	"../driver"
	//"definitions"
	"fmt"
	"time"
)

const (
	IDLE = iota
	ELEVATING
	DOOR_OPEN
)

type ElevatorState struct {
	fsmState   int //State
	floor, dir int
	//destination int
}

//Initializing FSM
func (state *ElevatorState) InitFsm() {
	fmt.Println("Currently initializin fsm going to state: IDLE")
	state.floor = 0
	state.dir = 0
	//elev.destination = 0
	state.fsmState = IDLE
}

/*func (state ElevatorState) GetDestination() int {
	return state.destination()
}*/

func (state *ElevatorState) SetDirection(dir int) {
	state.dir = dir
}

func (state *ElevatorState) Setfloor(f int) {
	state.floor = f
}

func (state ElevatorState) GetDirection() int {
	return state.dir
}

func (state ElevatorState) GetFloor() int {
	return state.floor
}

func (state ElevatorState) GetState() int {
	return state.fsmState
}

func (state *ElevatorState) IDLE() { // we have to use pointer reciever beacause we want to read and write as oposed to just read.
	//fmt.Println("Going to state: IDLE")
	driver.StopElevate()
	state.fsmState = IDLE
	state.dir = 0
}

func (state *ElevatorState) Elevating(direction int) {
	fmt.Println("Getting in to motion by going to state: ELEVATING")
	state.dir = direction
	fmt.Println("The Direction is:", state.dir)
	if direction == 1 {
		state.fsmState = ELEVATING
		driver.ElevateUp()
	} else if direction == -1 {
		state.fsmState = ELEVATING
		driver.ElevateDown()
	} else {
		time.Sleep(200 * time.Millisecond)
		state.IDLE()
	}

}

func (state *ElevatorState) DoorOpen() {
	fmt.Println("Going in to state: DOOR_OPEN")
	driver.StopElevate()
	state.fsmState = DOOR_OPEN
	driver.ElevSetDoorOpenLamp(1)
	time.Sleep(2 * time.Second) //The lines under should happend after sleep.
	driver.ElevSetDoorOpenLamp(0)
	state.IDLE()
}
