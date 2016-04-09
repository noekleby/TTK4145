package queue

import (
	"../driver"
	"fmt"
	//"../eventhandler"
)

type Order struct {
	InternalOrders [driver.N_FLOORS]int
	ExternalUp     [driver.N_FLOORS]int
	ExternalDown   [driver.N_FLOORS]int
	PrevFloor      int
	dir            int
}

func GetQueueInfo() Order {
	var info Order
	return info
}

func (info Order) GetIntQ(floor int) int {
	fmt.Println("Inside GetIntQ", info.InternalOrders)
	return info.InternalOrders[floor]
}

func (info Order) GetExUp(floor int) int {
	return info.ExternalUp[floor]
}
func (info Order) GetExDown(floor int) int {
	return info.ExternalDown[floor]
}

func (info Order) Add(floor, button int) {
	info.AddOrder(floor, button)
}

func (Order *Order) EmptyQueue() bool {
	for i := 0; i < driver.N_FLOORS; i++ {
		if Order.ExternalUp[i] == 0 || Order.ExternalDown[i] == 0 || Order.InternalOrders[i] == 0 {
			return true
		}
	}
	return false
}

// 1:  first checks if it should stop in the current floor
// 2:  if not, check Orders furter in the current direction
// 3:  if not, check Orders in opposite direction in current floor
// 4:  if not, change direction

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

func (Order *Order) RemoveOrders() { //Orders that are finished
	Order.InternalOrders[Order.PrevFloor] = 0
	driver.SetButtonLamp(Order.PrevFloor, 0, false) // resets previous floor which always will be done when function is called
	if Order.dir == 1 {                             // Resets all Orders in up direction
		Order.ExternalUp[Order.PrevFloor] = 0
		driver.SetButtonLamp(Order.PrevFloor, 1, false)
	}
	if Order.dir == -1 {
		Order.ExternalDown[Order.PrevFloor] = 0
		driver.SetButtonLamp(Order.PrevFloor, -1, false)
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
