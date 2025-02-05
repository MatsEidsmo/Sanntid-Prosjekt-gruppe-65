package requests

import (
	"Driver-go/elevio"
	ec "Driver-go/elev_config"
)

func requests_above(e *ec.Elevator) int {
	for f := e.Floor+1; f < ec.N_floors; f++ {
		for btn := 0; btn < ec.N_buttons; btn++ {
			if e.RequestMatrix[f][btn] == 1 {
				return 1
			}
		}
	}
	return 0
}

func requests_below(e *ec.Elevator) int {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < ec.N_buttons; btn++ {
			if e.RequestMatrix[f][btn] == 1 {
				return 1
			}
		}
	}
	return 0
}

func requests_here(e *ec.Elevator) int {
	for btn := 0; btn < ec.N_buttons; btn++ {
		if e.RequestMatrix[e.Floor][btn] == 1 {
			return 1
		}
	}
	return 0
}

func add_order(e *ec.Elevator, floor int, btnType elevio.ButtonType) {
	e.RequestMatrix[floor][btnType] = 1
}

func choose_dir(e *ec.Elevator) {
	if e.Dir == elevio.MD_Up {
		
	}
}