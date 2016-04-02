package driver

import(
	"fmt"
) 

const( 
	MOTOR_SPEED = 2800
 	N_FLOORS = 4
	N_BUTTONS = 3 // needs to be N_FLOORS-1 for the init function
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

var button=[N_floors][N_buttons]int{
	{0,0,0},
	{0,0,0},
	{0,0,0},
	{0,0,0},
}

var Sensors = [N_floors]int{SENSOR_FLOOR1,SENSOR_FLOOR2,SENSOR_FLOOR3,SENSOR_FLOOR4}

//Check initialization of hardware, drive down to IDLE, clear all lights except floor indicator. 
func elev_init(){
	init_success := c.io_init()
	
	if (init_success){
		for (f := 0; f <= N_FLOORS; f++){		//iterates over the 4 floors
			for (b := 0; b <= N_BUTTONS; b++){	//iterates over buttons for all other floors than the one you are in
				elev_set_button_lamp(b, f, false)//clears every button lamp (false = light off)
			}
		}	
	}
	else{
		fmt.println("Unable to initialize elevator hardware!")
	}
	Set_stop_lamp(0)
	Set_door_open_lamp(0)

	//Finne ut hvor vi er, dersom vi ikke er i 1. etasje skal heisen forflytte seg ned dit. 
	// Så må vi sette på etasjeindikatoren. 
	Set_floor_indicator(0)
	floor = elevGetFloorSignal()
	if (floor =! 0){
		elevSetMotorDirection(-1)
	}
	for Elev_get_floor_sensor_signal() == -1 {}
  	Elev_set_motor_direction(DIR_STOP)
  	Elev_set_floor_indicator(Elev_get_floor_sensor_signal())

    return 1
}


func elevSetMotorDirection(direction int) {
	if (direction == 0){
        io_write_analog(MOTOR, 0)
    }
    if (direction > 0) {
        io_clear_bit(MOTORDIR)
        io_write_analog(MOTOR, MOTOR_SPEED)
    }
    if (direction < 0) {
        io_set_bit(MOTORDIR)
        io_write_analog(MOTOR, MOTOR_SPEED)
    }
}

func elevSetButtonLamp(floor int, button int, value bool){
	// floor can be any N_FLOOR
	// button indicates UP (= 1), DOWN (=-1) or COMMAND (=0)
	// value sets the light on/off
	if ((N_FLOORS > floor >= 0) && N_BUTTONS > button >= 0) {
    		if (value) {
        		io_set_bit(lamp_channel_matrix[floor][button])
    		} else {
        		io_clear_bit(lamp_channel_matrix[floor][button])
    		}
	}
	else {
		fmt.println("ERROR: Unable to update the button lamps")
	}
}

func elevSetFloorIndicator(floor int) {
	if (floor & 0x02) != 0 { // handles the odd numbered floors
		Io_set_bit(LIGHT_FLOOR_IND1)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND1)
	}

	if (floor & 0x01) != 0 { // handles the even numbered floors
		Io_set_bit(LIGHT_FLOOR_IND2)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND2)
	}
}

func elevSetDoorOpenLamp(door int) {
	if door == 1 {
		io_set_bit(LIGHT_DOOR_OPEN)
	} else{
		io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func elevSetStopLamp(stop int) {
	if stop == 1 {
		io_set_bit(LIGHT_STOP)
	} else{
		io_clear_bit(LIGHT_STOP)
	}
}

func elevGetButtonSignal(button int, floor int) int {
	if ((N_FLOORS > floor >= 0) && N_BUTTONS > button >= 0){ // checks if floor and button are valid
    		if (io_read_bit(button_channel_matrix[floor][button])) { // what's the purpose of read_bit(?)
        		return 1
    		} else {
        		return 0
    		}  
	} else{
		fmt.println("ERROR: Unable to read the button signal!")
	}
}

func elevGetStopSignal() int {
	return(io_read_bit(STOP))
}
func elevGetObstructionSignal() int {
	return (io_read_bit(OBSTRUCTION))
}
func elevGetFloorSignal() {
	for (f := 0; f <= N_FLOORS; f++){
		if io_read_bit(Sensors[f]) == 1 {
			return f
		}
	}
	return -1
}