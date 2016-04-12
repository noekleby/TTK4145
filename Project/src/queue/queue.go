package queue

import (
	"../driver"
	."../definitions"
	"fmt"
	//"../eventhandler"
	."../network"
)

const(
	N_FLOORS = 4
)

func ShouldStop(floor, dir int) bool {
	if Elevators[GetLocalIP()].InternalOrders[floor] == true {
		return true
	}
	if dir == 1 {
		if Elevators[GetLocalIP()].ExternalUp[floor] == true || floor == driver.N_FLOORS-1 {
			return true
		} else if QueueDirectionUp(floor){
			return false
		} else {
			return true
		}
	} else if dir == -1 {
		if Elevators[GetLocalIP()].ExternalDown[floor] == true || floor == 0 {
			return true
		} else if QueueDirectionDown(floor) {
			return false
		} else {
			return true
		}
	}
	return true
}


/*
func QueueDirectionDown(floor int) int {
	//Already checked if there are any Elevators in queue and there are
	check := false
	if floor == 0 {
		return 1
	} else {
		for i := floor; i >= 0; i-- {
			if Elevators[GetLocalIP()].InternalOrders[i]  == true || Elevators[GetLocalIP()].ExternalDown[i] == true {
				check = true
			}
		}
		if check == true {
			return -1
		} else {
			return 1
		}
	}
}

func QueueDirectionUp(floor int) int {
	check := false
	if floor == 3 {
		return -1
	} else {
		for i := floor; i < driver.N_FLOORS; i++ {
			if Elevators[GetLocalIP()].InternalOrders[i]  == true || Elevators[GetLocalIP()].ExternalUp[i] == true {
				check = true
			}
		}
		if check == true {
			return 1
		} else {
			return -1
		}
	}
}
/*

func (Elevator *Elevator) QueueDirection(direction, floor int) int {
	if Elevator.EmptyQueue() == true {
		return 0

	} else if direction == 1 {
		return (Elevator.QueueDirectionUp(floor))

	} else if direction == -1 {
		return (Elevator.QueueDirectionDown(floor))
	} else {
		fmt.Println("Something wrong with queue.")
		return 0
	}
}*/

func QueueDirection(direction, floor int) int {
	if EmptyQueue() == true {
		return 0

	} else if direction == 1 {
		if QueueDirectionUp(floor) {
			return 1
		} else if QueueDirectionDown(floor) {
			return -1
		}
	} else if direction == -1 {
		if QueueDirectionDown(floor)  {
			return -1
		} else if QueueDirectionUp(floor) {
			return 1
		}
	}
	return 0
}

func QueueDirectionUp(floor int) bool {
	for f := floor + 1; f < driver.N_FLOORS; f++ {
		if Elevators[GetLocalIP()].InternalOrders[f] == true || Elevators[GetLocalIP()].ExternalUp[f] == true || Elevators[GetLocalIP()].ExternalDown[f] == true {
			return true
		}
	}
	return false
}

func QueueDirectionDown(floor int) bool {
	for f := floor - 1; f > -1; f-- {
		if Elevators[GetLocalIP()].InternalOrders[f] == true || Elevators[GetLocalIP()].ExternalUp[f] == true || Elevators[GetLocalIP()].ExternalDown[f] == true {
			return true
		}
	}
	return false
}

func RemoveOrder(floor, dir int) {
	if dir == 1 {
		Elevators[GetLocalIP()].ExternalUp[floor] = false
		Elevators[GetLocalIP()].InternalOrders[floor] = false
		//Send message
		driver.SetButtonLamp(floor, UP, false)
		driver.SetButtonLamp(floor, COMMAND, false)
		fmt.Println("inside remove Elevator with dir == 1")
		if floor == 3 {
			Elevators[GetLocalIP()].ExternalDown[floor] = false
			driver.SetButtonLamp(floor, DOWN, false)
		}
	} else if dir == -1 {
		Elevators[GetLocalIP()].ExternalDown[floor] = false
		Elevators[GetLocalIP()].InternalOrders[floor] = false
		//Send message
		driver.SetButtonLamp(floor, DOWN, false)
		driver.SetButtonLamp(floor, COMMAND, false)
		fmt.Println("inside remove Elevator with dir == -1")
		fmt.Println(floor)
		if floor == 0 {
			fmt.Println("Inside here?")
			driver.SetButtonLamp(floor, UP, false)
			Elevators[GetLocalIP()].ExternalUp[floor] = false
		}
	} else {
		//Send Message
		Elevators[GetLocalIP()].ExternalDown[floor] = false
		Elevators[GetLocalIP()].InternalOrders[floor] = false
		Elevators[GetLocalIP()].ExternalUp[floor] = false
		driver.SetButtonLamp(floor, COMMAND, false)
		driver.SetButtonLamp(floor, DOWN, false)
		driver.SetButtonLamp(floor, UP, false)
	}

}

