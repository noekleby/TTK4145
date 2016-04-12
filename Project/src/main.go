package main

import (
	"./driver"
	"./eventhandler"
	//"./fsm"
	//"./queue"
	"fmt"
	"./network"
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
	go network.HeartbeatTransceiver(newElevatorChan, deadElevatorChan)
	go eventhandler.HeartbeatEventHandler(newElevatorChan, deadElevatorChan)

	//Starting gorutines to check for events on buttons and floor sensors
	eventhandler.ButtonandFloorEventHandler()
	
}
