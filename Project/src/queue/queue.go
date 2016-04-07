package queue

import (
	"../eventHandler"
	"../driver"
)
//var cost_matrix = [driver.N_FLOORS*driver.N_BUTTONS]

func EmptyQueue bool {
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			if (driver.Msg.ExternalUp[i] == 0 || driver.Msg.ExDownOrders[i] == 0 || driver.Msg.InOrders[i] == 0) {
				return true
			}
		}
	}
	return false
}

/*func (q* queue) sortQueue int(){
	for i := 0; i < length(q); i++ {
		cost_matrix(i) = cost(currFloor, targetFloor, direction)
	}
	for j := 0; j < 
}

func NewDirection  int () {
	if EmptyQueue {
		driver.Init()
	} else {

	}

}
*/

// 1:  first checks if it should stop in the current floor
// 2:  if not, check orders furter in the current direction 
// 3:  if not, check orders in opposite direction in current floor
// 4:  if not, change direction 

func changeDirection () int{
	//Should stop (= 0) if emptyqueue is true
	//Should stop if the new order is the same as the previous (you're in the same floor)
	//Make if statement for msg = direction up
	//	Make loop
	//Make if statement for msg = direction down (or just use else)
	//	Make loop
}

func ordersInDirection () bool {
	//Checks first in up direction
	//then in down direction using if statements
}

