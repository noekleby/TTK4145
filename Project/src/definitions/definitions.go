package definitions

import (
	"time"
	."../driver"
)

const (

	LocalPort = 30000
	ServerPort = 20012
	MsgSize = 1024
	//HeartBeatPort = 30117
	
)

type Heartbeat struct {
	Id string
	Time time.Time
}

type Message struct {
	MessageTpe string
	SenderIP string
	TargetIP string //Which elevator that changes status
	Elevator Elevator
	Order Order
}

type Elevator struct {
	Active	bool
	Floor	int 
	Direction	int
	PrevFloor	int
	fsmState	int
	InternalOrders	[N_FLOORS]int
	ExternalUp	[N_FLOORS]int
	ExternalDown	[N_FLOORS]int
}

type Button_info struct {
	Button int
	Floor  int
}
