package definitions

import (
	"time"
	"../driver"
)

const (
	LocalPort = 30000
	ServerPort = 20012
	MsgSize = 1024
	HeartBeatPort = 30117
	
)

type Heartbeat struct {
	ID string
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
	Active          bool
	InFloor         bool
	Direction       int
	LastPassedFloor int

	UpOrders      []bool
	DownOrders    []bool
	CommandOrders []bool
}

type Button_info struct {
	Button int
	Floor  int
}

type Order struct {
	InternalOrders [driver.N_FLOORS]int
	ExternalUp     [driver.N_FLOORS]int
	ExternalDown   [driver.N_FLOORS]int
	PrevFloor      int
	dir            int
	Type  int
	Floor int
}

type ElevatorState struct {
	fsmState   int //State
	floor, dir int
	//destination int
}