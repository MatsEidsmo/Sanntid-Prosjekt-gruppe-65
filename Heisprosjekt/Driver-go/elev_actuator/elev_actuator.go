package elev_actuator

import (
	ec "Driver-go/elev_config"
	el "Driver-go/elev_logic"
	eio "Driver-go/elevio"
	//"Driver-go/orders"
	"fmt"
	"time"
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
	e.Behaviour = ec.EB_DoorOpen
	fmt.Println("Door is Open!")
	
	
}

//Checks if obstruction is active, if not closes the door and chooses direction.
func Upon_Door_Timeout(e *ec.Elevator) {
	

	if e.Obstruction {
		println("Obstruction true")
		Open_Door(e)
		return
	}

	
	curr_dir := el.Choose_Dir(e)
	eio.SetDoorOpenLamp(false)
	el.Clear_Floor_Requests(e, false)
	if curr_dir == eio.MD_Stop {
		e.Behaviour = ec.EB_Idle
	} else {
		e.Behaviour = ec.EB_Moving
	}

	eio.SetMotorDirection(curr_dir)
}