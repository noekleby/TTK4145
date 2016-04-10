package driver

import (
	"fmt"
	"time"
)

const (
	MOTOR_SPEED = 2800
	N_FLOORS    = 4
	N_BUTTONS   = 3
)

var button_channel_matrix = [N_FLOORS][N_BUTTONS]int{ //button command for 4 floors
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}
var lamp_channel_matrix = [N_FLOORS][N_BUTTONS]int{ //floor lights for 4 floors
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var sensors = [N_FLOORS]int{SENSOR_FLOOR1, SENSOR_FLOOR2, SENSOR_FLOOR3, SENSOR_FLOOR4}

//Check initialization of hardware, drive down to IDLE, clear all lights except floor indicator.
func Init() int {
	init_success := ioInit()

	if init_success == 1 {
		StopElevate()
		ElevateBottomFloor()
		fmt.Println("I'm currently initializing elevator and hardware")
		ElevSetStopLamp(0)
		ElevSetDoorOpenLamp(0)
		for f := 0; f < N_FLOORS; f++ { //iterates over the 4 floors
			for b := 0; b < N_BUTTONS; b++ { //iterates over buttons for all other floors than the one you are in
				SetButtonLamp(f, b, false) //clears every button lamp (false = light off)
			}
		}
	} else {
		fmt.Println("Unable to initialize elevator hardware!")
	}
	return init_success

}

func ElevateDown() {
	ioSetBit(MOTORDIR)
	ioWriteAnalog(MOTOR, MOTOR_SPEED)
}

func ElevateUp() {
	ioClearBit(MOTORDIR)
	ioWriteAnalog(MOTOR, MOTOR_SPEED)
}

func StopElevate() {
	ioWriteAnalog(MOTOR, 0)
}

func ElevateBottomFloor() {
	if GetFloorSignal() != 0 {
		ElevateDown()
		fmt.Println("I'm currently driving")
		for ioReadBit(sensors[0]) == 0 {
			time.Sleep(time.Millisecond * 200)
		}
		SetFloorIndicator(GetFloorSignal())
		StopElevate()
	}
}

func GetFloorSignal() int {
	for f := 0; f < N_FLOORS; f++ {
		if ioReadBit(sensors[f]) == 1 {
			return f
		}
	}
	return -1
}

func SetFloorIndicator(floor int) bool {
	if (floor & 0x02) != 0 { // handles the odd numbered floors
		ioSetBit(LIGHT_FLOOR_IND1)
	} else {
		ioClearBit(LIGHT_FLOOR_IND1)
	}
	if (floor & 0x01) != 0 { // handles the even numbered floors
		ioSetBit(LIGHT_FLOOR_IND2)
	} else {
		ioClearBit(LIGHT_FLOOR_IND2)
	}
	if floor < 0 || floor >= N_FLOORS {
		return false
	}
	return true
}

func SetButtonLamp(floor int, button int, value bool) {
	if (floor >= 0) && (button >= 0) {
		if value {
			ioSetBit(lamp_channel_matrix[floor][button])
		} else {
			ioClearBit(lamp_channel_matrix[floor][button])
		}
	} else {
		fmt.Println("ERROR: Unable to update button lamps")
	}
}

func ElevSetDoorOpenLamp(door int) {
	if door == 1 {
		ioSetBit(LIGHT_DOOR_OPEN)
	} else {
		ioClearBit(LIGHT_DOOR_OPEN)
	}
}

func ElevSetStopLamp(stop int) {
	if stop == 1 {
		ioSetBit(LIGHT_STOP)
	} else {
		ioClearBit(LIGHT_STOP)
	}
}

func ElevGetButtonSignal(button int, floor int) int {
	if floor < 0 || floor >= N_FLOORS || button < 0 || button >= N_BUTTONS {
		return 0
	} else {
		return ioReadBit(button_channel_matrix[floor][button])
	}
}

func ElevGetStopSignal() int {
	return (ioReadBit(STOP))
}
func ElevGetObstructionSignal() int {
	return (ioReadBit(OBSTRUCTION))
}

func GetDirection() int {
	return ioReadBit(MOTORDIR)
}
