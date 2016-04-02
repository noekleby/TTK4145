package main 

import(
	."driver"
	"fmt"
	"time"
)


func main() {
	if (elev_init() == 0){
		fmt.Println("The elevator was not able to initialize")
	}
}