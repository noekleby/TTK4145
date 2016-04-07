package driver
/*
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

var button = [N_FLOORS][N_BUTTONS]int{
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
}
var sensors = [N_FLOORS]int{SENSOR_FLOOR1, SENSOR_FLOOR2, SENSOR_FLOOR3, SENSOR_FLOOR4}

//Check initialization of hardware, drive down to IDLE, clear all lights except floor indicator.
func Init() int{
	init_success := ioInit()


	if init_success == 1 {
		RunStop()
		ElevateBottomFloor()
		elevSetStopLamp(0)
		elevSetDoorOpenLamp(0)
		for f := 0; f <= N_FLOORS; f++ { //iterates over the 4 floors
			for b := 0; b <= N_BUTTONS; b++ { //iterates over buttons for all other floors than the one you are in
				SetButtonLamp(b, f, false) //clears every button lamp (false = light off)
			}
		}
	} else {
		fmt.Println("Unable to initialize elevator hardware!")
	}
	return init_success

}



func GetFloorSignal() int {
	for f := 0; f <= N_FLOORS; f++ {
		if ioReadBit(sensors[f]) == 1 {
			return f 
		}
	}
	return -1 // instead of else return 
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
		for ioReadBit(sensors[0]) == 0 {
			SetFloorIndicator(GetFloorSignal())
			time.Sleep(200*time.Millisecond)
		}
		SetFloorIndicator(GetFloorSignal())
		StopElevate()
	}
}

func SetFloorIndicator() bool {
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
	if floor < 0 || floor >= N_FLOORS{
		return false 
	}
	return true 
}

func SetButtonLamp(floor int, button int, value bool) {
	// floor can be any N_FLOOR
	// button indicates UP (= 1), DOWN (=-1) or COMMAND (=0)
	// value sets the light on/off
	if (floor >= 0) && (button >= 0) {
		if value {
			ioSetBit(lamp_channel_matrix[floor][button])
		} else {
			ioClearBit(lamp_channel_matrix[floor][button])
		}
	} else {
		fmt.Println("ERROR: Unable to update the button lamps")
	}
}


func elevSetDoorOpenLamp(door int) {
	if door == 1 {
		ioSetBit(LIGHT_DOOR_OPEN)
	} else {
		ioClearBit(LIGHT_DOOR_OPEN)
	}
}

func elevSetStopLamp(stop int) {
	if stop == 1 {
		ioSetBit(LIGHT_STOP)
	} else {
		ioClearBit(LIGHT_STOP)
	}
}

func elevGetButtonSignal(button int, floor int) int {
	if (floor >= 0) && button >= 0 { // checks if floor and button are valid
		if ioReadBit(button_channel_matrix[floor][button]) == 1 { // what's the purpose of read_bit(?)
			return 1
		} else {
			return 0
		}
	} else {
		fmt.Println("ERROR: Unable to read the button signal!")
	}
	return -1
}

func elevGetStopSignal() int {
	return (ioReadBit(STOP))
}
func elevGetObstructionSignal() int {
	return (ioReadBit(OBSTRUCTION))
}

/*func ElevateTopFloor() {
	if GetFloorSignal() != 3 {
		ElevateUp()
		for ioReadBit(sensors[3]) == 0 {
			SetFloorIndicator(GetFloorSignal())
			time.Sleep(200*time.Millisecond)
		}
		SetFloorIndicator(GetFloorSignal())
		StopElevate()
	}
}*/


