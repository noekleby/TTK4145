package network

import (
	."../definitions"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const (
	HeartBeatPort = 30115
	StatusPort    = 30215
)

var broadcastChan = make(chan Message)

func BroadcastMessage(message Message) {
	broadcastChan <- message
}

func GetIP() string {

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

func HeartbeatTransceiver(newElevatorChan chan string, deadElevatorChan chan string) {

	receive := make(chan []byte, 1)
	heartbeats := make(map[string]*time.Time)
	go udpRx(receive, HeartBeatPort)
	go sendHeartBeat()

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
				fmt.Println("Waring:", dur)
				deadElevatorChan <- i
				delete(heartbeats, i)
			}
		}
	}
}

func MessageTransceiver(receiveChan chan Message) {

	receive := make(chan []byte)

	go udpRx(receive, StatusPort)
	go sendStatus(broadcastChan)

	for {
		RxMessageBs := <-receive
		RxMessage := Message{}
		error := json.Unmarshal(RxMessageBs, &RxMessage)
		if error != nil {
			fmt.Println("error:", error)
		}
		if RxMessage.SenderIP != GetIP() {
			receiveChan <- RxMessage
		}
	}
}

/****************************************************************************************************************
Private
*/

func sendHeartBeat() {
	send := make(chan []byte, 1)
	go udpTx(send, HeartBeatPort)

	for {
		myBeat := Heartbeat{GetIP(), time.Now()}
		myBeatBs, error := json.Marshal(myBeat)

		if error != nil {
			fmt.Println("error:", error)
		}
		send <- myBeatBs
		time.Sleep(100 * time.Millisecond)
	}
}

func sendStatus(toSend chan Message) {
	send := make(chan []byte)
	go udpTx(send, StatusPort)

	for {
		temp := <-toSend
		toSendBs, error := json.Marshal(temp)
		if error != nil {
			fmt.Println("error:", error)
		}
		send <- toSendBs
	}
}

func udpDial(port int) *net.UDPConn {
	casting, error := net.ResolveUDPAddr("udp", fmt.Sprintf("129.241.187.255:%d", port))
	if error != nil {
		fmt.Println("error:", error)
	}
	socket, error := net.DialUDP("udp", nil, casting)
	if error != nil {
		fmt.Println("error:", error)
	}
	return socket
}

func udpListen(port int) *net.UDPConn {
	local, error := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if error != nil {
		fmt.Println("error:", error)
	}

	socket, error := net.ListenUDP("udp", local)
	if error != nil {
		fmt.Println("error:", error)
	}
	return socket
}

func udpRx(rx chan []byte, port int) {
	for {
		socket := udpListen(port)
		buffer := make([]byte, 1024)
		n, _, error := socket.ReadFromUDP(buffer)

		if error != nil {
			fmt.Println("error:", error)
		}

		buffer = buffer[:n]
		rx <- buffer
		socket.Close()
	}
}

func udpTx(tx chan []byte, port int) {
	for {
		socket := udpDial(port)
		dummy := <-tx
		socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
		_, error := socket.Write(dummy)
		if error != nil {
			fmt.Println("error:", error)
		}
		socket.Close()
	}
}