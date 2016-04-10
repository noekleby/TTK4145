package eventhandler

import (
	"../driver"
	"fmt"
	"time"
	"../definitions"
)

var button = [driver.N_FLOORS][driver.N_BUTTONS]int{
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
}

func CheckEvents(floorChannel chan int, buttonChannel chan Button_info) {
	go FloorEventCheck(floorChannel)
	go ButtonEventCheck(buttonChannel)
}

func FloorEventCheck(event chan int) {
	prevFloor := 0
	for {
		newFloor := driver.GetFloorSignal()
		if newFloor != prevFloor {
			prevFloor = newFloor
			event <- newFloor
		}
	}
}

func ButtonEventCheck(event chan Button_info) {
	prevButton := button
	newButton := button
	fmt.Println(newButton)
	for {
		fmt.Println("Inside forever")
		for f := 0; f < driver.N_FLOORS; f++ {
			for b := 0; b < driver.N_BUTTONS; b++ {
				newButton[f][b] = driver.ElevGetButtonSignal(b, f)
				fmt.Println(newButton, prevButton)
				if (newButton[f][b] != prevButton[f][b]) && (newButton[f][b] == 1) {
					prevButton[f][b] = newButton[f][b]
					var buttonInf Button_info
					buttonInf.Button = b
					buttonInf.Floor = f
					event <- buttonInf
				}
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
}

/*
//----------------------------------------------------- eventuelt putte i en egen definition module
type MSG struct{
	MsgType			int
	State 			int
	PrevFloor 		int
	Dir   			int 	//never 0.
	ExUpOrders 		[N_FLOORS]int //External orders
	ExDownOrders	[N_FLOORS]int
	InOrders		[N_FLOORS]int //internal orders
}

type Order struct{
	Floor 	int
	Button 	int
}

type Udp_message struct {
	Raddr  string //if receiving raddr=senders address, if sending raddr should be set to "broadcast" or an ip:port
	Data   []byte //TODO: implement another encoding, strings are meh
	Length int    //length of received data, in #bytes // N/A for sending
}

var Msg = MSG{}
//---------------------------------------------------------

func InternalOrderDetector(orderChan chan Order) {
	var currSignalMatrix 	[3][N_FLOORS]int
	var prevSignalMatrix 	[3][N_FLOORS]int

	for {
		for floor:=0; floor < N_FLOORS; floor++ {
			for button:=0; button < N_FLOORS; button++ {
				currSignalMatrix[button][floor] = elevGetButtonSignal(button,floor) //fra driveren elev.go
				if (currSignalMatrix[button][floor] == 1 && prevSignalMatrix[button][floor] == 0) { //Hvis get button er 1 og det ikke finnes ordre for etasjen fra før, legges det inn en ny ordre
					orderChan <- Order{floor, button} //Sender Order på kanalen
				}
				prevSignalMatrix[button][floor] = currSignalMatrix[button][floor]
			}
		}
		time.Sleep(10*time.Millisecond)
	}
}


func floorReached(floorReachedChan chan int) {
	var prevFloor = elevGetFloorSignal() //Trenger en getfloorsignal-funksjon i driveren (som er elev.go midlertidig)
	for {
		if (elevGetFloorSignal() != -1 && prevFloor == -1) { //-1 er idle
			floorReachedChan <- elevGetFloorSignal() //Legger hvilke etasje man er i inn i floor reached-kanalen
		}
		prevFloor = elevGetFloorSignal()
		time.Sleep(10*time.Millisecond)
	}
}

func EmptyQueue() bool {
	return ExactlyOneOrder() //Funksjon som ligger i queue, som sjekker om man har en og bare en ordre. Returnerer da true.
}

func NewOrderInCurrentFloorEventDetector(order Order) bool {
	return order.Floor == Msg.PrevFloor
}

func ExternalOrdersUpdate(otherLift MSG){ //går gjennom panelene på utsiden av heisen i alle etasjene. Fjerner og legger til ordre (og styrer lampene)
	switch otherLift.MsgType{
	case AddOrders:
		for i:=0; i<N_FLOORS; i++ {
			if (otherLift.ExUpOrders[i] == 1) {
				Msg.ExUpOrders[i] = 1
				elevSetButtonLamp(i, 1, 1)
			}
			if	(otherLift.ExDownOrders[i] == 1) {
				Msg.ExDownOrders[i] = 1
				elevSetButtonLamp(i, -1, 1)
			}
		}

	case RemoveOrders:
		for i:=0; i<N_FLOORS; i++ {
			if (otherLift.ExUpOrders[i] == 0) {
				Msg.ExUpOrders[i] = 0
				elevSetButtonLamp(i, 1, 0)
			}
			if	(otherLift.ExDownOrders[i] == 0) {
				Msg.ExDownOrders[i] = 0
				elevSetButtonLamp(i, -1, 0)
			}
		}
	}
}


func EventHandler(timerChan chan string, timeOutChan chan int, send_ch, receive_ch chan Udp_message) {
	orderChan := make(chan Order) // lager channels
	floorReachedChan := make(chan int)
	go InternalOrderDetector(orderChan)
	go floorReached(floorReachedChan)


	for {

		select {

		case UDP_Rec := <- receive_ch: //mottas og legges i UDP_Rec ; Mottar her en ny ordre

			fmt.Println("HEIHEIHEHEHI", Laddr.String())

			if (Laddr.String() != UDP_Rec.Raddr) { // Sjekker om den nye ordren er forskjellig fra hva man har fra før
				fmt.Println("beat2")
				fmt.Println(UDP_Rec.Raddr)
				Dec_Msg := DecodeMsg(UDP_Rec.Data, UDP_Rec.Length)

				UpDateOrders(Dec_Msg) // Oppdaterer ordrene
				fmt.Println(Dec_Msg)
			}

		case order := <- orderEventChannel: //mottas og legges i order
			AddOrder(order)
			PrintMsg()

			Udp_msg.Data = EncodeMsg(Msg)
			send_ch <- Udp_msg //Sender udp_msg

			if (EmptyQueue()) {
				fmt.Println("NewOrderInEmptyQueue")
				EmptyQueue(timerChan)
				fmt.Println("Event : NewOrderInEmptyQueue")
			}
			if (NewOrderInCurrentFloorEventDetector(order)) {
				NewOrderInCurrentFloor(timerChan)
				fmt.Println("Event : NewOrderInCurrentFloor")
			}

		case floor := <- floorReachedEventChannel: //mottas og legges i floor ; Skjer hvis man har kommet til riktig etasje
			Msg.PrevFloor = floor
			fmt.Println("Event : New floor reached :", floor)
			stopped := false
			FloorReached(timerChan, stopped)
			if stopped {
				Msg.MsgType = REMOVE_ORDERS
				Udp_msg.Data = EncodeMsg(Msg)
				send_ch <- Udp_msg
				Msg.MsgType = NOTHING
			}

		case <- timerChan: // Resetter timerChan (definert i timer.go) / tømmer timerchan
			TimerOut()  // gir "time out" hvis man er for lenge i en state
		}
	}
}*/
