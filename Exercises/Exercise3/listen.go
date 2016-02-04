package main

import(
    "net"
    "fmt"
    "log"
    "strcov"
    "time"
    )

func getSenderID(sender *net.UDPAddr) ID {
    return ID(sender.IP.String())
}

func listen(socket *netUDPConn){
    for {
        buffer := make([]byte, 1024);
        read_bytes, sender, err := socket.ReadFromUDP(buffer);
        if err == nil{
            fmt.Println("Received: " + string(buffer[0:read_bytes]));
            fmt.Println(sender);
            fmt.Println(getSenderID(sender))
        }
        else{
            log.Println(err);
        }
    }
}

func CheckError(err error){
    if err != nil{
        fmt.Println("Error:", err)
    }

}

func main() {
    localIP := "78.91.21.240"; //HUSK!!! riktig IP
    localPort = 20005;
    localAddress, err := net.ResolveUDPAddr("udp", strconv.Itoa(localPort))
    

    for {

    }    

}