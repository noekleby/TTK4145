package queue

import (
	. "../definitions"
	"../driver"
	. "../network"
	"fmt"
)

const (
	N_FLOORS = 4
)

func ShouldStop(floor, dir int) bool {
	ip := Elevators[GetLocalIP()].CompletingOrdersForElevator
	fmt.Println(Elevators[ip].ExternalUp, Elevators[ip].InternalOrders, Elevators[ip].ExternalDown)
	if Elevators[ip].InternalOrders[floor] == true {
		return true
	}
	if dir == 1 {
		if Elevators[ip].ExternalUp[floor] == true || floor == driver.N_FLOORS-1 {
			return true
		} else if QueueDirectionUp(floor, ip) {
			return false
		} else {
			return true
		}
	} else if dir == -1 {
		if Elevators[ip].ExternalDown[floor] == true || floor == 0 {
			return true
		} else if QueueDirectionDown(floor, ip) {
			return false
		} else {
			return true
		}
	}
	return true
}

func QueueDirection(direction, floor int, ip string) int {
	if EmptyQueue(ip) == true {
		return 0

	} else if direction == 1 {
		if QueueDirectionUp(floor, ip) {
			return 1
		} else if QueueDirectionDown(floor, ip) {
			return -1
		}
	} else if direction == -1 {
		if QueueDirectionDown(floor, ip) {
			return -1
		} else if QueueDirectionUp(floor, ip) {
			return 1
		}
	}
	return 0
}

func QueueDirectionUp(floor int, ip string) bool {
	for f := floor + 1; f < driver.N_FLOORS; f++ {
		if Elevators[ip].InternalOrders[f] == true || Elevators[ip].ExternalUp[f] == true || Elevators[ip].ExternalDown[f] == true {
			return true
		}
	}
	return false
}

func QueueDirectionDown(floor int, ip string) bool {
	for f := floor - 1; f > -1; f-- {
		if Elevators[ip].InternalOrders[f] == true || Elevators[ip].ExternalUp[f] == true || Elevators[ip].ExternalDown[f] == true {
			return true
		}
	}
	return false
}

func RemoveRemoteOrder(floor int, direction int) {
	if direction == UP {
		Elevators[GetLocalIP()].ExternalUp[floor] = false
		driver.SetButtonLamp(floor, UP, false)
	} else if direction == DOWN {
		Elevators[GetLocalIP()].ExternalDown[floor] = false
		driver.SetButtonLamp(floor, DOWN, false)
	}
}

func AddRemoteOrder(IP string, elevator Elevator, order Order) {
	if order.Buttontype == UP {
		if !Elevators[IP].ExternalUp[order.Floor] && elevator.ExternalUp[order.Floor] {
			driver.SetButtonLamp(order.Floor, UP, true)
			Elevators[IP].ExternalUp[order.Floor] = elevator.ExternalUp[order.Floor]
		}
	} else {
		if !Elevators[IP].ExternalDown[order.Floor] && elevator.ExternalDown[order.Floor] {
			driver.SetButtonLamp(order.Floor, DOWN, true)
			Elevators[IP].ExternalDown[order.Floor] = elevator.ExternalDown[order.Floor]
		}
	}
}

func RemoveOrder(floor int, dir int) {
	order := Order{-1, -1, ""}
	if dir == 1 {
		Elevators[GetLocalIP()].ExternalUp[floor] = false
		Elevators[GetLocalIP()].InternalOrders[floor] = false
		newMsg := Message{"Remove order up", GetLocalIP(), "", *(Elevators[GetLocalIP()]), order}
		BroadcastMessage(newMsg)
		driver.SetButtonLamp(floor, UP, false)
		driver.SetButtonLamp(floor, COMMAND, false)
		if floor == 3 {
			Elevators[GetLocalIP()].ExternalDown[floor] = false
			driver.SetButtonLamp(floor, DOWN, false)
		}
	} else if dir == -1 {
		Elevators[GetLocalIP()].ExternalDown[floor] = false
		Elevators[GetLocalIP()].InternalOrders[floor] = false
		newMsg := Message{"Remove order down", GetLocalIP(), "", *(Elevators[GetLocalIP()]), order}
		BroadcastMessage(newMsg)
		driver.SetButtonLamp(floor, DOWN, false)
		driver.SetButtonLamp(floor, COMMAND, false)
		if floor == 0 {
			driver.SetButtonLamp(floor, UP, false)
			Elevators[GetLocalIP()].ExternalUp[floor] = false
		}
	} else {
		Elevators[GetLocalIP()].ExternalDown[floor] = false
		Elevators[GetLocalIP()].InternalOrders[floor] = false
		Elevators[GetLocalIP()].ExternalUp[floor] = false
		if floor == 0 {
			newMsg := Message{"Remove order up", GetLocalIP(), "", *(Elevators[GetLocalIP()]), order}
			BroadcastMessage(newMsg)
		} else if floor == 3 {
			newMsg := Message{"Remove order down", GetLocalIP(), "", *(Elevators[GetLocalIP()]), order}
			BroadcastMessage(newMsg)
		}
		driver.SetButtonLamp(floor, COMMAND, false)
		driver.SetButtonLamp(floor, DOWN, false)
		driver.SetButtonLamp(floor, UP, false)
	}

}

