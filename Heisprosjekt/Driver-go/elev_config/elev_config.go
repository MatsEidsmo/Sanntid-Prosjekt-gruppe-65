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
	RequestMatrix 	[][] int
	Behaviour 		ElevatorBehavior
	Obstruction 	bool
	ElevID			string
	TimeToIdle 		int
}



func InitElev(id string) Elevator{
	rm := make([][] int, 0)
	for floor := 0; floor < N_floors; floor++ {
		rm = append(rm, make([]int, N_buttons))
		for btn := range rm[floor]{
			rm[floor][btn] = 0
		}
	}
	return Elevator{
		Floor: 			0,
		Dir:        	eio.MD_Stop,
		RequestMatrix:  rm,
		ElevID: 		id,	
		TimeToIdle: 	0,
	}

	
}