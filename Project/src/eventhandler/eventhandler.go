package eventhandler

import (
	. "../definitions"
	"../driver"
	"../fsm"
	. "../network"
	"../queue"
	"fmt"
	"time"
)

func MessageTypeHandler(messageReciveChan chan Message, floorChan chan int, upOrderChan chan int, downOrderChan chan int, commandOrderChan chan int) {

	for {
		msg := <-messageReciveChan
		switch msg.MessageType {

		case "Remove order up":
			if msg.SenderIP != GetLocalIP() {
				updateElevatorStatus(msg.SenderIP, msg.Elevator)
				queue.RemoveRemoteOrder(msg.Elevator.Floor, UP)
			}

		case "Remove order down":
			if msg.SenderIP != GetLocalIP() {
				updateElevatorStatus(msg.SenderIP, msg.Elevator)
				queue.RemoveRemoteOrder(msg.Elevator.Floor, DOWN)
			}

		case "Add order up":
			fmt.Println("Hi!")
			if msg.SenderIP != msg.TargetIP {
				queue.AddRemoteOrder(msg.TargetIP, msg.Elevator.ExternalUp, UP)
			}
			if msg.SenderIP != GetLocalIP() {
				updateElevatorStatus(msg.SenderIP, msg.Elevator)
			}

		case "Add order down":
			if msg.SenderIP != msg.TargetIP {
				queue.AddRemoteOrder(msg.TargetIP, msg.Elevator.ExternalDown, DOWN)
			}
			if msg.SenderIP != GetLocalIP() {
				updateElevatorStatus(msg.SenderIP, msg.Elevator)
			}

		case "Add internal order":
			if msg.SenderIP != GetLocalIP() {
				updateElevatorStatus(msg.SenderIP, msg.Elevator)
			}
		}
	}
}

func updateElevatorStatus(MessageIP string, elevator Elevator) {
	Elevators[MessageIP].Active = elevator.Active
	Elevators[MessageIP].Floor = elevator.Floor
	Elevators[MessageIP].Direction = elevator.Direction
	Elevators[MessageIP].PrevFloor = elevator.PrevFloor
	Elevators[MessageIP].FsmState = elevator.FsmState

	for i := 0; i < driver.N_FLOORS; i++ {
		Elevators[MessageIP].InternalOrders[i] = elevator.InternalOrders[i]
		Elevators[MessageIP].ExternalUp[i] = elevator.ExternalUp[i]
		Elevators[MessageIP].ExternalDown[i] = elevator.ExternalDown[i]
	}
}

func HeartbeatEventHandler(newElevatorChan chan string, deadElevatorChan chan string) {
	for {
		select {
		case IP := <-newElevatorChan:
			fmt.Println("A new Elevator online:", IP)
			_, exist := Elevators[IP]
			if exist {
				Elevators[IP].Active = true
			} else {
				fmt.Println("Meeting new elevator")
				Elevators[IP] = &Elevator{true, -1, 0, -1, IDLE, [driver.N_FLOORS]bool{false, false, false, false}, [driver.N_FLOORS]bool{false, false, false, false}, [driver.N_FLOORS]bool{false, false, false, false}}
			}
		case IP := <-deadElevatorChan:
			fmt.Println("We have lost an elevator:", IP)
			Elevators[IP].Active = false
		}
	}
}

