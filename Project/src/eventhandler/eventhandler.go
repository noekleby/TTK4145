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
			Elevators[msg.SenderIP].ExternalUp[msg.Elevator.Floor] = false
			lightEventChan <- 1
			Elevators[msg.SenderIP].Floor = msg.Elevator.Floor
			Elevators[msg.SenderIP].Direction = 1

		case "Remove order down":
			Elevators[msg.SenderIP].ExternalDown[msg.Elevator.Floor] = false
			lightEventChan <- 1
			Elevators[msg.SenderIP].Floor = msg.Elevator.Floor
			Elevators[msg.SenderIP].Direction = -1

		case "Add order":
			if msg.Order.Buttontype == UP {
				Elevators[msg.TargetIP].ExternalUp[msg.Order.Floor] = true
				lightEventChan <- 1
				Elevators[msg.SenderIP].Floor = msg.Elevator.Floor

			} else if msg.Order.Buttontype == DOWN {
				Elevators[msg.TargetIP].ExternalDown[msg.Order.Floor] = true
				lightEventChan <- 1
				Elevators[msg.SenderIP].Floor = msg.Elevator.Floor

			} else {
				Elevators[msg.TargetIP].InternalOrders[msg.Order.Floor] = true
				Elevators[msg.SenderIP].Floor = msg.Elevator.Floor
			}
		case "Status Update":
			if msg.TargetIP == GetLocalIP(){
				for floor := 0; floor < N_FLOORS; floor++{
					Elevators[msg.TargetIP].NewlyInit = false
					Elevators[msg.TargetIP].InternalOrders[floor] = msg.Elevator.InternalOrders[floor]
					Elevators[msg.TargetIP].ExternalUp[floor] = msg.Elevator.ExternalUp[floor]
					Elevators[msg.TargetIP].ExternalDown[floor] = msg.Elevator.ExternalDown[floor]
				}
			} else if msg.TargetIP == "" && msg.Elevator.NewlyInit == true && msg.SenderIP != GetLocalIP() {
				newMsg := Message{"Status update", GetLocalIP(), msg.SenderIP, *(Elevators[msg.SenderIP]), Order{-1,-1,""}}
				BroadcastMessage(newMsg)
			} 
		}
	}
}

func HeartbeatEventHandler(newElevatorChan chan string, deadElevatorChan chan string) {
	for {
		select {
		case ip := <- newElevatorChan:
			fmt.Println("A new Elevator online with IP:", ip)
			_, exist := Elevators[ip]
			if exist {
				Elevators[ip].Active = true
			} else {
				fmt.Println("Meeting new elevator")
				Elevators[ip] = &Elevator{true, -1, -1, IDLE, false, [N_FLOORS]bool{false, false, false, false}, [N_FLOORS]bool{false, false, false, false}, [N_FLOORS]bool{false, false, false, false}}
			}
		case ip := <- deadElevatorChan:
			fmt.Println("We have lost an elevator with IP:", ip)
			Elevators[ip].Active = false
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
			if floor != -1 {
				driver.SetFloorIndicator(floor)
				Elevators[GetLocalIP()].Floor = floor
				if queue.ShouldStop(floor, dir) {
					Elevators[GetLocalIP()].NewlyInit = false
					queue.RemoveLocalOrder(floor, PrevDirection, lightEventChan)
					fsm.GoToDoorOpen()
				}
			} else {
				PrevDirection = dir
			}

		case order := <-buttonChan:
			Elevators[GetLocalIP()].NewlyInit = false
			queue.AddLocalOrder(order, lightEventChan)
			if (Elevators[GetLocalIP()].Direction) != queue.NextDirection(PrevDirection, Elevators[GetLocalIP()].Floor) {
				Elevators[GetLocalIP()].Direction = queue.NextDirection(PrevDirection, Elevators[GetLocalIP()].Floor)
				if Elevators[GetLocalIP()].Direction != 0 {
					fsm.GoToElevating(Elevators[GetLocalIP()].Direction)
				}
			}

		default:
			switch Elevators[GetLocalIP()].FsmState {
			case IDLE:
				direction := queue.NextDirection(PrevDirection, Elevators[GetLocalIP()].Floor)
				if direction == 0 && queue.EmptyQueue() {
					fsm.GoToIDLE()
				} else if direction == 0 && !queue.EmptyQueue() {
					fsm.GoToDoorOpen()
					queue.RemoveLocalOrder(Elevators[GetLocalIP()].Floor, PrevDirection, lightEventChan)
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

	for {
		var order Order
		for floor := 0; floor < N_FLOORS; floor++ {
			for buttonType := 0; buttonType < N_BUTTONS; buttonType++ {
				if driver.GetButtonSignal(buttonType, floor) == 1 {
					switch buttonType {
					case UP:
						if !Elevators[GetLocalIP()].ExternalUp[floor]{
							order.Buttontype = buttonType
							order.Floor = floor
							order.FromIP = ""
							buttonChan <- order
						}
					case DOWN:
						if !Elevators[GetLocalIP()].ExternalDown[floor]{
							order.Buttontype = buttonType
							order.Floor = floor
							order.FromIP = ""
							buttonChan <- order
						}
					case COMMAND:
						if !Elevators[GetLocalIP()].InternalOrders[floor]{
							order.Buttontype = buttonType
							order.Floor = floor
							order.FromIP = ""
							buttonChan <- order
						}
					}
				} 
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
}


func LampHandler(lightEventChan chan int) {
	InternalLamp := make([]bool, N_FLOORS)
	ExternalUpLamp := make([]bool, N_FLOORS)
	ExternalDownLamp := make([]bool, N_FLOORS)

	for {
		<-lightEventChan
		for floor := 0; floor < N_FLOORS; floor++ {
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
		for floor := 0; floor < N_FLOORS; floor++ {
			if floor > 0 && ExternalDownLamp[floor] {
				driver.SetButtonLamp(floor, DOWN, true)
			} else {
				driver.SetButtonLamp(floor, DOWN, false)
			}
			if floor < N_FLOORS-1 && ExternalUpLamp[floor] {
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
