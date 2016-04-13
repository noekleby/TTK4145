package queue

import (
	. "../definitions"
	"../driver"
	"fmt"
	//"../eventhandler"
	. "../network"
)

const (
	N_FLOORS = 4
)

func ShouldStop(floor, dir int) bool {
	if Elevators[GetLocalIP()].InternalOrders[floor] == true {
		return true
	}
	if dir == 1 {
		if Elevators[GetLocalIP()].ExternalUp[floor] == true || floor == driver.N_FLOORS-1 {
			return true
		} else if QueueDirectionUp(floor) {
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
		if QueueDirectionDown(floor) {
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

func RemoveRemoteOrder(floor int, direction int) {
	if direction == UP {
		Elevators[GetLocalIP()].ExternalUp[floor] = false
		driver.SetButtonLamp(floor, UP, false)
	} else if direction == DOWN {
		Elevators[GetLocalIP()].ExternalDown[floor] = false
		driver.SetButtonLamp(floor, DOWN, false)
	}
}

func AddRemoteOrder(IP string, queue [driver.N_FLOORS]bool, direction int) {
	for i := 0; i < N_FLOORS; i++ {
		if direction == UP {
			if !Elevators[IP].ExternalUp[i] && queue[i] {
				driver.SetButtonLamp(i, UP, true)
				Elevators[IP].ExternalUp[i] = queue[i]
			}
		} else {
			if !Elevators[IP].ExternalDown[i] && queue[i] {
				driver.SetButtonLamp(i, DOWN, true)
				Elevators[IP].ExternalDown[i] = queue[i]
			}
		}
	}
}

func RemoveOrder(floor int, dir int) {
	if dir == 1 {
		Elevators[GetLocalIP()].ExternalUp[floor] = false
		Elevators[GetLocalIP()].InternalOrders[floor] = false
		newMsg := Message{"Remove order up", GetLocalIP(), "", *(Elevators[GetLocalIP()])}
		BroadcastMessage(newMsg)
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
		newMsg := Message{"Remove order down", GetLocalIP(), "", *(Elevators[GetLocalIP()])}
		BroadcastMessage(newMsg)
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
	if buttonType != COMMAND {
		cheapestElevator = findCheapestElevator(floor)
	}
	switch buttonType {
	case UP:
		Elevators[cheapestElevator].ExternalUp[floor] = true
		newMsg := Message{"Add order up", GetLocalIP(), cheapestElevator, *(Elevators[GetLocalIP()])}
		BroadcastMessage(newMsg)
		driver.SetButtonLamp(floor, buttonType, true)
		fmt.Println("Elevator added in ExternalUp queue")
	case DOWN:
		Elevators[cheapestElevator].ExternalDown[floor] = true
		newMsg := Message{"Add order down", GetLocalIP(), cheapestElevator, *(Elevators[GetLocalIP()])}
		BroadcastMessage(newMsg)
		driver.SetButtonLamp(floor, buttonType, true)
		fmt.Println("Elevator added in ExternalDown queue")
	case COMMAND:
		Elevators[GetLocalIP()].InternalOrders[floor] = true
		newMsg := Message{"Add internal order", GetLocalIP(), cheapestElevator, *(Elevators[GetLocalIP()])}
		BroadcastMessage(newMsg)
		driver.SetButtonLamp(floor, buttonType, true)
		fmt.Println("New internal order to floor:", floor, " added")
	}
}

func findCheapestElevator(floor int) string {
	//length := len(Elevators)
	costs := [5]int{999, 999, 999, 999, 999} // How can we solve this?
	//var costs []int   //Have to hardcode.. That does not work..
	i := 0
	for _, info := range Elevators {
		costs[i] = calculateOrderCostForOnlyOneElevator(info.Floor, floor, info.Direction)
		i++
	}
	lowestnumber := 0
	for elev := 1; elev < len(Elevators); elev++ {
		if costs[elev] < costs[lowestnumber] {
			lowestnumber = elev
		}
	}
	j := 0
	for ip, _ := range Elevators {
		if j == lowestnumber {
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
