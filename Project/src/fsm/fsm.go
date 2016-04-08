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

//type State int

type ElevatorState struct{
	fsmState int //State
	floor, dir, destination int 
}


func InitFsm() ElevatorState {
	var elev ElevatorState
	fmt.Println("Currently initializin fsm going to state: IDLE")
	elev.floor = 0
	elev.dir = 0
	elev.destination = 0
	elev.fsmState = IDLE
	return elev
}

func (elev ElevatorState) GetElevatorDirection() int{
	return elev.dir
}


func (elev *ElevatorState) IDLE() { // we have to use pointer reciever beacause we want to read and write as oposed to just read. 	
	fmt.Println("Going to state: IDLE")
	driver.StopElevate()
	elev.fsmState = IDLE
	elev.dir = 0
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
	fmt.Println("Going in to state: DOOR_OPEN")
	driver.StopElevate()
	elev.fsmState = DOOR_OPEN
	driver.ElevSetDoorOpenLamp(1)
	time.Sleep(4000*time.Millisecond) //The lines under should happend after sleep. 
	driver.ElevSetDoorOpenLamp(0)
}



