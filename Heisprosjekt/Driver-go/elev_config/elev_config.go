package elev_config

import (
	eio "Driver-go/elevio"
	"time"
)

const N_floors int = 4
const N_buttons int = 3

const TRAVEL_TIME = 3000 * time.Millisecond
const DOOR_TIMEOUT = 3000 * time.Millisecond

type ElevatorBehavior int
const (
	EB_Moving	ElevatorBehavior = iota -1
	EB_DoorOpen 
	EB_Idle
)

type Elevator struct {
	Floor 			int
	Dir 			eio.MotorDirection 
	RequestMatrix 	[N_floors][N_buttons] int
	Behaviour 		ElevatorBehavior
	Obstruction 	bool
	ElevID			string
}



func DefaultElev(e *Elevator) {
	e.Floor 	= 		0
	e.Behaviour = EB_Idle
	
}