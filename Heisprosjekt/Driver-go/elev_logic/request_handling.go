package requests

import (
	eio "Driver-go/elevio"
	ec "Driver-go/elev_config"
)

func requests_above(e *ec.Elevator) bool {
	for f := e.Floor+1; f < ec.N_floors; f++ {
		for btn := 0; btn < ec.N_buttons; btn++ {
			if e.RequestMatrix[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func requests_below(e *ec.Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < ec.N_buttons; btn++ {
			if e.RequestMatrix[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func requests_here(e *ec.Elevator) int {
	for btn := 0; btn < ec.N_buttons; btn++ {
		if e.RequestMatrix[e.Floor][btn] == 1 {
			return 1
		}
	}
	return 0
}

func Add_Request(e *ec.Elevator, floor int, btnType eio.ButtonType) {
	e.RequestMatrix[floor][btnType] = 1
}

func Remove_Request(e *ec.Elevator) {
	for btn := 0; btn < ec.N_buttons; btn++{
		e.RequestMatrix[e.Floor][btn] = 0
	}
}

func Choose_Dir(e *ec.Elevator) eio.MotorDirection {
	switch e.Dir {
	case eio.MD_Up:
		if requests_above(e) {
			return eio.MD_Up
		}
		if requests_below(e) {
			return eio.MD_Down
		}
	case eio.MD_Down:
		if requests_below(e) {
			return eio.MD_Down
		}
		if requests_above(e) {
			return eio.MD_Up
		}
	}
	return eio.MD_Stop
}

func Stop_Here(e *ec.Elevator) bool {
	if e.RequestMatrix[e.Floor][eio.BT_Cab] == 1 {
		return true
	}
	if e.RequestMatrix[e.Floor][eio.BT_HallUp] == 1 && e.Dir == eio.MD_Up{
		return true
	}
	if e.RequestMatrix[e.Floor][eio.BT_HallDown] == 1 && e.Dir == eio.MD_Down{
		return true
	}
	return false
}