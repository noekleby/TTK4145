package main

import (
	. "./definitions"
	. "./driver"
	. "./eventhandler"
	. "./fsm"
	. "./network"
	"fmt"
)

func main() {

	//Initializes all elevators.
	if Init() == 1 {
		fmt.Println("The elevator was able to initialize")
	} else {
		fmt.Println("The elevator was not able to initialize")
	}
	InitFsm()

	floorChan := make(chan int)
	buttonChan := make(chan Order)
	messageReciveChan := make(chan Message)

	// Handels floor and button events: Does whatever needs to be done when buttons are pushed or floors are reached
	// This includes floorpanel from both local and other elevators, elevatorpanel from local elevator and floor events on local elevator.
	go ButtonandFloorEventHandler(floorChan, buttonChan)

	// Handels incoming and outgoing messages and what needs to be done when they get one.
	go MessageReciever(messageReciveChan)
	go MessageTypeHandler(messageReciveChan, floorChan, buttonChan)
	go MessageBroadcast(MessageBroadcastChan)

	// Handels heartbeats, finds new elevators, tells us wether known elevators are dead or alive.
	newElevatorChan := make(chan string)
	deadElevatorChan := make(chan string)
	go HeartbeatEventCheck(newElevatorChan, deadElevatorChan)
	go HeartbeatEventHandler(newElevatorChan, deadElevatorChan)
	go SendHeartBeat()

	KeepElevatorGoingChan := make(chan string)
	<-KeepElevatorGoingChan

}