func AddLocalOrder(order Order) {
	var cheapestElevator string
	if order.Buttontype != COMMAND {
		cheapestElevator = findCheapestElevator(order.Floor)
		//cheapestElevator = "129.241.187.26"
		//fmt.Println("Inside addLocalOrder: ", cheapestElevator)
	}
	switch order.Buttontype {
	case UP:
		Elevators[cheapestElevator].ExternalUp[order.Floor] = true
		if order.FromIP != GetLocalIP() {
			newMsg := Message{"Add order", GetLocalIP(), cheapestElevator, *(Elevators[GetLocalIP()]), order}
			BroadcastMessage(newMsg)
		}
		driver.SetButtonLamp(order.Floor, order.Buttontype, true)
		fmt.Println("Elevator added in ExternalUp queue")
	case DOWN:
		Elevators[cheapestElevator].ExternalDown[order.Floor] = true
		if order.FromIP != GetLocalIP() {
			newMsg := Message{"Add order", GetLocalIP(), cheapestElevator, *(Elevators[GetLocalIP()]), order}
			BroadcastMessage(newMsg)
		}
		driver.SetButtonLamp(order.Floor, order.Buttontype, true)
		fmt.Println("Elevator added in ExternalDown queue")
	case COMMAND:
		Elevators[GetLocalIP()].InternalOrders[order.Floor] = true
		if order.FromIP != GetLocalIP() {
			newMsg := Message{"Add order", GetLocalIP(), cheapestElevator, *(Elevators[GetLocalIP()]), order}
			BroadcastMessage(newMsg)
		}
		driver.SetButtonLamp(order.Floor, order.Buttontype, true)
		fmt.Println("New internal order to floor:", order.Floor, " added")
	}
}

func findCheapestElevator(floor int) string { // Think this is our obstacle
	cheapestElevator := ""
	minCost := 9999
	for IP, elevator := range Elevators {
		if Elevators[IP].Active == true {
			cost := costFunction(elevator.Floor, floor, elevator)
			fmt.Println("Cost for order is ", cost, " for IP ", IP)
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
	return cheapestElevator // Does not calculate the first order?
}

func costFunction(currFloor int, orderedFloor int, elevator *Elevator) int {

	cost := 0

	if elevator.Direction == UP && orderedFloor < currFloor {
		cost += 5

	} else if elevator.Direction == DOWN && orderedFloor > currFloor {
		cost += 5

	}
	return cost
}

/*
func findCheapestElevator(floor int) string {
	//length := len(Elevators)
	costs := [5]int{999, 999, 999, 999, 999} // How can we solve this?
	//var costs []int   //Have to hardcode.. That does not work..
	i := 0
	for _, info := range Elevators {
		costs[i] = calculateOrderCostForOnlyOneElevator(info.Floor, floor, info.Direction)
		i++
		fmt.Println("Cost for order:", calculateOrderCostForOnlyOneElevator(info.Floor, floor, info.Direction))
	}
	fmt.Println("Cost for first IP: ", costs[0], " Cost for second IP: ", costs[1]) // Får printet ut riktig cost, men begge heisene stopper alltid når de skal innom 3.etasje
	lowestnumber := 0
	for elev := 1; elev < len(Elevators); elev++ {
		if costs[elev] < costs[lowestnumber] {
			lowestnumber = elev
		}
	}
	j := 0
	for ip, _ := range Elevators {
		if j == lowestnumber {
			fmt.Println("Returning from findCheapestElevator function:", ip) // Printer også ut riktig ip basert på hvem som har billigst cost
			return ip
		}
		j++
	}
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

}*/

func EmptyQueue(ip string) bool {
	check := true
	for floor := 0; floor < driver.N_FLOORS; floor++ {
		if Elevators[ip].ExternalUp[floor] == true || Elevators[ip].ExternalDown[floor] == true || Elevators[ip].InternalOrders[floor] == true {
			check = false
		}
	}
	return check
}
