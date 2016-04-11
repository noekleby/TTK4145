package main

import (
	"./driver"
	"./eventhandler"
	"./fsm"
	"./queue"
	"fmt"
	"./network"
	//"./transManager.go"
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

	//transManager.Init()
	elevator.InitFsm()
	floorChan := make(chan int)
	UpOrderChan := make(chan int)
	DownOrderChan := make(chan int)
	CommandOrderChan := make(chan int)


	//putt dette i main: 
	newElevatorChan := make(chan string)
	deadElevatorChan := make(chan string)
	
	go network.SendHeartbeat()
	go network.HeartbeatEventCheck(newElevatorChan, deadElevatorChan)
	go eventhandler.HeartbeatEventHandler(newElevatorChan, deadElevatorChan)

	//Starting gorutines to check for events on buttons and floor sensors
	eventhandler.CheckEvents(UpOrderChan, DownOrderChan, CommandOrderChan, floorChan)
	PrevDirection := -1

	// infinite loop, to keep the elevator going
	for {
		select {

		case floor := <-floorChan:
			fmt.Println("I do get a floor signal.")
			dir := elevator.GetDirection()
			if floor != -1 {
				fmt.Println("I do go inside the floor if sentence")
				elevator.Setfloor(floor)
				if queue.ShouldStop(floor, dir) {
					fmt.Println("Should I stop? yes.")
					queue.RemoveOrder(floor, PrevDirection)
					elevator.DoorOpen()
				}
			} else {
				PrevDirection = dir
			}

		case floor := <-DownOrderChan:
			queue.AddOrder(floor, 1)
			if elevator.GetDirection() != queue.QueueDirection(PrevDirection, elevator.GetFloor()) {
				elevator.SetDirection(queue.QueueDirection(PrevDirection, elevator.GetFloor()))
				if elevator.GetDirection() != 0 {
					elevator.Elevating(elevator.GetDirection())
				}
			}
		case floor := <-UpOrderChan:
			queue.AddOrder(floor, 0)
			if elevator.GetDirection() != queue.QueueDirection(PrevDirection, elevator.GetFloor()) {
				elevator.SetDirection(queue.QueueDirection(PrevDirection, elevator.GetFloor()))
				if elevator.GetDirection() != 0 {
					elevator.Elevating(elevator.GetDirection())
				}
			}
		case floor := <-CommandOrderChan:
			queue.AddOrder(floor, 2)
			if elevator.GetDirection() != queue.QueueDirection(PrevDirection, elevator.GetFloor()) {
				elevator.SetDirection(queue.QueueDirection(PrevDirection, elevator.GetFloor()))
				if elevator.GetDirection() != 0 {
					elevator.Elevating(elevator.GetDirection())
				}
			}
		default:
			switch elevator.GetState() {
			case fsm.IDLE:
				//fmt.Println("Inside default")
				//fmt.Println(PrevDirection)
				direction := queue.QueueDirection(PrevDirection, elevator.GetFloor())
				if direction == 0 && queue.EmptyQueue() {
					elevator.IDLE()
				} else if direction == 0 && !queue.EmptyQueue() {
					elevator.DoorOpen()
					queue.RemoveOrder(elevator.GetFloor(), 0)
				} else {
					elevator.Elevating(direction)
				}

			}

		}

	}
}
