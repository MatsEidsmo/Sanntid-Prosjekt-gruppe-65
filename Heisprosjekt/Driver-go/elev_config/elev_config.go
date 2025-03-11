package elev_config

import (
	eio "Driver-go/elevio"
	"time"
)

const N_floors int = 4
const N_buttons int = 3


const DOOR_TIMEOUT = 3 * time.Second

type ElevatorBehavior int
const (
	EB_Moving	ElevatorBehavior = iota -1
	EB_DoorOpen 
	EB_Idle
)

type Elevator struct {
	Floor 			int
	Dir 			eio.MotorDirection // except MD_Stop
	RequestMatrix 	[N_floors][N_buttons] int
	Behaviour 		ElevatorBehavior
}