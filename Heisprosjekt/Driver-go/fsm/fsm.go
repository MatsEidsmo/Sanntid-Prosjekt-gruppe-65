package fsm

import (
	ea "Driver-go/elev_actuator"
	ec "Driver-go/elev_config"
	el "Driver-go/elev_logic"
	eio "Driver-go/elevio"
	"fmt"
)

func Run(e *ec.Elevator, pushed_btn chan eio.ButtonEvent, obstr_chann chan bool, floor_sensor chan int) {
	for {
		select {
		case btn := <- pushed_btn:
			el.Add_Request(e, btn.Floor, btn.Button)
			if e.Dir == eio.MD_Stop {
				el.Choose_Dir(e)

			}
		case floor := <- floor_sensor:
			eio.SetFloorIndicator(floor)
			fmt.Println("Floor: " , floor)
			fmt.Println("Dir: " , e.Dir)
			e.Floor = floor
			if el.Stop_Here(e) {
				fmt.Println("Elevator stopping")
				ea.Open_Door()
				e.Dir = eio.MD_Stop
				el.Clear_Floor_Requests(e)
				fmt.Println("Req: ", e.RequestMatrix[3][2])  // Nye bestillinger blir ike lagret i matrisen mens døra er åpen
				el.Choose_Dir(e)
			}

		case obstr := <- obstr_chann:
			if obstr {
				//Stop()
			}
		}
	}
}