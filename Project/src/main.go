package main

import (
	."./driver"
	."./eventhandler"
	"./fsm"
	//"./queue"
	"fmt"
	."./network"
	//"./transManager.go"
	//"time"
)

func main() {

	//Initzialization of elevator Hardware and fsm.
	if Init() == 0 {
		fmt.Println("The elevator was not able to initialize")
	}
	if Init() == 1 {
		fmt.Println("The elevator was able to initialize")
	}
	fsm.InitFsm()


	go ButtonandFloorEventHandler()

	newElevatorChan := make(chan string)
	deadElevatorChan := make(chan string)
	go HeartbeatEventCheck(newElevatorChan, deadElevatorChan)
	go HeartbeatEventHandler(newElevatorChan, deadElevatorChan)
	go SendHeartBeat()

	alwaysOnChan := make(chan string)
	<- alwaysOnChan
	
}
