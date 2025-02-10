package elev_actuator

import (
	eio "Driver-go/elevio"
	"time"
)

func Open_Door() {
	eio.SetMotorDirection(eio.MD_Stop)
	eio.SetDoorOpenLamp(true)
	time.Sleep(2*time.Second)
	eio.SetDoorOpenLamp(false)
}