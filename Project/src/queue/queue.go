package queue

import (
	. "../definitions"
	"../driver"
	. "../network"
	"fmt"
)


func ShouldStop(floor, direction int) bool {
	if Elevators[LocalIP].InternalOrders[floor] == true {
		return true
	} else if direction == 1 {
		if Elevators[LocalIP].ExternalUp[floor] == true || floor == N_FLOORS-1 {
			return true
		} else if queueDirectionUp(floor) {
			return false
		} else {
			return true
		}
	} else if direction == -1 {
		if Elevators[LocalIP].ExternalDown[floor] == true || floor == 0 {
			return true
		} else if queueDirectionDown(floor) {
			return false
		} else {
			return true
		}
	}
	return true
}


func NextDirection(direction, floor int) int {
	if EmptyQueue() {
		return 0

	} else if direction == 1 {
		if queueDirectionUp(floor) {
			return 1
		} else if Elevators[LocalIP].ExternalDown[floor]{
			return 0 
		} else if queueDirectionDown(floor) {
			return -1
		}
	} else if direction == -1 {
		if queueDirectionDown(floor) {
			return -1
		} else if Elevators[LocalIP].ExternalUp[floor]{
			return 0
		}else if queueDirectionUp(floor) {
			return 1
		}
	}
	return 0
}

func RemoveLocalOrder(floor int, direction int, lightEventChan chan int) {
	order := Order{-1, -1, ""}

	if direction == 1 && driver.GetLampSignal(UP, floor) == 1 {
		Elevators[LocalIP].ExternalUp[floor] = false
		Elevators[LocalIP].InternalOrders[floor] = false

		newMsg := Message{"Remove order up", LocalIP, "", *(Elevators[LocalIP]), order}
		BroadcastMessage(newMsg)
		lightEventChan <- 1

	} else if direction == 1 && driver.GetLampSignal(DOWN, floor) == 1 {
		Elevators[LocalIP].ExternalDown[floor] = false
		Elevators[LocalIP].InternalOrders[floor] = false
		
		newMsg := Message{"Remove order down", LocalIP, "", *(Elevators[LocalIP]), order}
		BroadcastMessage(newMsg)
		lightEventChan <- 1

	} else if direction == -1 && driver.GetLampSignal(DOWN, floor) == 1 {
		Elevators[LocalIP].ExternalDown[floor] = false
		Elevators[LocalIP].InternalOrders[floor] = false

		newMsg := Message{"Remove order down",LocalIP, "", *(Elevators[LocalIP]), order}
		BroadcastMessage(newMsg)
		lightEventChan <- 1

	} else if direction == -1 && driver.GetLampSignal(UP, floor) == 1 {
		Elevators[LocalIP].ExternalUp[floor] = false
		Elevators[LocalIP].InternalOrders[floor] = false

		newMsg := Message{"Remove order up", LocalIP, "", *(Elevators[LocalIP]), order}
		BroadcastMessage(newMsg)
		lightEventChan <- 1

	} else {
		Elevators[LocalIP].ExternalDown[floor] = false
		Elevators[LocalIP].InternalOrders[floor] = false
		Elevators[LocalIP].ExternalUp[floor] = false

		if floor == 0 {
			newMsg := Message{"Remove order up", LocalIP, "", *(Elevators[LocalIP]), order}
			BroadcastMessage(newMsg)
		} else if floor == 3 {
			newMsg := Message{"Remove order down", LocalIP, "", *(Elevators[LocalIP]), order}
			BroadcastMessage(newMsg)
		}
		lightEventChan <- 1
	}
}


func AddLocalOrder(order Order, lightEventChan chan int) {
	var cheapestElevator string
	if order.Buttontype != COMMAND {
		cheapestElevator = findCheapestElevator(order)
	}
	switch order.Buttontype {

	case UP:
		Elevators[cheapestElevator].ExternalUp[order.Floor] = true
		if order.FromIP != LocalIP {
			newMsg := Message{"Add order", LocalIP, cheapestElevator, *(Elevators[LocalIP]), order}
			BroadcastMessage(newMsg)
		}
		lightEventChan <- 1

	case DOWN:
		Elevators[cheapestElevator].ExternalDown[order.Floor] = true
		if order.FromIP != LocalIP {
			newMsg := Message{"Add order", LocalIP, cheapestElevator, *(Elevators[LocalIP]), order}
			BroadcastMessage(newMsg)
		}
		lightEventChan <- 1

	case COMMAND:
		Elevators[LocalIP].InternalOrders[order.Floor] = true
		lightEventChan <- 1
	}
}



func EmptyQueue() bool {
	check := true
	for floor := 0; floor < N_FLOORS; floor++ {
		if Elevators[LocalIP].ExternalUp[floor] == true || Elevators[LocalIP].ExternalDown[floor] == true || Elevators[LocalIP].InternalOrders[floor] == true {
			check = false
		}
	}
	return check
}


func queueDirectionUp(floor int) bool {
	for f := floor + 1; f < N_FLOORS; f++ {
		if Elevators[LocalIP].InternalOrders[f] == true || Elevators[LocalIP].ExternalUp[f] == true || Elevators[LocalIP].ExternalDown[f] == true {
			return true
		}
	}
	return false
}

func queueDirectionDown(floor int) bool {
	for f := floor - 1; f > -1; f-- {
		if Elevators[LocalIP].InternalOrders[f] == true || Elevators[LocalIP].ExternalUp[f] == true || Elevators[LocalIP].ExternalDown[f] == true {
			return true
		}
	}
	return false
}

func findCheapestElevator(order Order) string { // Think this is our obstacle
	cheapestElevator := ""
	minCost := 9999
	for IP, elevator := range Elevators {
		if Elevators[IP].Active == true {
			cost := calculateCost(order, elevator)
			if cost < minCost {
				minCost = cost
				cheapestElevator = IP
			}
			if cost == 0 {
				break
			}
		}
	}
	fmt.Println("The cheapest IP is ", cheapestElevator)
	return cheapestElevator
}

func calculateCost(order Order, elevator *Elevator) int {
	cost := 0

	if order.Buttontype == DOWN {
		for floor := elevator.Floor; floor == 0; floor-- {
			if elevator.ExternalDown[floor] {
				cost += 1
			}
		}
	}

	if order.Buttontype == UP {
		for floor := elevator.Floor; floor == (N_FLOORS - 1); floor++ {
			if elevator.ExternalUp[floor] {
				cost += 1
			}
		}
	}

	if order.Buttontype == UP && elevator.ExternalDown[order.Floor] {
		cost += 7
	}

	if order.Buttontype == DOWN && elevator.ExternalUp[order.Floor] {
		cost += 7
	}

	for floor := 0; floor < N_FLOORS; floor++ {
		if elevator.InternalOrders[floor] {
			cost += 1
		}
	}

	if elevator.Direction == UP && order.Floor < elevator.Floor {
		cost += 5

	} else if elevator.Direction == DOWN && order.Floor > elevator.Floor {
		cost += 5

	}

	return cost
}
