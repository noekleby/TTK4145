package transManager

import (
    "../queue"
    "../network"
    "time"
    "fmt"
    "definitions"
    "builtin"
)

func messageTransmitter(msgType string, targetIP string, order Order) { 
	newMessage := Message{
		msgType,
		myIP,
		targetIP,
		*(elevators[targetIP]), // elevators er en map definert i queue (hos Truls)
		order,
	}
	network.BroadcastMessage(newMessage)
}

func MessageTRX(receiveChan chan Message) {

	storedChan := make(chan []byte)

	go network.Listen(network.GetListenSocket(), storedChan)
	go network.SendStatus(broadcastChan)

	for {
		buffer := <-receive
		RxMessage := Message{}
		err := json.Unmarshal(buffer, &RxMessage)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		if RxMessage.SenderIP != network.GetLocalIP() {
			receiveChan <- RxMessage
		}
	}
}

func HeartbeatTRX(newElevatorChan chan string, deadElevatorChan chan string) {

	storedChan := make(chan []byte, 1)
	heartbeats := make(map[string]*time.Time)
	go network.Listen(network.GetListenSocket(), storedChan)
	go network.SendHeartbeat()

	for {
		buffer := <-receive
		storeBeat := Heartbeat{}
		err := json.Unmarshal(buffer, &storeBeat)
		if err!= nil {
			fmt.Println("Error: ", err)
		}
		_, exist := heartbeats[storeBeat.Id]
		if exist {
			heartbeats[storeBeat.Id] = &storeBeat.Time
		} else {
			newElevatorChan <- storeBeat.Id
			heartbeats[storeBeat.Id] = &storeBeat.Time
		}
		for i, t := range heartbeats {
			dur := time.Since(*t)
			if dur.Seconds() > 1 {
				fmt.Println("Warning, time is running out: ", dur)
				deadElevatorChan <- i
				delete(heartbeats, i)
			}
		}
	}
}

func MessageReceiver(incommingMsgChan chan Message, orderOnSameFloorChan chan int, orderInEmptyQueueChan chan int) {
	for {
		message := <-incommingMsgChan
		switch message.MessageType {
		case "newOrder":
			i := queue.addExternalOrder(message.TargetIP, message.Order)
			switch i {
			case "empty":
				orderInEmptyQueueChan <- message.Order.Floor
			case "sameFloor":
				orderOnSameFloorChan <- message.Order.Floor
			}
		case "newDirection":
			elevators[message.TargetIP].Direction = message.Order.Type
		case "newFloor":
			elevators[message.TargetIP].LastPassedFloor = message.Order.Floor
			elevators[message.TargetIP].InFloor = true
		case "completedOrder":
			queue.OrderCompleted(message.Order.Floor, message.TargetIP)
		case "statusUpdate":
			if message.SenderIP != network.GetLocalIP( {
				_, exist := elevators[message.TargetIP]
				if !exist {
					newElev := Elevator{true, true, 1, 0, []bool{false, false, false, false}, []bool{false, false, false, false}, []bool{false, false, false, false}}
					elevators[message.TargetIP] = &newElev
				}
				elevators[message.TargetIP].InFloor = message.Elevator.InFloor
				elevators[message.TargetIP].LastPassedFloor = message.Elevator.LastPassedFloor
				elevators[message.TargetIP].Direction = message.Elevator.Direction

				for floor := 0; floor < N_FLOORS; floor++ {
					elevators[message.TargetIP].UpOrders[floor] = elevators[message.TargetIP].UpOrders[floor] || message.Elevator.UpOrders[floor]
					elevators[message.TargetIP].DownOrders[floor] = elevators[message.TargetIP].DownOrders[floor] || message.Elevator.DownOrders[floor]
					elevators[message.TargetIP].CommandOrders[floor] = elevators[message.TargetIP].CommandOrders[floor] || message.Elevator.CommandOrders[floor]
				}
				queue.orderInEmptyQueueChan <- 1
				queue.lightUpdateChan <- 1
			}

		case "leftFloor":
			fmt.Printf("Heis %s har forlatt etasjen:\n", message.TargetIP)
			queue.LeftFloor(message.TargetIP)
		}
	}
}

func HeartbeatReceiver(newElevatorChan chan string, deadElevatorChan chan string) {
	for {
		select {
		case IP := <-newElevatorChan:
			if IP != network.GetLocalIP( {
				fmt.Printf("Det er dukket opp en ny heis me IP: %s\n", IP)
				_, exist := elevators[IP]
				if exist {
					elevators[IP].Active = true

					messageTransmitter("statusUpdate", myIP, Order{-1, -1})
					time.Sleep(1 * time.Millisecond)

					messageTransmitter("statusUpdate", IP, Order{-1, -1})
				} else {
					newElev := Elevator{true, true, 1, 0, []bool{false, false, false, false}, []bool{false, false, false, false}, []bool{false, false, false, false}}
					elevators[IP] = &newElev
					messageTransmitter("statusUpdate", myIP, Order{-1, -1})
				}
			}
		case IP := <-deadElevatorChan:
			elevators[IP].Active = false
			fmt.Printf("Det er fjernet en heis med IP: %s\n", IP)
		}
	}
}

func Init(){
	runtime.GOMAXPROCS(runtime.NumCPU()) //sets the number of cpu cores the program can use simultaneously.
	//sets it here to Numcpu which is the number of cores available. 
 
	orderOnSameFloorChan := make(chan int)
	orderInEmptyQueueChan := make(chan int)
 
	newElevatorChan := make(chan string)
	deadElevatorChan := make(chan string)
	go HeartbeatTRX(newElevatorChan, deadElevatorChan)
	go HeartbeatReceiver(newElevatorChan, deadElevatorChan)
	
	receiveChan := make(chan Message)
	go MessageTRX(receiveChan)
    go queue.MessageReceiver(receiveChan, orderOnSameFloorChan, orderInEmptyQueueChan)
	
	time.Sleep(time.Second*5)
}