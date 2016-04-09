package fsm

import(
	"fmt"
	"time"
	"../driver"
)

const(
	IDLE = iota
	ELEVATING 
	DOOR_OPEN 
)

type ElevatorState struct{
	fsmState int //State
	floor, dir int
	//destination int 
}


//Initializing FSM 
func (state *ElevatorState) InitFsm() {
	fmt.Println("Currently initializin fsm going to state: IDLE")
	elev.floor = 0
	elev.dir = 0
	//elev.destination = 0
	elev.fsmState = IDLE
}

/*func (state ElevatorState) GetDestination() int {
	return state.destination()
}*/

func (state *ElevatorState) SetDirection(dir int) {
	state.dir = dir 
}

func (state ElevatorState) GetDirection() int {
	return state.dir
}

func (state ElevatorState) Getfloor() int {
	return state.floor
}

func (state *ElevatorState) IDLE() { // we have to use pointer reciever beacause we want to read and write as oposed to just read. 	
	fmt.Println("Going to state: IDLE")
	driver.StopElevate()
	elev.fsmState = IDLE
	//elev.dir = 0
}

func (state *ElevatorState) Elevating(direction int) {
	fmt.Println("Getting in to motion by going to state: ELEVATING")
	state.dir = direction
	if direction == 1 {
		driver.ElevateUp()
	} else {
		driver.ElevateDown()
	}
	state.fsmState = ELEVATING

}

func (state *ElevatorState) DoorOpen() {
	fmt.Println("Going in to state: DOOR_OPEN")
	driver.StopElevate()
	state.fsmState = DOOR_OPEN
	driver.ElevSetDoorOpenLamp(1)
	time.Sleep(4000*time.Millisecond) //The lines under should happend after sleep. 
	driver.ElevSetDoorOpenLamp(0)
	state.IDLE()
}



