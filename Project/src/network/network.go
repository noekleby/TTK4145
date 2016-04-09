//send, recieve, error
package network

import (
	"fmt"
	"net"
	"time"
	"log" // to log errors (time and writes to standard errors)
	"encoding/json"
)

const (
	//broadcastIP := "129.241.187.255"
	localPort := 30000
	serverPort := 20005
	msgSize := 1024
)

type msg struct(	//used for decoding the message (unmarshalling)
	destination int
	currFloor int
	// sende kø?
)

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

func Listen(socket *net.UDPConn) {
	var storedChan := make(chan msg)
	for {
		buffer := make([]byte, msgSize)
		length,_,err := socket.ReadFromUDP(buffer) 
		if err == nil {
			buffer = buffer[:length] // To just get the length of the message
			fmt.Println("A message was received.", string(buffer))
			var storedData msg
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

func Transmit(socket *net.UDPConn, sendMsg chan msg) {
	for {
		buffer, err := json.Marshal(sendMsg)
		if (err != nil){ // kan eventuelt lage en egen check error function
			fmt.Println("Could not encode message.")
			log.Println(err)
		}
		socket.Write([]byte(buffer))
		time.Sleep(2*time.Second)
	}
}

/*func PrintMessages (storedChan chan msg) {
	for {
		msg := <- storedChan
		fmt.Println(msg)
	}
}*/

func Init(sendMsg chan msg){
	runtime.GOMAXPROCS(runtime.NumCPU()) //sets the number of cpu cores the program can use simultaneously.
	//sets it here to Numcpu which is the number of cores available. 

	localAddress, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(localPort))
	if err != nil {
		fmt.Println("There is an error in resolving.")
	} 
	listenSocket, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		fmt.Println("There is an error in listening.")
		defer listenSocket.Close() // defer utsetter kallet til de andre funksjonene har kjørt (trengs kanskje ikke her)
	} 

	serverAddress, err := net.ResolveUDPAddr("udp", GetLocalIP()+":"+strconv.Itoa(serverPort))
	if err != nil {
		fmt.Println("There is an error in resolving.")
	} 
	transmitSocket, _ := net.DialUDP("udp", nil, serverAddress)
	if err != nil {
		fmt.Println("There is an error in dialing.")
		defer transmitSocket.Close()
	} 

	go listen(listenSocket)
	go transmit(transmitSocket, sendMsg)
	time.Sleep(time.Second*5)
}

// Bruke marshall og unmarshall for encoding og decoding. Trengs for konvertere data til og fra byte-nivå og 
// tekstrepresentasjon.
// JSON er JavaScript Object Notation. Syntaks for å lagre og utveksle data. Lettere å bruke enn XML. 
// Sjekk blog.golang json and go (google it!)