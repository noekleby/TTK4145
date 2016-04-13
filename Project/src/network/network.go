package network

import (
	. "../definitions"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	HeartBeatPort = 30113
	BroadcastPort = 30215
)

func BroadcastMessage(message Message) {
	MessageBroadcastChan <- message
}

func MessageReciever(messageRecieveChan chan Message) {

	bufferRecieveChan := make(chan []byte)
	go udpRecieve(bufferRecieveChan, BroadcastPort)

	for {
		buffer := <-bufferRecieveChan //Can not call it buffer
		msg := Message{}              //
		error := json.Unmarshal(buffer, &msg)
		if error != nil {
			fmt.Println("Error:", error)
			time.Sleep((4 * time.Second))
		}
		if msg.SenderIP != GetLocalIP() {
			messageRecieveChan <- msg
		}
	}
}
func MessageBroadcast(MessageBroadcastChan chan Message) {

	bufferSendChan := make(chan []byte)
	go udpSend(bufferSendChan, BroadcastPort)

	for {
		msg := <-MessageBroadcastChan
		buffer, error := json.Marshal(msg) //Can not call it buffer
		if error != nil {
			fmt.Println("Error:", error)
			time.Sleep((4 * time.Second))
		}
		fmt.Println("Broadcasting message")
		bufferSendChan <- buffer
	}
}

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

func HeartbeatEventCheck(newElevatorChan chan string, deadElevatorChan chan string) {

	receive := make(chan []byte, 1)
	heartbeats := make(map[string]*time.Time)
	go udpRecieve(receive, HeartBeatPort)

	for {
		otherBeatBs := <-receive
		otherBeat := Heartbeat{}
		error := json.Unmarshal(otherBeatBs, &otherBeat)
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
	}
}

func SendHeartBeat() {
	send := make(chan []byte, 1)
	go udpSend(send, HeartBeatPort)

	for {
		myBeat := Heartbeat{GetLocalIP(), time.Now()}
		myBeatBs, error := json.Marshal(myBeat)

		if error != nil {
			fmt.Println("error:", error)
		}
		send <- myBeatBs
		time.Sleep(100 * time.Millisecond)
	}
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func getTransmitSocket(port int) *net.UDPConn {

	serverAddress, err := net.ResolveUDPAddr("udp", fmt.Sprintf("129.241.187.255:%d", port))
	if err != nil {
		fmt.Println("There is an error in resolving server:", err)
	}
	transmitSocket, _ := net.DialUDP("udp", nil, serverAddress)
	if err != nil {
		fmt.Println("There is an error in dialing:", err)
	}
	return transmitSocket
}

func getListenSocket(port int) *net.UDPConn {
	localAddress, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("There is an error in resolving local:", err)
	}
	listenSocket, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		fmt.Println("There is an error in listening:", err)
	}
	return listenSocket
}

func udpRecieve(msg chan []byte, port int) {
	for {
		socket := getListenSocket(port)
		buffer := make([]byte, 1024)
		n, _, err := socket.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("error:", err)
		}
		buffer = buffer[:n]
		msg <- buffer
		socket.Close()
	}
}

func udpSend(msg chan []byte, port int) {
	for {
		socket := getTransmitSocket(port)
		temp := <-msg
		socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
		_, error := socket.Write(temp)
		if error != nil {
			fmt.Println("error:", error)
		}
		socket.Close()
	}
}
