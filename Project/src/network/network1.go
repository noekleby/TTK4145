package network

import (
	"fmt"
	"net"
	"time"
	"log" // to log errors (time and writes to standard errors)
	"encoding/json"
	"../definitions"
	"strconv"
)

type Heartbeat struct {
	ID string
	Tid string

}
var broadcastChan = make(chan definitions.Message)
var port = definitions.HeartBeatPort

func BroadcastMessage(message definitions.Message) {
	broadcastChan <- message
}
//Checks if any of the following events happens: one elevator dies or one elevator awakes. 
func HeartbeatEventCheck(newElevatorChan chan string, deadElevatorChan chan string){

	//Time time.Time // Cannot marshall this typ
	
	bufferHeartbeate := make(chan []byte, 1)
	//storedElevators := make(map[string]*time.Time)

	go UdpRx(bufferHeartbeate, port)
	go SendHeartbeat()

	for {

		RecievedHeartbeat := <- bufferHeartbeate
		//var newHeartbeat Heartbeat
		var newHeartbeat string

		error := json.Unmarshal(RecievedHeartbeat, &newHeartbeat)
		if error != nil {
			//fmt.Println("This unmarshal", newHeartbeat.ID)
			fmt.Println("error:", error)
		}
		if error == nil {
			fmt.Println("Its working", newHeartbeat)
		}
		/*
		_,exist := storedElevators[newHeartbeat.ID]

		if exist {

			t1, _ := time.Parse(time.RFC3339, newHeartbeat.Tid)
			storedElevators[newHeartbeat.ID] = &t1

		} else {

			newElevatorChan <- newHeartbeat.ID
			t1, _ := time.Parse(time.RFC3339, newHeartbeat.Tid)
			storedElevators[newHeartbeat.ID] = &t1
			fmt.Println(t1)
		}
		for id, t := range storedElevators{
			fmt.Println(storedElevators[id])
			duration := time.Since(*t)
			fmt.Println("This is duration: \n", duration)
			fmt.Println(duration.Seconds())
			if duration.Seconds() > 2 {
				deadElevatorChan <- id 
				delete(storedElevators, id)
			}
		}*/

	}

}

func UdpRx(rx chan []byte, port int) {
	for {
		socket := GetListenSocket(port)  
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

func GetLocalIP() string { // hjelpefunksjon fra stack overflow
	addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback (localhost 127.0.0.1) then display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

func SendHeartbeat() {
	send := make(chan [] byte,1)
	go Transmit(GetTransmitSocket(), send)

	for {
		//t := time.Now()
		//str := fmt.Sprintf("%q",t.Format(time.RFC3339))
		//str1 := GetLocalIP()
		localBeat := Heartbeat{"hello", "hello"}
		buffer, err := json.Marshal(localBeat)

		if err != nil {
			fmt.Println("error:", err)
		}
		send <- buffer
		time.Sleep(1000 * time.Millisecond)
	}
}

func SendStatus(toSend chan definitions.Message, port int) {
	send := make(chan [] byte)
	go Transmit(GetTransmitSocket(port), send) 

	for {
		temp := <-toSend
		buffer, err := json.Marshal(temp)
		if err != nil {
			fmt.Println("error:", err)
		}
		send <- buffer
	}
}

func Listen(socket *net.UDPConn, storedChan chan definitions.Message) {
	for {
		buffer := make([]byte, definitions.MsgSize)
		length,_,err := socket.ReadFromUDP(buffer) 
		if err == nil {
			buffer = buffer[:length] // To just get the length of the message
			fmt.Println("A message was received.", string(buffer))
			var storedData definitions.Message
			err = json.Unmarshal(buffer, &storedData) // data from buffer will be stored in storedData
			if (err != nil){ // kan eventuelt lage en error function
				fmt.Println("Could not decode message.")
				log.Println(err)
			}
			storedChan <- storedData 
		}else {
			log.Println(err)
		}
		time.Sleep(400*time.Millisecond)
	}
}

func GetListenSocket (port int) *net.UDPConn {
	localAddress, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("There is an error in resolving.")
	} 
	listenSocket, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		fmt.Println("There is an error in listening.")
		defer listenSocket.Close() // defer utsetter kallet til de andre funksjonene har kjørt (trengs kanskje ikke her)
	}
	return listenSocket
}

func Transmit(socket *net.UDPConn, sendMsg chan [] byte) {
	for {
		temp := <- sendMsg
		buffer, err := json.Marshal(temp)
		if (err != nil){ // kan eventuelt lage en egen check error function
			fmt.Println("Could not encode message.")
			log.Println(err)
		}
		socket.Write([]byte(buffer))
		time.Sleep(2*time.Second)
	}
}


func GetTransmitSocket (port int) *net.UDPConn {
	serverAddress, err := net.ResolveUDPAddr("udp", GetLocalIP()+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("There is an error in resolving.")
	} 
	transmitSocket, _ := net.DialUDP("udp", nil, serverAddress)
	if err != nil {
		fmt.Println("There is an error in dialing.")
		defer transmitSocket.Close()
	}
	return transmitSocket
}

// Bruke marshall og unmarshall for encoding og decoding. Trengs for konvertere data til og fra byte-nivå og 
// tekstrepresentasjon.
// JSON er JavaScript Object Notation. Syntaks for å lagre og utveksle data. Lettere å bruke enn XML. 
// Sjekk blog.golang json and go (google it!)
