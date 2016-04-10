package eventhandler

import (
	"../driver"
	"fmt"
	"time"
)

func CheckEvents(UpOrderChan chan int, DownOrderChan chan int, CommandOrderChan chan int, floorSensorChan) {
	go FloorEventCheck(floorSensorChannel)
	go ButtonEventCheck(UpOrderChan, DownOrderChan, CommandOrderChan)
	go ButtonEventHandler(UpOrderChan, DownOrderChan, CommandOrderChan)
}


func ButtonEventHandler(UpOrderChan chan int, DownOrderChan chan int, CommandOrderChan chan int){

	for{
		select{
		case floor := <- UpOrderChan:
			
		case floor := <- DownOrderChan:

		case floor := <- CommandOrderChan:

		}

	}
}



func FloorEventCheck(floorSignal chan int) {
	prevFloor := 0
	for {
		newFloor := driver.GetFloorSignal()
		if newFloor != prevFloor {
			prevFloor = newFloor
			floorSignal <- newFloor
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func ButtonEventCheck(UpOrderChan chan int, DownOrderChan chan int, CommandOrderChan chan int){
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