func AddLocalOrder(floor, buttonType int) {
	var cheapestElevator string
	if buttonType != COMMAND{
		cheapestElevator = findCheapestElevator(floor)
	}
	switch buttonType {
	case UP:
		Elevators[cheapestElevator].ExternalUp[floor] = true
		//Send message type "ExternalUpOrderAdded"
		driver.SetButtonLamp(floor, buttonType, true)
		fmt.Println("Elevator added in ExternalUp queue")
	case DOWN:
		Elevators[cheapestElevator].ExternalDown[floor] = true
		//Send message type "ExternalDownOrderAdded"
		driver.SetButtonLamp(floor, buttonType, true)
		fmt.Println("Elevator added in ExternalDown queue")
	case COMMAND:
		Elevators[GetLocalIP()].InternalOrders[floor] = true
		//Send message type "lokalOrderAdded"
		driver.SetButtonLamp(floor, buttonType, true)
		fmt.Println("New internal order to floor:", floor, " added")
	}
}

func findCheapestElevator(floor int) string {
	/*length := 1 //len(Elevators)
	var costs[length]int
	i := 0
	for ip,info := range Elevators {
		costs[i] = calculateOrderCostForOnlyOneElevator(info.Floor, floor, info.Direction)
		i++
	}
	lowestnumber := 0
	for elev := 1; elev < len(Elevators); elev++{
		if costs[elev] < costs[lowestnumber]{
			lowestnumber = elev
		}
	}
	j := 0
	for ip,_ := range Elevators {
		if j == lowestnumber{
			return ip
		}
		j++
	}*/
	return GetLocalIP()
}

func calculateOrderCostForOnlyOneElevator(currFloor int, orderedFloor int, direction int) int {
	cost := 0
	if currFloor == -1 { //Hvis heisen er mellom etasjer
		cost++
	} else if direction != 0 { //Heis på etasje, men i full fart
		cost += 2
	}
	if currFloor < orderedFloor {
		for floor := currFloor; floor <= N_FLOORS; floor++ {
			cost++
		}
		if direction < 0 {
			cost += 5
		}
	}
	if currFloor > orderedFloor {
		for floor := N_FLOORS; floor >= currFloor; floor-- {
			cost++
		}
		if direction > 0 {
			cost += 5
		}
	}
	return cost
	
}

func EmptyQueue() bool {
	check := true
	for floor := 0; floor < driver.N_FLOORS; floor++ {
		if Elevators[GetLocalIP()].ExternalUp[floor] == true || Elevators[GetLocalIP()].ExternalDown[floor] == true || Elevators[GetLocalIP()].InternalOrders[floor] == true {
			check = false
		}
	}
	return check
}

/* 
func MessageReceiver(receivedMsgChan chan Message, sameFloorChan chan int, emptyQueueChan chan int){
	msg := <- receivedMsgChan
	switch msg.MessageType {
		case "IncomingElevator":
			ElevatorMsgHandler := <- AddElevator(msg.ReceiverIP, msg.Elevator)
			switch ElevatorMsgHandler {
				case "empty":
					emptyQueueChan <- msg.Elevator.Floor
				case "inSameFloor":
					sameFloorChan <- msg.Elevator.Floor
			}
		case "QueueDirection":
			elevators[msg.ReceiverIP].Direction = msg.Elevator.Direction
		case "newFloor":
			elevators[msg.ReceiverIP].LastPassedFloor = msg.Elevator.Floor //LastPassedFloor er en ny funksjon som må legges til
			elevators[msg.ReceiverIP].InFloor = true 
		case "completedElevator":
			ElevatorCompleted(msg.Elevator.Floor, msg.ReceiverIP)
	 	case "sendStatus":
		if message.SenderIP != myIP {
				_, exist := elevators[message.TargetIP]
				if !exist {
					newElev := Elevator{true, true, 1, 0, []bool{false, false, false, false}, []bool{false, false, false, false}, []bool{false, false, false, false}}
					elevators[message.TargetIP] = &newElev
				}
				elevators[message.TargetIP].InFloor = message.Elevator.InFloor
				elevators[message.TargetIP].LastPassedFloor = message.Elevator.LastPassedFloor
				elevators[message.TargetIP].Direction = message.Elevator.Direction

				for floor := 0; floor < N_FLOORS; floor++ {
					elevators[message.TargetIP].UpElevators[floor] = elevators[message.TargetIP].UpElevators[floor] || message.Elevator.UpElevators[floor]
					elevators[message.TargetIP].DownElevators[floor] = elevators[message.TargetIP].DownElevators[floor] || message.Elevator.DownElevators[floor]
					elevators[message.TargetIP].CommandElevators[floor] = elevators[message.TargetIP].CommandElevators[floor] || message.Elevator.CommandElevators[floor]
				}
				ElevatorInEmptyQueueChan <- 1
				lightUpdateChan <- 1
			}

		case "leftFloor":
			fmt.Printf("Heis %s har forlatt etasjen:\n", message.TargetIP)
			LeftFloor(message.TargetIP)
		}
	}
}
*/