package queue

import (
	"../eventHandler"
	"../driver"
)

func EmptyQueue bool {
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			if Msg.ExternalUp[i] == 0 || Msg.ExDownOrders[i] == 0 || Msg.InOrders[i] == 0 {
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
	if InternalOrders[PrevFloor] { 	//Should stop if the new order is the same as the previous (you're in the same floor)
		return 0
	}
	if dir { // direction == up
		//Check if (Prev.floor+1) which is current floor until the top. If we detect an order in the queue, the lift continues upward
		for i := PrevFloor + 1; i < N_FLOORS; i++ {
			if InternalOrders[i] || ExternalUp[i] || ExternalDown[i] {
				return 1 //lift goes up
			}
		}
		return -1
	}
	if dir == -1 { // or just use else
		for i := PrevFloor - 1; i >= 0; i-- { // just the opposite of the previous if statement
			if InternalOrders[i] || ExternalUp[i] || ExternalDown[i] {
				return -1 //lift goes down
			}
		}
		return 1
	}
}

func RemoveOrders() { //Orders that are finished
	InternalOrders[PrevFloor] = 0
	SetButtonLamp(PrevFloor,0,false) // resets previous floor which always will be done when function is called 
	if dir { // Resets all orders in up direction
		ExternalUp[PrevFloor] = 0
		SetButtonLamp(PrevFloor,1,false)
	}
	if dir == -1 {
		ExternalDown[PrevFloor] = 0
		SetButtonLamp(PrevFloor,-1,false)
	}
}

func AdddOrder() {
	
}

