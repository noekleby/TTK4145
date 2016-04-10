package eventhandler

import (
	"../driver"
	"../definitions"
	"fmt"
	"time"
	"fsm"
)

func CheckEvents(UpOrderChan chan int, DownOrderChan chan int, CommandOrderChan chan int, floorSensorChan chan int, betweenFloorSensorChan chan int) {
	elevator.InitFsm()
	go floorEventCheck(floorSensorChan, betweenFloorSensorChan)
	go buttonEventCheck(UpOrderChan, DownOrderChan, CommandOrderChan)
	go buttonEventHandler(UpOrderChan, DownOrderChan, CommandOrderChan)
	go floorEventHandler(floorSensorChannel, betweenFloorSensorChan)
}

func floorEventHandler(floorSensorChan chan int, betweenFloorSensorChan chan int) {

	for{
		select{
			case floor <- betweenFloorSensorChan:
				//give a message to queue: "Not in floor"
				fmt.Println("Between floors, should make a function to take car of this situation.")

			case floor <- floorSensorChan:
				if queue.ShouldStop(floor){
					elevator.GoToDoorOpen()
					queue.CompletedOrder(floor, "local")
				}
		}
	}
	
}

func buttonEventHandler(UpOrderChan chan int, DownOrderChan chan int, CommandOrderChan chan int){

	for{
		select{
		case floor := <- UpOrderChan:
			newOrder := Order{ UP, floor}
			i := AddLocalOrder(newOrder)
			switch i {
				case "empty": 
					if elevator.GetState() == IDLE{
						direction := queue.NextDirection()
						elevator.GoToElevating(direction)
					} else {
						fmt.Println("Elevator should always be in IDLE when empty queue, but it is not")
					}
				case "sameFloor":
					switch elevator.GetState(){
					case IDLE:
						elevator.GoToDoorOpen()
						queue.OrderCompleted()

					case DOOR_OPEN: 
						elevator.GoToDoorOpen()
						queue.OrderCompleted()
					}

			}
			
		case floor := <- DownOrderChan:
			newOrder := Order{ DOWN, floor}
			i := AddLocalOrder(newOrder)
			switch i {
				case "empty": 
				case "sameFloor":
			}

		case floor := <- CommandOrderChan:
			newOrder := Order{ COMMAND, floor}
			i := AddLocalOrder(newOrder)
			switch i {
				case "empty": 
				case "sameFloor":
			}

		}

	}
}



func floorEventCheck(floorSensorChan chan int, betweenFloorSensorChan chan int) {
	prevFloor := 0
	for {
		newFloor := driver.GetFloorSignal()
		if newFloor != prevFloor {
			prevFloor = newFloor
			if newFloor == -1 {
				betweenFloorSensorChan <- newFloor
			} else{
				floorSensorChan <- newFloor
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func buttonEventCheck(UpOrderChan chan int, DownOrderChan chan int, CommandOrderChan chan int){
	buttons := [driver.N_FLOORS][driver.N_BUTTONS]int{ 
		{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
		{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
		{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
		{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
	}

	buttonPressed := make([][]bool, driver.N_BUTTONS) //Makes row
	for i := range buttonPressed{
		buttonPressed[i] = make([]bool, driver.N_FLOORS) //Makes column 
	}

	for {
		for floor := 0; floor < driver.N_FLOORS; floor++ {

			if ElevGetButtonSignal(0, floor) == 1 && !buttonPressed[floor][0] {
				buttonPressed[floor][0] = true
				UpOrderChan <- floor 

			} else if ElevGetButtonSignal(1, floor) == 1 && !buttonPressed[floor][1] {
				buttonPressed[floor][1] = true
				DownOrderChan <- floor 

			}else if ElevGetButtonSignal(2, floor) == 1 && !buttonPressed[floor][2] {
				buttonPressed[floor][2] = true
				CommandOrderChan <- floor 

			}else if buttonPressed([floor][0]){
				buttonPressed[floor][0] = false

			}else if buttonPressed([floor][1]){
				buttonPressed[floor][1] = false

			}else if buttonPressed([floor][2]){
				buttonPressed[floor][2] = false	
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
}




