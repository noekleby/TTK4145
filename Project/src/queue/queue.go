package queue

import (
	"../driver"
	//"definitions"
	"fmt"
	//"../eventhandler"
)

const (
	UP = iota
	DOWN
	COMMAND
)

type Order struct {
	InternalOrders [driver.N_FLOORS]int
	ExternalUp     [driver.N_FLOORS]int
	ExternalDown   [driver.N_FLOORS]int
	PrevFloor      int
	dir            int
}

func (Order *Order) ShouldStop(floor, dir int) bool {
	if Order.InternalOrders[floor] == 1 {
		return true
	}
	if dir == 1 {
		if Order.ExternalUp[floor] == 1 || floor == driver.N_FLOORS-1 {
			return true
		} else if Order.QueueDirectionUp(floor) {
			return false
		} else {
			return true
		}
	} else if dir == -1 {
		if Order.ExternalDown[floor] == 1 || floor == 0 {
			return true
		} else if Order.QueueDirectionDown(floor) {
			return false
		} else {
			return true
		}
	}
	return true
}

func (Order *Order) SetDirection() {
	Order.dir = driver.GetDirection()
}

/*
func (Order *Order) QueueDirectionDown(floor int) int {
	//Already checked if there are any orders in queue and there are
	check := false
	if floor == 0 {
		return 1
	} else {
		for i := floor; i >= 0; i-- {
			if Order.InternalOrders[i] == 1 || Order.ExternalDown[i] == 1 {
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

func (Order *Order) QueueDirectionUp(floor int) int {
	check := false
	if floor == 3 {
		return -1
	} else {
		for i := floor; i < driver.N_FLOORS; i++ {
			if Order.InternalOrders[i] == 1 || Order.ExternalUp[i] == 1 {
				check = true
			}
		}
		if check == true {
			return 1
		} else {
			return -1
		}
	}
}*/

/*
func (Order *Order) QueueDirection(direction, floor int) int {
	if Order.EmptyQueue() == true {
		return 0

	} else if direction == 1 {
		return (Order.QueueDirectionUp(floor))

	} else if direction == -1 {
		return (Order.QueueDirectionDown(floor))
	} else {
		fmt.Println("Something wrong with queue.")
		return 0
	}
}*/

func (Order *Order) QueueDirection(direction, floor int) int {
	if Order.EmptyQueue() == true {
		return 0

	} else if direction == 1 {
		if Order.QueueDirectionUp(floor) {
			return 1
		} else if Order.QueueDirectionDown(floor) {
			return -1
		}
	} else if direction == -1 {
		if Order.QueueDirectionDown(floor) {
			return -1
		} else if Order.QueueDirectionUp(floor) {
			return 1
		}
	}
	return 0
}

func (Order *Order) QueueDirectionUp(floor int) bool {
	for f := floor + 1; f < driver.N_FLOORS; f++ {
		if Order.InternalOrders[f] == 1 || Order.ExternalUp[f] == 1 || Order.ExternalDown[f] == 1 {
			return true
		}
	}
	return false
}

func (Order *Order) QueueDirectionDown(floor int) bool {
	for f := floor - 1; f > -1; f-- {
		if Order.InternalOrders[f] == 1 || Order.ExternalUp[f] == 1 || Order.ExternalDown[f] == 1 {
			return true
		}
	}
	return false
}
func (Order *Order) RemoveOrder(floor, dir int) {
	if dir == 1 {
		Order.ExternalUp[floor] = 0
		Order.InternalOrders[floor] = 0
		driver.SetButtonLamp(floor, UP, false)
		driver.SetButtonLamp(floor, COMMAND, false)
		fmt.Println("inside remove order with dir == 1")
		if floor == 3 {
			Order.ExternalDown[floor] = 0
			driver.SetButtonLamp(floor, DOWN, false)
		}
	} else if dir == -1 {
		Order.ExternalDown[floor] = 0
		Order.InternalOrders[floor] = 0
		driver.SetButtonLamp(floor, DOWN, false)
		driver.SetButtonLamp(floor, COMMAND, false)
		fmt.Println("inside remove order with dir == -1")
		fmt.Println(floor)
		if floor == 0 {
			fmt.Println("Inside here?")
			driver.SetButtonLamp(floor, UP, false)
			Order.ExternalUp[floor] = 0
		}
	} else {
		driver.SetButtonLamp(floor, COMMAND, false)
		driver.SetButtonLamp(floor, DOWN, false)
		driver.SetButtonLamp(floor, UP, false)
		Order.ExternalDown[floor] = 0
		Order.InternalOrders[floor] = 0
		Order.ExternalUp[floor] = 0
	}

}

func (Order *Order) AddOrder(newOrder, New int) {
	switch New {
	case 0:
		Order.ExternalUp[newOrder] = 1
		driver.SetButtonLamp(newOrder, New, true)
		fmt.Println("Order added in ExternalUp queue")
		fmt.Println(Order.ExternalUp)
	case 1:
		Order.ExternalDown[newOrder] = 1
		driver.SetButtonLamp(newOrder, New, true)
		fmt.Println("Order added in ExternalDown queue")
		fmt.Println(Order.ExternalDown)
	case 2:
		Order.InternalOrders[newOrder] = 1
		driver.SetButtonLamp(newOrder, New, true)
		fmt.Println("Order added in Internal queue")
		fmt.Println(Order.InternalOrders)
	}
}

func (Order *Order) EmptyQueue() bool {
	check := true
	for i := 0; i < driver.N_FLOORS; i++ {
		if Order.ExternalUp[i] == 1 || Order.ExternalDown[i] == 1 || Order.InternalOrders[i] == 1 {
			check = false
		}
	}
	return check
}
