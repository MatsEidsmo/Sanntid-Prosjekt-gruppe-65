package elev_actuator

import (
	ec "Driver-go/elev_config"
	el "Driver-go/elev_logic"
	eio "Driver-go/elevio"
	"time"
	"fmt"
)

var DoorTimer *time.Timer

func Timer_init() {
	DoorTimer = time.NewTimer(ec.DOOR_TIMEOUT)
	if !DoorTimer.Stop() {
		<- DoorTimer.C
	}
}

func Timer_start() {

	DoorTimer.Stop()
	select {
	case <- DoorTimer.C:
	default:
	}
	DoorTimer.Reset(ec.DOOR_TIMEOUT)
}



func Open_Door(e *ec.Elevator) {
	
	Timer_start()
	eio.SetMotorDirection(eio.MD_Stop)
	eio.SetDoorOpenLamp(true)
	fmt.Println("Door is Open!")
	el.Clear_Floor_Requests(e)
	


	// eio.SetMotorDirection(eio.MD_Stop)
	// eio.SetDoorOpenLamp(true)
	// time.Sleep(2*time.Second)
	// eio.SetDoorOpenLamp(false)
}