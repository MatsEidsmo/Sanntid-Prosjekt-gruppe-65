package elev_actuator

import (
	eio "Driver-go/elevio"
	"time"
	ec "Driver-go/elev_config"
	el "Driver-go/elev_logic"
)

func Open_Door(door_timer chan bool, e *ec.Elevator) {
	timer := time.NewTimer(3*time.Second)

	go func() {
		eio.SetMotorDirection(eio.MD_Stop)
		eio.SetDoorOpenLamp(true)
		<-timer.C
		eio.SetDoorOpenLamp(false)
		el.Clear_Floor_Requests(e)
		door_timer <- true
	}()


	// eio.SetMotorDirection(eio.MD_Stop)
	// eio.SetDoorOpenLamp(true)
	// time.Sleep(2*time.Second)
	// eio.SetDoorOpenLamp(false)
}