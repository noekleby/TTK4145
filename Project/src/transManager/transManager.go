package transManager

import (
    "../queue"
    "../network"
    "net"
    "time"
    "fmt"
    "definitions"
)

func MessageTx(receiveChan chan Message) {

	storedChan := make(chan []byte)

	go network.Listen(network.GetListenSocket(), storedChan)
	go network.sendStatus(broadcastChan)

	for {
		RxMessageBs := <-receive
		RxMessage := Message{}
		err := json.Unmarshal(RxMessageBs, &RxMessage)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		if RxMessage.SenderIP != GetLocalIP() {
			receiveChan <- RxMessage
		}
	}
}

func HeartbeatTx(newElevatorChan chan string, deadElevatorChan chan string) {

	storedChan := make(chan []byte, 1)
	heartbeats := make(map[string]*time.Time)
	go network.Listen(network.GetListenSocket(), storedChan)
	go network.sendHeartBeat()

	for {
		otherBeatBs := <-receive
		otherBeat := network.Heartbeat{}
		err := json.Unmarshal(otherBeatBs, &otherBeat)
		if err!= nil {
			fmt.Println("Error: ", err)
		}
		_, exist := heartbeats[otherBeat.Id]
		if exist {
			heartbeats[otherBeat.Id] = &otherBeat.Time
		} else {
			newElevatorChan <- otherBeat.Id
			heartbeats[otherBeat.Id] = &otherBeat.Time
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
			i := addExternalOrder(message.TargetIP, message.Order)
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
			OrderCompleted(message.Order.Floor, message.TargetIP)
		case "statusUpdate":
			if message.SenderIP != myIP {
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
				orderInEmptyQueueChan <- 1
				lightUpdateChan <- 1
			}

		case "leftFloor":
			fmt.Printf("Heis %s har forlatt etasjen:\n", message.TargetIP)
			LeftFloor(message.TargetIP)
		}
	}
}

func HeartbeatReceiver(newElevatorChan chan string, deadElevatorChan chan string) {
	for {
		select {
		case IP := <-newElevatorChan:
			if IP != myIP {
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
	go HeartbeatTx(newElevatorChan, deadElevatorChan)
	go HeartbeatReceiver(newElevatorChan, deadElevatorChan)
	
	receiveChan := make(chan Message)
	go MessageTx(receiveChan)
    go queue.MessageReceiver(receiveChan, orderOnSameFloorChan, orderInEmptyQueueChan)
	
	time.Sleep(time.Second*5)
}