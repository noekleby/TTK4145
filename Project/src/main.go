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
		case NewEvent := <-floorChannel:
			Queueinfo = queue.GetQueueInfo()
			fmt.Println("New floor event happened.", NewEvent)
			if NewEvent != -1 {
				fmt.Println(Queueinfo.GetIntQ(NewEvent), "Inside if sentence")
				fmt.Println(Queueinfo.GetIntQ(0), "IntQ0")
				fmt.Println(Queueinfo.GetIntQ(1), "IntQ1")
				fmt.Println(Queueinfo.GetIntQ(2), "IntQ2")
				fmt.Println(Queueinfo.GetIntQ(3), "IntQ3")
				fmt.Println(Queueinfo.GetIntQ(NewEvent+1), "IntQ3")
				if (Queueinfo.GetIntQ(NewEvent) == 1 || Queueinfo.GetExUp(NewEvent) == 1) && (Elevatorinfo.GetElevatorDirection() == 1) {
					Elevatorinfo.DoorOpen()
				}
				if (Queueinfo.GetIntQ(NewEvent) == 1 || Queueinfo.GetExDown(NewEvent) == 1) && (Elevatorinfo.GetElevatorDirection() == -1) {
					Elevatorinfo.DoorOpen()
				}
			}

		case NewEvent := <-buttonChannel:
			driver.SetButtonLamp(NewEvent.Floor, NewEvent.Button, true)
			Queueinfo.Add(NewEvent.Floor, NewEvent.Button)
			queueDir := Queueinfo.QueueDirection()
			Dir := Elevatorinfo.GetElevatorDirection()
			fmt.Println("queue direction", queueDir)
			if Dir == queueDir {
				fmt.Println("Order in same direction")

			} else {
				Elevatorinfo.Elevating(queueDir)

			}

		}
	}
}
