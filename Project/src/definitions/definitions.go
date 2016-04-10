package definitions

const (
	localPort := 30000
	serverPort := 20005
	msgSize := 1024
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

type Order struct {
	Type  int
	Floor int
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
}

type ElevatorState struct {
	fsmState   int //State
	floor, dir int
	//destination int
}