func ButtonandFloorEventHandler(floorChan chan int, upOrderChan chan int, downOrderChan chan int, commandOrderChan chan int) {

	go FloorEventCheck(floorChan)
	go ButtonEventCheck(upOrderChan, downOrderChan, commandOrderChan)
	PrevDirection := -1

	for {
		select {

		case floor := <-floorChan:
			dir := Elevators[GetLocalIP()].Direction
			if floor != -1 {
				Elevators[GetLocalIP()].Floor = floor
				if queue.ShouldStop(floor, dir) {
					queue.RemoveOrder(floor, PrevDirection)
					fsm.GoToDoorOpen()
				}
			} else {
				PrevDirection = dir
			}

		case floor := <-downOrderChan:
			queue.AddLocalOrder(floor, DOWN)
			if (Elevators[GetLocalIP()].Direction) != queue.QueueDirection(PrevDirection, Elevators[GetLocalIP()].Floor) {
				Elevators[GetLocalIP()].Direction = queue.QueueDirection(PrevDirection, Elevators[GetLocalIP()].Floor)
				if Elevators[GetLocalIP()].Direction != 0 {
					fsm.GoToElevating(Elevators[GetLocalIP()].Direction)
				}
			}
		// All of the cases below does the same thing with different floor, make a new routine to take care of this.
		// We do not want duplicate code. Function made, you can find it at the bottom of the page. 
		case floor := <-upOrderChan:
			queue.AddLocalOrder(floor, UP)
			if Elevators[GetLocalIP()].Direction != queue.QueueDirection(PrevDirection, Elevators[GetLocalIP()].Floor) {
				Elevators[GetLocalIP()].Direction = queue.QueueDirection(PrevDirection, Elevators[GetLocalIP()].Floor)
				if Elevators[GetLocalIP()].Direction != 0 {
					fsm.GoToElevating(Elevators[GetLocalIP()].Direction)
				}
			} else {
				fmt.Println("Still in the same floor.")
			}
		case floor := <-commandOrderChan:
			queue.AddLocalOrder(floor, COMMAND)
			if Elevators[GetLocalIP()].Direction != queue.QueueDirection(PrevDirection, Elevators[GetLocalIP()].Floor) {
				Elevators[GetLocalIP()].Direction = queue.QueueDirection(PrevDirection, Elevators[GetLocalIP()].Floor)
				if Elevators[GetLocalIP()].Direction != 0 {
					fsm.GoToElevating(Elevators[GetLocalIP()].Direction)
				}
			}
		default:
			switch Elevators[GetLocalIP()].FsmState {
			case IDLE:
				direction := queue.QueueDirection(PrevDirection, Elevators[GetLocalIP()].Floor)
				if direction == 0 && queue.EmptyQueue() {
					fsm.GoToIDLE()
				} else if direction == 0 && !queue.EmptyQueue() {
					fsm.GoToDoorOpen()
					queue.RemoveOrder(Elevators[GetLocalIP()].Floor, 0)
				} else {
					fsm.GoToElevating(direction)
				}

			}

		}

	}
}

func FloorEventCheck(event chan int) {
	prevFloor := 0
	for {
		newFloor := driver.GetFloorSignal()
		if newFloor != prevFloor {
			prevFloor = newFloor
			event <- newFloor
		}
		time.Sleep(200 * time.Millisecond)
	}
}


func ButtonEventCheck(UpOrderChan chan int, DownOrderChan chan int, CommandOrderChan chan int) {

	//We might not need this, we could change it to check the elevator[IP].queue instead. 
	buttonPressed := make([][]bool, driver.N_FLOORS) //Makes row
	for i := range buttonPressed {
		buttonPressed[i] = make([]bool, driver.N_BUTTONS) //Makes column
	}

	for {
		for floor := 0; floor < driver.N_FLOORS; floor++ {

			if driver.ElevGetButtonSignal(0, floor) == 1 && !buttonPressed[floor][0] {
				buttonPressed[floor][0] = true
				UpOrderChan <- floor

			} else if driver.ElevGetButtonSignal(1, floor) == 1 && !buttonPressed[floor][1] {
				buttonPressed[floor][1] = true
				DownOrderChan <- floor

			} else if driver.ElevGetButtonSignal(2, floor) == 1 && !buttonPressed[floor][2] {
				buttonPressed[floor][2] = true
				CommandOrderChan <- floor

			} else if buttonPressed[floor][0] {
				buttonPressed[floor][0] = false

			} else if buttonPressed[floor][1] {
				buttonPressed[floor][1] = false

			} else if buttonPressed[floor][2] {
				buttonPressed[floor][2] = false
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func ElevateTowardNextDirection() {
	if Elevators[GetLocalIP()].Direction != queue.QueueDirection(PrevDirection, Elevators[GetLocalIP()].Floor) {
		Elevators[GetLocalIP()].Direction = queue.QueueDirection(PrevDirection, Elevators[GetLocalIP()].Floor)
		if Elevators[GetLocalIP()].Direction != 0 {
			fsm.GoToElevating(Elevators[GetLocalIP()].Direction)
		}
	}
}
