package main

import (
	"./driver"
	"fmt"
	//"time"
)

func main() {
	if driver.elev_init() == 0 {
		fmt.Println("The elevator was not able to initialize")
	}
	if driver.elev_init() == 1 {
		fmt.Println("The elevator was able to initialize")
	}
}

// cannot refer to unexported name driver.elev_init
// ./main.go:10: undefined: driver.elev_init  :(
