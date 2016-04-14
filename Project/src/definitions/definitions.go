package definitions

import (
	. "../driver"
	"time"
)

const (
	LocalPort  = 30000
	ServerPort = 20012
	MsgSize    = 1024
	UP         = 0
	DOWN       = 1
	COMMAND    = 2
	//HeartBeatPort = 30117

)
const (
	IDLE = iota
	ELEVATING
	DOOR_OPEN
)
const (
	HeartBeatPort = 30113
	BroadcastPort = 30215
)
var Elevators = map[string]*Elevator{}
var MessageBroadcastChan = make(chan Message)

type Heartbeat struct {
	Id   string
	Time time.Time
}

type Message struct {
	MessageType string
	SenderIP    string
	TargetIP    string //Cheapest Elevator
	Elevator    Elevator
}

type Elevator struct {
	Active         bool
	Floor          int
	Direction      int
	PrevFloor      int
	FsmState       int
	InternalOrders [N_FLOORS]bool
	ExternalUp     [N_FLOORS]bool
	ExternalDown   [N_FLOORS]bool
}
