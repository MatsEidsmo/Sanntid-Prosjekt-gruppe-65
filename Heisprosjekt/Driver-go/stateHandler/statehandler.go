package statehandler

import (
	ec "Driver-go/elev_config"
	//"fmt"
)

func RecieveAndUpdateStates(RecieveChan chan ec.Elevator, elevator_states map[string]ec.Elevator) {
	for {
		state_update := <- RecieveChan
		elevator_states[state_update.ElevID] = state_update
		//fmt.Println(elevator_states)
	}
}