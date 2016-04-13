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
		// We do not want duplicate code.
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

/*func CheckEvents(UpOrderChan chan int, DownOrderChan chan int, CommandOrderChan chan int, floorSensorChan chan int, betweenFloorSensorChan chan int) {
	elevator.InitFsm()
	go floorEventCheck(floorSensorChan, betweenFloorSensorChan)
	go buttonEventCheck(UpOrderChan, DownOrderChan, CommandOrderChan)
	go buttonEventHandler(UpOrderChan, DownOrderChan, CommandOrderChan)
	go floorEventHandler(floorSensorChannel, betweenFloorSensorChan)
}

func floorEventHandler(floorSensorChan chan int, betweenFloorSensorChan chan int) {

	for {
		select {
		case floor <- betweenFloorSensorChan:
			//give a message to queue: "Not in floor"
			fmt.Println("Between floors, should make a function to take car of this situation.")

		case floor <- floorSensorChan:
			if queue.ShouldStop(floor) {
				elevator.GoToDoorOpen()
				queue.CompletedOrder(floor, "local")
			}
		}
	}

}


func buttonEventHandler(UpOrderChan chan int, DownOrderChan chan int, CommandOrderChan chan int) {

	for {
		select {
		case floor := <-UpOrderChan:
			newOrder := Order{UP, floor}
			i := AddLocalOrder(newOrder)
			switch i {
			case "empty":
				if elevator.GetState() == IDLE {
					direction := queue.NextDirection()
					elevator.GoToElevating(direction)
				} else {
					fmt.Println("Elevator should always be in IDLE when empty queue, but it is not")
				}
			case "sameFloor":
				switch elevator.GetState() {
				case IDLE:
					elevator.GoToDoorOpen()
					queue.OrderCompleted()

				case DOOR_OPEN:
					elevator.GoToDoorOpen()
					queue.OrderCompleted()
				}

			}

		case floor := <-DownOrderChan:
			newOrder := Order{DOWN, floor}
			i := AddLocalOrder(newOrder)
			switch i {
			case "empty":
			case "sameFloor":
			}

		case floor := <-CommandOrderChan:
			newOrder := Order{COMMAND, floor}
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
			} else {
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

*/
/*func ButtonEventCheck(event chan Button_info) {
	prevButton := button
	newButton := button
	fmt.Println(newButton)
	for {
		fmt.Println("Inside forever")
		for f := 0; f < driver.N_FLOORS; f++ {
			for b := 0; b < driver.N_BUTTONS; b++ {
				newButton[f][b] = driver.ElevGetButtonSignal(b, f)
				fmt.Println(newButton, prevButton)
				if (newButton[f][b] != prevButton[f][b]) && (newButton[f][b] == 1) {
					prevButton[f][b] = newButton[f][b]
					var buttonInf Button_info
					buttonInf.Button = b
					buttonInf.Floor = f
					event <- buttonInf
				}
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
}*/
