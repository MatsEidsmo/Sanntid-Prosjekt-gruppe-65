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
	eio.SetButtonLamp(btnType, floor, true)
}

func Clear_Floor_Requests(e *ec.Elevator) {
	for btn := 0; btn < ec.N_buttons; btn++{
		e.RequestMatrix[e.Floor][btn] = 0
	}
	eio.SetButtonLamp(eio.BT_HallUp, e.Floor, false)
	eio.SetButtonLamp(eio.BT_HallDown, e.Floor, false)
	eio.SetButtonLamp(eio.BT_Cab, e.Floor, false)
}

func Choose_Dir(e *ec.Elevator) {
	switch e.Dir {
	case eio.MD_Up:
		if requests_above(e) {
			e.Dir = eio.MD_Up
			eio.SetMotorDirection(e.Dir)
		} else if requests_below(e) {
			e.Dir = eio.MD_Down
			eio.SetMotorDirection(e.Dir)
		} else {
			eio.SetMotorDirection(eio.MD_Stop)
		}
	case eio.MD_Down:
		if requests_below(e) {
			e.Dir = eio.MD_Down
			eio.SetMotorDirection(e.Dir)
		} else if requests_above(e) {
			e.Dir = eio.MD_Up
			eio.SetMotorDirection(e.Dir)
		} else {
			eio.SetMotorDirection(eio.MD_Stop)
		}
	case eio.MD_Stop:
		if requests_below(e) {
			e.Dir = eio.MD_Down
			eio.SetMotorDirection(e.Dir)
		} else if requests_above(e) {
			e.Dir = eio.MD_Up
			eio.SetMotorDirection(e.Dir)
		} else {
			eio.SetMotorDirection(eio.MD_Stop)
		}
	}

}

func Stop_Here(e *ec.Elevator) bool {
	if e.RequestMatrix[e.Floor][eio.BT_Cab] == 1 {
		return true
	}
	if e.RequestMatrix[e.Floor][eio.BT_HallUp] == 1 && (e.Dir == eio.MD_Up || !requests_below(e)){
		return true
	}
	if e.RequestMatrix[e.Floor][eio.BT_HallDown] == 1 && (e.Dir == eio.MD_Down || !requests_above(e)){
		return true
	}
	return false
}

func Clear_RequestMatrix(e *ec.Elevator) {
	for i := 0; i < ec.N_floors; i++ {
		for j := 0; j < ec.N_buttons; j++ {
			e.RequestMatrix[i][j] = 0
		}
		eio.SetButtonLamp(eio.BT_HallUp, i, false)
		eio.SetButtonLamp(eio.BT_HallDown, i, false)
		eio.SetButtonLamp(eio.BT_Cab, i, false)
	}
}