package main

import (
	. "./definitions"
	. "./driver"
	. "./eventhandler"
	. "./network"
	"fmt"
)

func main() {

	//Initializes local elevator.
	if Init(){
		fmt.Println("The elevator was able to initialize")
	} else {
		fmt.Println("The elevator was not able to initialize")
	}

	floorChan := make(chan int)
	buttonChan := make(chan Order)
	messageReciveChan := make(chan Message)
	lightEventChan := make(chan int)
	newElevatorChan := make(chan string)
	deadElevatorChan := make(chan string)

	// Handels incoming and outgoing messages and what needs to be done when they get one.
	go MessageReciever(messageReciveChan)
	go MessageTypeHandler(messageReciveChan, floorChan, buttonChan, lightEventChan)
	go MessageBroadcast(MessageBroadcastChan)

	// Handels heartbeats, finds new elevators, tells us wether known elevators are dead or alive.
	go HeartbeatEventCheck(newElevatorChan, deadElevatorChan)
	go HeartbeatEventHandler(newElevatorChan, deadElevatorChan)
	go SendHeartBeat()

	//Tells other elevators if I need an update on my queue. 
	StatusUpdate()

	// Handels floor and button events: Does whatever needs to be done when buttons are pushed or floors are reached
	go ButtonandFloorEventHandler(floorChan, buttonChan, lightEventChan)
	go LampHandler(lightEventChan)




	keepElevatorGoingChan := make(chan string)
	<-keepElevatorGoingChan

}
