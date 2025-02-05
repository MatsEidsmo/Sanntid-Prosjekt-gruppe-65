package requests

import "elev_config"


func requests_above(e Elevator) int {
	for f := e.floor+1; f < n_floors; f++ {
		for btn := e.button; btn < n_buttons; btn++ {
			if e.requestMatrix[f][btn] {
				return 1
			}
		}
	}
	return 0
}

func requests_below(e Elevator) int {
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < n_buttons; btn++ {
			if e.requestMatrix[f][btn] {
				return 1
			}
		}
	}
	return 0
}

//Git works