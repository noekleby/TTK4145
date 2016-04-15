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

func MessageTypeHandler(messageReciveChan chan Message, floorChan chan int, buttonChan chan Order, lightEventChan chan int) {

	for {
		msg := <-messageReciveChan
		switch msg.MessageType {

		case "Remove order up":
			fmt.Println("IP in Remove order up: ", msg.SenderIP, msg.TargetIP)
			Elevators[msg.SenderIP].ExternalUp[msg.Elevator.Floor] = false
			lightEventChan <- 1
			Elevators[msg.SenderIP].Floor = msg.Elevator.Floor
			Elevators[msg.SenderIP].Direction = 1

		case "Remove order down":
			fmt.Println("IP in Remove order down: ", msg.SenderIP, msg.TargetIP)
			Elevators[msg.SenderIP].ExternalDown[msg.Elevator.Floor] = false
			lightEventChan <- 1
			Elevators[msg.SenderIP].Floor = msg.Elevator.Floor
			Elevators[msg.SenderIP].Direction = -1

		case "Add order":
			fmt.Println("IP in Add: ", msg.SenderIP, msg.TargetIP)
			if msg.Order.Buttontype == UP {
				Elevators[msg.TargetIP].ExternalUp[msg.Order.Floor] = true
				lightEventChan <- 1
			} else if msg.Order.Buttontype == DOWN {
				Elevators[msg.TargetIP].ExternalDown[msg.Order.Floor] = true
				lightEventChan <- 1
			} else {
				Elevators[msg.TargetIP].InternalOrders[msg.Order.Floor] = true
			}
		}
	}
}

func updateElevatorStatus(MessageIP string, elevator Elevator) {
	Elevators[MessageIP].Active = elevator.Active
	Elevators[MessageIP].Floor = elevator.Floor
	Elevators[MessageIP].Direction = elevator.Direction
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
				Elevators[IP] = &Elevator{true, -1, -1, IDLE, [driver.N_FLOORS]bool{false, false, false, false}, [driver.N_FLOORS]bool{false, false, false, false}, [driver.N_FLOORS]bool{false, false, false, false}}
			}
		case IP := <-deadElevatorChan:
			fmt.Println("We have lost an elevator:", IP)
			Elevators[IP].Active = false
		}
	}
}

func ButtonandFloorEventHandler(floorChan chan int, buttonChan chan Order, lightEventChan chan int) {

	go FloorEventCheck(floorChan)
	go ButtonEventCheck(buttonChan)
	PrevDirection := -1

	for {
		select {

		case floor := <-floorChan:
			dir := Elevators[GetLocalIP()].Direction
			fmt.Println(Elevators[GetLocalIP()].Direction)
			fmt.Println(PrevDirection)
			if floor != -1 {
				driver.SetFloorIndicator(floor)
				Elevators[GetLocalIP()].Floor = floor
				if queue.ShouldStop(floor, dir) {
					queue.RemoveOrder(floor, PrevDirection, lightEventChan)
					fsm.GoToDoorOpen()
				}
			} else {
				PrevDirection = dir
			}

		case order := <-buttonChan:
			queue.AddLocalOrder(order, lightEventChan)
			if (Elevators[GetLocalIP()].Direction) != queue.QueueDirection(PrevDirection, Elevators[GetLocalIP()].Floor) {
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
					queue.RemoveOrder(Elevators[GetLocalIP()].Floor, 0, lightEventChan)
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

func ButtonEventCheck(buttonChan chan Order) {

	//We might not need this, we could change it to check the elevator[IP].queue instead.
	buttonPressed := make([][]bool, driver.N_FLOORS) //Makes row
	for i := range buttonPressed {
		buttonPressed[i] = make([]bool, driver.N_BUTTONS) //Makes column
	}

	for {
		var order Order
		for floor := 0; floor < driver.N_FLOORS; floor++ {
			for buttonType := 0; buttonType < driver.N_BUTTONS; buttonType++ {
				if driver.ElevGetButtonSignal(buttonType, floor) == 1 && !buttonPressed[floor][buttonType] {
					buttonPressed[floor][buttonType] = true
					order.Buttontype = buttonType
					order.Floor = floor
					order.FromIP = ""
					buttonChan <- order
				} else if buttonPressed[floor][buttonType] {
					buttonPressed[floor][buttonType] = false
				}
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func LampHandler(lightEventChan chan int) {
	InternalLamp := make([]bool, driver.N_FLOORS)
	ExternalUpLamp := make([]bool, driver.N_FLOORS)
	ExternalDownLamp := make([]bool, driver.N_FLOORS)

	for {
		<-lightEventChan
		for floor := 0; floor < driver.N_FLOORS; floor++ {
			for IP, _ := range Elevators {
				if Elevators[GetLocalIP()].InternalOrders[floor] {
					InternalLamp[floor] = true
				}
				if Elevators[IP].ExternalUp[floor] {
					ExternalUpLamp[floor] = true
				}
				if Elevators[IP].ExternalDown[floor] {
					ExternalDownLamp[floor] = true
				}
			}
		}
		for floor := 0; floor < driver.N_FLOORS; floor++ {
			if floor > 0 && ExternalDownLamp[floor] {
				driver.SetButtonLamp(floor, DOWN, true)
			} else {
				driver.SetButtonLamp(floor, DOWN, false)
			}
			if floor < driver.N_FLOORS-1 && ExternalUpLamp[floor] {
				driver.SetButtonLamp(floor, UP, true)
			} else {
				driver.SetButtonLamp(floor, UP, false)
			}
			if InternalLamp[floor] {
				driver.SetButtonLamp(floor, COMMAND, true)
			} else {
				driver.SetButtonLamp(floor, COMMAND, false)
			}
			InternalLamp[floor] = false
			ExternalUpLamp[floor] = false
			ExternalDownLamp[floor] = false

		}
	}
}
