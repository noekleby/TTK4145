package definitions

import (
	"time"
)

const (
	LocalPort  = 30000
	ServerPort = 20012
	MsgSize    = 1024
	UP         = 0
	DOWN       = 1
	COMMAND    = 2
	HeartBeatPort = 30113
	BroadcastPort = 30215
	MOTOR_SPEED = 2800
	N_FLOORS    = 4
	N_BUTTONS   = 3
)
const (
	IDLE = iota
	ELEVATING
	DOOR_OPEN
	LocalIP = "129.241.187.152"
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
	Order       Order
}

type Elevator struct {
	Active         	bool
	Floor          	int
	Direction      	int
	FsmState       	int
	NewlyInit		bool
	InternalOrders [N_FLOORS]bool
	ExternalUp     [N_FLOORS]bool
	ExternalDown   [N_FLOORS]bool
}

type Order struct {
	Buttontype int
	Floor      int
	FromIP     string
}
