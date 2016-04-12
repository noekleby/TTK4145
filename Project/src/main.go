package main

import (
	"./driver"
	."./eventhandler"
	//"./fsm"
	//"./queue"
	"fmt"
	."./network"
	//"./transManager.go"
	//"time"
)

func main() {

	//Initzialization of elevator Hardware and fsm.
	if driver.Init() == 0 {
		fmt.Println("The elevator was not able to initialize")
	}
	if driver.Init() == 1 {
		fmt.Println("The elevator was able to initialize")
	}


	//putt dette i main: 
	newElevatorChan := make(chan string)
	deadElevatorChan := make(chan string)
	
	
	//go network.HeartbeatEventCheck(newElevatorChan, deadElevatorChan)
	go HeartbeatEventCheck(newElevatorChan, deadElevatorChan)
	go HeartbeatEventHandler(newElevatorChan, deadElevatorChan)
	go SendHeartBeat()

	//Starting gorutines to check for events on buttons and floor sensors
	ButtonandFloorEventHandler()
	
}
