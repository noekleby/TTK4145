// All UDP functions
package network

/*import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

type ID string

func getSenderID(sender *net.UDPAddr) ID {
	return ID(sender.IP.String())
}

func listen(socket *net.UDPConn) {
	for {
		buffer := make([]byte, 1024)
		read_bytes, _, err := socket.ReadFromUDP(buffer)
		if err == nil {
			fmt.Println("Received: " + string(buffer[0:read_bytes]))
		} else {
			log.Println(err)
		}
	}
}

func transmit(socket *net.UDPConn) {
	for {
		time.Sleep(2000 * time.Millisecond)

		message := "From server: Hello I got your mesage"
		socket.Write([]byte(message))
		fmt.Println("Sendte: " + message)
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func main() {
	broadcastIP := "129.241.187.255" //HUSK!!! riktig IP
	localPort := 30000
	localAddress, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(localPort))

	serverPort := 20005
	serverAddress, _ := net.ResolveUDPAddr("udp", broadcastIP+":"+strconv.Itoa(serverPort))

	fmt.Println("Local adress:", localAddress)
	fmt.Println("Server adress", serverAddress)

	listenSocket, _ := net.ListenUDP("udp", localAddress)
	transmitSocket, _ := net.DialUDP("udp", nil, serverAddress)

	listen(listenSocket)
	transmit(transmitSocket)
}*/