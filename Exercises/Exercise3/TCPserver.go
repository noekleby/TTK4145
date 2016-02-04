package main;

import (
    "fmt"
    "net"
    "os"
)

func main() {
    listener, err := net.Listen("tcp", "localhost:3000") //Listen to the tcp port
    if err!= nil {
        fmt.Println("Error listening", err.Error())
        os.Exit(1)
    }
    for{
        if conn, err := listener.Accept(); err == nil{
            go handle(conn) //Kjøres hvis vi ikke har noen errors
        }
        else{
            continue
        }
    }
}

func handle(conn net.Conn){
    fmt.Println("Connection established")
    defer conn.close()
    data := make([]byte, 1024) //Lager en slice av størrelse 1024
    n,err:=conn.Read(data) //legger tilkoblingen i bufferen
    if err!= nil {
        fmt.Println("Error reading", err.Error())
    }
    fmt.Println(string(data))
    conn.Write([]byte("Data received"))
    conn.Close()
}
