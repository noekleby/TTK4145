package network

import (
	. "../definitions"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)


func GetLocalIP() string {
	addrs, error := net.InterfaceAddrs()
	if error != nil {
		fmt.Println("error:", error)
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func BroadcastMessage(message Message) {
	//printMessage(message)
	MessageBroadcastChan <- message
}

func MessageBroadcast(MessageBroadcastChan chan Message) {

	bufferSendChan := make(chan []byte)
	go udpSend(bufferSendChan, BroadcastPort)

	for {
		msg := <-MessageBroadcastChan
		buffer, error := json.Marshal(msg)
		if error != nil {
			fmt.Println("Error:", error)
			time.Sleep((4 * time.Second))
		}
		bufferSendChan <- buffer
	}
}

func MessageReciever(messageRecieveChan chan Message) {

	bufferRecieveChan := make(chan []byte)
	go udpRecieve(bufferRecieveChan, BroadcastPort)

	for {
		buffer := <-bufferRecieveChan
		msg := Message{}
		error := json.Unmarshal(buffer, &msg)
		if error != nil {
			fmt.Println("Error:", error)
			time.Sleep((4 * time.Second))
		}
		if msg.SenderIP != GetLocalIP() {
			//printMessage(msg)
			messageRecieveChan <- msg
		}
	}
}

func HeartbeatEventCheck(newElevatorChan chan string, deadElevatorChan chan string) {

	receive := make(chan []byte, 1)
	heartbeats := make(map[string]*time.Time)
	go udpRecieve(receive, HeartBeatPort)

	for {
		if GetLocalIP() != "" {
			otherBeatBuffer := <-receive
			otherBeat := Heartbeat{}
			error := json.Unmarshal(otherBeatBuffer, &otherBeat)
			if error != nil {
				fmt.Println("error:", error)
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
					fmt.Println("Warning:", dur)
					deadElevatorChan <- i
					delete(heartbeats, i)
				}
			}
		} else {
		deadElevatorChan <- LocalIP
		}
	}
}

func SendHeartBeat() {
	send := make(chan []byte, 1)
	go udpSend(send, HeartBeatPort)

	for {
		myBeat := Heartbeat{GetLocalIP(), time.Now()}
		myBeatBuffer, error := json.Marshal(myBeat)

		if error != nil {
			fmt.Println("error:", error)
		}
		send <- myBeatBuffer
		time.Sleep(100 * time.Millisecond)
	}
}

func getTransmitSocket(port int) *net.UDPConn {

	serverAddress, error := net.ResolveUDPAddr("udp", fmt.Sprintf("129.241.187.255:%d", port))
	if error != nil {
		fmt.Println("error:", error)
	}
	transmitSocket, _ := net.DialUDP("udp", nil, serverAddress)
	if error != nil {
		fmt.Println("error:", error)
	}
	return transmitSocket
}

func getListenSocket(port int) *net.UDPConn {
	localAddress, error := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if error != nil {
		fmt.Println("error:", error)
	}
	listenSocket, error := net.ListenUDP("udp", localAddress)
	if error != nil {
		fmt.Println("error:", error)
	}
	return listenSocket
}

func udpRecieve(msg chan []byte, port int) {
	for {
		socket := getListenSocket(port)
		buffer := make([]byte, 1024)
		n, _, error := socket.ReadFromUDP(buffer)
		if error != nil {
			fmt.Println("error:", error)
		}
		buffer = buffer[:n]
		msg <- buffer
		socket.Close()
	}
}

func udpSend(msg chan []byte, port int) {
	for {
		socket := getTransmitSocket(port)
		buffer := <-msg
		if Elevators[LocalIP].Active == true {
			socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
			_, error := socket.Write(buffer)
			if error != nil {
				fmt.Println("error:", error)
			}
		}
		socket.Close()
	}
}

//Makes printing of messages easy, if needed.  
/*func printMessage(message Message) {
	fmt.Println("This message is being sent from elevator with IP: ", message.SenderIP)
	fmt.Println("MessageType: ", message.MessageType)
	fmt.Println("Target IP: ", message.TargetIP)
	fmt.Println("Active: ", message.Elevator.Active)
	fmt.Println("Floor: ", message.Elevator.Floor)
	fmt.Println("Direction: ", message.Elevator.Direction)
	fmt.Println("NewlyInit: ", message.Elevator.NewlyInit)
	fmt.Println("Internal Orders: ", message.Elevator.InternalOrders)
	fmt.Println("External Up Orders: ", message.Elevator.ExternalUp)
	fmt.Println("External down orders: ", message.Elevator.ExternalDown)
	fmt.Println("Order: ", message.Order, "\n")
}*/
