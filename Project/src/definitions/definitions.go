package definitions

import (
	"time"
	."../driver"
)

const (

	LocalPort = 30000
	ServerPort = 20012
	MsgSize = 1024
	UP = 0
	DOWN = 1
	COMMAND = 2
	//HeartBeatPort = 30117
	
)
const (
	IDLE = iota
	ELEVATING
	DOOR_OPEN
)

var Elevators = map[string]*Elevator{}

type Heartbeat struct {
	Id string
	Time time.Time
}

type Message struct {
	MessageTpe string
	SenderIP string
	TargetIP string //Which elevator that changes status
	Elevator Elevator
	Button_info	Button_info
}

type Elevator struct {
	Active	bool
	Floor	int 
	Direction	int
	PrevFloor	int
	FsmState	int
	InternalOrders	[N_FLOORS]bool
	ExternalUp	[N_FLOORS]bool
	ExternalDown	[N_FLOORS]bool
}

type Button_info struct {
	Button int
	Floor  int
}


