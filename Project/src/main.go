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
	go ButtonandFloorEventHandler(floorChan, buttonChan)

	// Handels incoming and outgoing messages
	go MessageReciever(messageReciveChan)
	go MessageTypeHandler(messageReciveChan, floorChan, buttonChan)
	go MessageBroadcast(MessageBroadcastChan)

	// Handels heartbeats, finds new elevator and tell us continuly wether they ar alive or dead.
	newElevatorChan := make(chan string)
	deadElevatorChan := make(chan string)
	go HeartbeatEventCheck(newElevatorChan, deadElevatorChan)
	go HeartbeatEventHandler(newElevatorChan, deadElevatorChan)
	go SendHeartBeat()

	KeepElevatorGoingChan := make(chan string)
	<-KeepElevatorGoingChan

}
