package queue

import (
	"../driver"
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
	if dir == 1 {
		if Order.InternalOrders[floor] == 1 || Order.ExternalUp[floor] == 1 {
			return true
		} else {
			return false
		}
	} else {
		if Order.InternalOrders[floor] == 1 || Order.ExternalDown[floor] == 1 {
			return true
		} else {
			return false
		}
	}
}

func (Order *Order) SetDirection() {
	Order.dir = driver.GetDirection()
}

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
}

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
}

func (Order *Order) RemoveOrder(floor, dir int) {
	if dir == 1 {
		Order.ExternalUp[floor] = 0
		Order.InternalOrders[floor] = 0
		driver.SetButtonLamp(floor, UP, false)
		driver.SetButtonLamp(floor, COMMAND, false)
		fmt.Println("inside remove order with dir == 1")
	} else if dir == -1 {
		Order.ExternalDown[floor] = 0
		Order.InternalOrders[floor] = 0
		driver.SetButtonLamp(floor, DOWN, false)
		driver.SetButtonLamp(floor, COMMAND, false)
		fmt.Println("inside remove order with dir == -1")
	} else {
		fmt.Println("Can not remove order from queue if direction is set to zero")
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
