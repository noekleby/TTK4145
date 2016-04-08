package queue

import (
	"../driver"
	"../eventhandler"
)

type order struc {
InternalOrders [N_FLOORS] int
ExternalUp [N_FLOORS] int
ExternalDown [N_FLOORS] int
PrevFloor, Dir int }

func EmptyQueue bool {
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			if order.ExternalUp[i] == 0 || order.ExternalDown[i] == 0 || order.InternalOrders[i] == 0 {
				return true
			}
		}
	}
	return false
}
// 1:  first checks if it should stop in the current floor
// 2:  if not, check orders furter in the current direction 
// 3:  if not, check orders in opposite direction in current floor
// 4:  if not, change direction 

func Direction () int{
	if EmptyQueue == 0 {    //Should stop (= 0) if emptyqueue is true
		return 0
	}
	if order.InternalOrders[PrevFloor] { 	//Should stop if the new order is the same as the previous (you're in the same floor)
		return 0
	}
	if order.dir { // direction == up
		//Check if (Prev.floor+1) which is current floor until the top. If we detect an order in the queue, the lift continues upward
		for i := order.PrevFloor + 1; i < driver.N_FLOORS; i++ {
			if order.InternalOrders[i] || order.ExternalUp[i] || order.ExternalDown[i] {
				return 1 //lift goes up
			}
		}
		return -1
	}
	if order.dir == -1 { // or just use else
		for i := order.PrevFloor - 1; i >= 0; i-- { // just the opposite of the previous if statement
			if order.InternalOrders[i] || order.ExternalUp[i] || order.ExternalDown[i] {
				return -1 //lift goes down
			}
		}
		return 1
	}
}

func RemoveOrders() { //Orders that are finished
	order.InternalOrders[order.PrevFloor] = 0
	SetButtonLamp(order.PrevFloor,0,false) // resets previous floor which always will be done when function is called 
	if order.dir { // Resets all orders in up direction
		order.ExternalUp[order.PrevFloor] = 0
		SetButtonLamp(order.PrevFloor,1,false)
	}
	if order.dir == -1 {
		ExternalDown[order.PrevFloor] = 0
		SetButtonLamp(order.PrevFloor,-1,false)
	}
}

func AddOrder(newOrder button_info) {
	switch newOrder.button {
		case 1:
			order.ExternalUp[newOrder.Floor] = 1
		case -1:
			order.ExternalDown[newOrder.Floor] = 1
		case 0:
			order.InternalOrders[newOrder.Floor] = 1
	}
}

