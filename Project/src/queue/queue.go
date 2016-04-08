package queue

import (
	"fmt"
	"../driver"
	//"../eventhandler"
)

type order struct {
InternalOrders [driver.N_FLOORS] int
ExternalUp [driver.N_FLOORS] int
ExternalDown [driver.N_FLOORS] int
PrevFloor int
dir int 
}

func GetQueueInfo() order {
	var info order 
	return info
}
/*
func (order *order) GetIntQ(floor int) int {
	return InternalOrders[floor]
}

func (order *order) GetExUp(floor int) int {
	return ExternalUp[floor]
}
func (order *order) GetExDown(floor int) int {
	return ExternalDown[floor]
}
*/
func (order *order) EmptyQueue() bool {
	for i := 0; i < driver.N_FLOORS; i++ {
		if order.ExternalUp[i] == 0 || order.ExternalDown[i] == 0 || order.InternalOrders[i] == 0 {
			return true
		}
	}
	return false
}
// 1:  first checks if it should stop in the current floor
// 2:  if not, check orders furter in the current direction 
// 3:  if not, check orders in opposite direction in current floor
// 4:  if not, change direction 

func (order *order) QueueDirection() int{
	order.dir = driver.GetDirection()
	if order.EmptyQueue() == false {    //Should stop (= 0) if emptyqueue is true
		return 0
	}
	/*if order.InternalOrders[order.PrevFloor] { 	//Should stop if the new order is the same as the previous (you're in the same floor)
		return 0
	}*/
	if order.dir == 1 { // direction == up
		//Check if (Prev.floor+1) which is current floor until the top. If we detect an order in the queue, the lift continues upward
		for i := order.PrevFloor + 1; i < driver.N_FLOORS; i++ {
			if order.InternalOrders[i] == 1 || order.ExternalUp[i] == 1|| order.ExternalDown[i] == 1{
				return -1 //lift goes up
			} else {
				return 1
			}
		}
	}
	if order.dir == -1 { // or just use else
		for i := order.PrevFloor - 1; i >= 0; i-- { // just the opposite of the previous if statement
			if order.InternalOrders[i] == 1|| order.ExternalUp[i] == 1|| order.ExternalDown[i] == 1 {
				return 1 //lift goes down
			} else {
				return -1
			}
		}
	}
	return 0
}

func (order *order) RemoveOrders() { //Orders that are finished
	order.InternalOrders[order.PrevFloor] = 0
	driver.SetButtonLamp(order.PrevFloor,0,false) // resets previous floor which always will be done when function is called 
	if order.dir == 1 { // Resets all orders in up direction
		order.ExternalUp[order.PrevFloor] = 0
		driver.SetButtonLamp(order.PrevFloor,1,false)
	}
	if order.dir == -1 {
		order.ExternalDown[order.PrevFloor] = 0
		driver.SetButtonLamp(order.PrevFloor,-1,false)
	}
}

func AddOrder(newOrder, New int) {
	var order order
	switch New {
		case 0:
			order.ExternalUp[newOrder] = 1
			fmt.Println("Order added in ExternalUp queue")
		case 1:
			order.ExternalDown[newOrder] = 1
			fmt.Println("Order added in ExternalDown queue")
		case 2:
			order.InternalOrders[newOrder] = 1
			fmt.Println("Order added in Internal queue")
	}
}

