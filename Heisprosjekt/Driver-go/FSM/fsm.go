package FSM

import (
	eio "Driver-go/elevio"
	ec "Driver-go/elev_config"
	el "Driver-go/elev_logic"
)

func run(e *ec.Elevator, pushed_btn chan eio.ButtonEvent, obstr_chann chan bool, floor_sensor chan int) {
	for {
		select {
		case btn := <- pushed_btn:
			el.Add_Request(e, btn.Floor, btn.Button)

		case floor := <- floor_sensor:
			e.Floor = floor
			if el.Stop_Here(e) {
				el.Open_Door()
				el.Remove_Request(e)
			}

		case obstr := <- obstr_chann:
			if obstr {
				//Stop()
			}
		}
	}
}