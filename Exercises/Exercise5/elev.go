package Driver

import(
) 

const( 
	MOTOR_SPEED = 2800
 	N_FLOORS = 3
 	)

//floor light
//button lights

var button_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func elev_init() {
	init_success := c.io_init()
	for (f := 0; f < N_FLOORS; f++){
		for (elev_button_type)
	}	
	// Set every set function to zero
	// Check initialization of hardware
}


func elev_set_motor_direction(direction int) {
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

func elev_set_button_lamp(floor, button int){
	// floor can be any N_FLOOR
	// button indicates UP, DOWN or COMMAND.
	if button == UP{
		if floor == 0 {
			io_set_bit(LIGHT_UP1)
		}
		if floor == 1 {
			io_set_bit(LIGHT_UP2)
		}
		if floor == 2 {
			io_set_bit(LIGHT_UP3)
		}
	}
	if button == DOWN {
		if floor == 1 {
			io_set_bit(LIGHT_DOWN2)
		}
		if floor == 2 {
			io_set_bit(LIGHT_DOWN3)
		}
		if floor == 3 {
			io_set_bit(LIGHT_DOWN4)
		}
		
	}
	if button == COMMAND {
		if floor == 0 {
			io_set_bit(COMMAND1)
		}
		if floor == 1 {
			io_set_bit(COMMAND2)
		}
		if floor == 2 {
			io_set_bit(COMMAND3)
		}
		if floor == 3 {
			io_set_bit(COMMAND4)
		}
		
	}
}

func elev_set_floor_indicator(floor int) {
	if (floor & 0x02) != 0 {
		Io_set_bit(LIGHT_FLOOR_IND1)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND1)
	}

	if (floor & 0x01) != 0 {
		Io_set_bit(LIGHT_FLOOR_IND2)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND2)
	}
}

func elev_set_door_open_lamp(door int) {
	if door == 1 {
		io_set_bit(LIGHT_DOOR_OPEN)
	} else{
		io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func elev_set_stop_lamp(stop int) {
	if stop == 1 {
		io_set_bit(LIGHT_STOP)
	} else{
		io_clear_bit(LIGHT_STOP)
	}
}

func elev_get_button_signal(button, floor int) int{
	if io_read_bit(button_channel_matrix[floor][button]) {
		return 1
	} else {
		return 0
	}
}
func Get_button_signal(button int, floor int) int {
	
}

func elev_get_stop_signal() int {
	return(io_read_bit(STOP))
}
func elev_get_obstruction_signal() int {
	return (io_read_bit(OBSTRUCTION))
}