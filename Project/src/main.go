package main

import (
	"./driver"
	"fmt"
	"./eventhandler"
	"./queue"
	"./fsm"
	//"time"
)



func main() {
	if driver.Init() == 0 {
		fmt.Println("The elevator was not able to initialize")
	}
	if driver.Init() == 1 {
		fmt.Println("The elevator was able to initialize")
	}
	Elevatorinfo := fsm.InitFsm()
	Queueinfo := queue.GetQueueInfo()
	floorChannel := make(chan int)
	buttonChannel := make(chan eventhandler.Button_info)

	eventhandler.CheckEvents(floorChannel, buttonChannel)

	for {
		select {
			case NewEvent := <- floorChannel:
				fmt.Println("New floor event happend.", NewEvent)
				if NewEvent != -1{
					if (Queueinfo.InternalOrders[NewEvent] == 1 || Queueinfo.ExternalUp[NewEvent] == 1) && (Elevatorinfo.GetElevatorDirection() == 1){
						Elevatorinfo.DoorOpen()
					}
					if (Queueinfo.InternalOrders[NewEvent] == 1 || Queueinfo.ExternalDown[NewEvent] == 1) && (Elevatorinfo.GetElevatorDirection() == -1) {
						Elevatorinfo.DoorOpen()
					}
				}


			case NewEvent := <- buttonChannel:
				driver.SetButtonLamp(NewEvent.Floor, NewEvent.Button, true)
				queue.AddOrder(NewEvent.Floor, NewEvent.Button)
				queueDir := Queueinfo.QueueDirection()
				Dir := Elevatorinfo.GetElevatorDirection()
				fmt.Println("queue direction", queueDir)
				if Dir == queueDir{
					fmt.Println("Order in same direction")

				} else {
					Elevatorinfo.Elevating(queueDir)

				}

		}
	}
}

