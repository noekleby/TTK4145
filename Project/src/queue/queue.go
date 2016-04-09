package queue

import (
	"../driver"
	"fmt"
	//"../eventhandler"
)

Const (
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
	if dir == 1{
		if InternalOrders[floor] == 1 || ExternalUp[floor] == 1  {
			return true 
		} else {
			return false 
		}
	} else {
		if InternalOrders[floor] == 1 || ExternalDown[floor] == 1 {
			return true 
		} else {
			return false 
		} 
	}
}

func (Order *Order) QueueDirection() int {
	Order.dir = driver.GetDirection()
	if Order.EmptyQueue() == false { //Should stop (= 0) if emptyqueue is true
		return 0
	}
	/*if Order.InternalOrders[Order.PrevFloor] { 	//Should stop if the new Order is the same as the previous (you're in the same floor)
		return 0
	}*/
	if Order.dir == 1 { // direction == up
		//Check if (Prev.floor+1) which is current floor until the top. If we detect an Order in the queue, the lift continues upward
		for i := Order.PrevFloor + 1; i < driver.N_FLOORS; i++ {
			if Order.InternalOrders[i] == 1 || Order.ExternalUp[i] == 1 || Order.ExternalDown[i] == 1 {
				return -1 //lift goes up
			} else {
				return 1
			}
		}
	}
	if Order.dir == -1 { // or just use else
		for i := Order.PrevFloor - 1; i >= 0; i-- { // just the opposite of the previous if statement
			if Order.InternalOrders[i] == 1 || Order.ExternalUp[i] == 1 || Order.ExternalDown[i] == 1 {
				return 1 //lift goes down
			} else {
				return -1
			}
		}
	}
	return 0
}

func (Order *Order) RemoveOrder(floor, dir int) {
	if dir == 1 {
		ExternalUp[floor] = 0
		InternalOrders[floor] = 0 
		driver.SetButtonLamp(floor, UP, false)
		driver.SetButtonLamp(floor, COMMAND, false)
	} 
	if dir == -1 {
		ExternalDown[floor] = 0 
		InternalOrders[floor] = 0 
		driver.SetButtonLamp(floor, DOWN, false)
		driver.SetButtonLamp(floor, COMMAND, false)
	} else {
		fmt.Println("Can not remove order from queue if direction is set to zero")
	}

}


func (Order *Order) AddOrder(newOrder, New int) {
	switch New {
	case 0:
		Order.ExternalUp[newOrder] = 1
		fmt.Println("Order added in ExternalUp queue")
	case 1:
		Order.ExternalDown[newOrder] = 1
		fmt.Println("Order added in ExternalDown queue")
	case 2:
		Order.InternalOrders[newOrder] = 1
		fmt.Println("Order added in Internal queue")
		fmt.Println(Order.InternalOrders[0], Order.InternalOrders[1], Order.InternalOrders[2], Order.InternalOrders[3])
	}
}







