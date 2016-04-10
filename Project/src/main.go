package main

import (
	"./driver"
	"./eventhandler"
	"./fsm"
	"./queue"
	"fmt"
	//"time"
)

func main() {

	//Initzialization of elevator Hardware and fsm.
	if driver.Init() == 0 {
		fmt.Println("The elevator was not able to initialize")
	}
	if driver.Init() == 1 {
		fmt.Println("The elevator was able to initialize")
	}

	var queue queue.Order
	var elevator fsm.ElevatorState

	elevator.InitFsm()

	floorChannel := make(chan int)
	buttonChannel := make(chan eventhandler.Button_info)

	//Starting gorutines to check for events on buttons and floor sensors
	eventhandler.CheckEvents(floorChannel, buttonChannel)
	PrevDirection := -1
	// infinite loop, to keep the elevator going
	for {
		select {

		case NewEvent := <-floorChannel: // Gets 0,1,2 or 3, never -1 
			dir := elevator.GetDirection()
			if NewEvent != -1 {
				elevator.Setfloor(NewEvent)
				if queue.ShouldStop(NewEvent, dir) {
					queue.RemoveOrder(NewEvent, PrevDirection)
					elevator.DoorOpen()
				}
			} else {
				PrevDirection = dir
			}

		case NewEvent := <-buttonChannel:
			queue.AddOrder(NewEvent.Floor, NewEvent.Button)
			/*if elevator.GetDirection() != queue.QueueDirection(PrevDirection, elevator.GetFloor()) {
				elevator.SetDirection(queue.QueueDirection(PrevDirection, elevator.GetFloor()))
				if elevator.GetDirection() != 0 {
					elevator.Elevating(elevator.GetDirection())
				}
			}*/

		default:
			switch elevator.GetState() {
			case fsm.IDLE:
				fmt.Println("Inside deafulte")
				fmt.Println(PrevDirection)
				direct := queue.QueueDirection(PrevDirection, elevator.GetFloor())
				fmt.Println("The direction from que is set to:", direct)
				elevator.Elevating(direct)
			}

		}

	}
}
