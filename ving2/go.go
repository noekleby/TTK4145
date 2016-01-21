// Go 1.2

package main

import (
    "fmt"
    "runtime"
    "time"
)

var i int = 0

var channel chan int
var done1 chan bool
var done2 chan bool


func someGoroutine1(channel chan int, done1 chan bool) {
	time.Sleep(time.Second)
	for j := 0; j < 1000000; j++ {
		channel <- 1
		i++
		<- channel
	}
	done1 <- true

}
func someGoroutine2(channel chan int, done2 chan bool) {
	time.Sleep(time.Second)
	for j := 0; j < 1000000; j++ {
		channel <- 1
		i--
		<- channel
	}
	done2 <- true
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())    // I guess this is a hint to what GOMAXPROCS does...
                                            // Try doing the exercise both with and without it!

    channel := make(chan int, 1)
    done1 := make(chan bool, 1)
    done2 := make(chan bool, 1)
    
    go someGoroutine1(channel, done1)                  // This spawns someGoroutine() as a goroutine
    go someGoroutine2(channel, done2)
    // We have no way to wait for the completion of a goroutine (without additional syncronization of some sort)
    // We'll come back to using channels in Exercise 2. For now: Sleep.
    <-done1
    <-done2

    //i := <-channel
    fmt.Println(i)
}
