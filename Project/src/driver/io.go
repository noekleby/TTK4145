package driver

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"

func ioInit() bool {
	if int(C.io_init()) == 1 {
		return true
	} else{
		return false
	}
}

func ioSetBit(channel int) {
	C.io_set_bit(C.int(channel))
}

func ioClearBit(channel int) {
	C.io_clear_bit(C.int(channel))
}

func ioReadBit(channel int) int {
	return int(C.io_read_bit(C.int(channel)))
}

func ioReadAnalog(channel int) int {
	return int(C.io_read_analog(C.int(channel)))
}

func ioWriteAnalog(channel, value int) {
	C.io_write_analog(C.int(channel), C.int(value))
}
