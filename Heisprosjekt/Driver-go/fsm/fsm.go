package fsm

import (
	ea "Driver-go/elev_actuator"
	ec "Driver-go/elev_config"
	el "Driver-go/elev_logic"
	eio "Driver-go/elevio"
	"fmt"
)

func Run(e *ec.Elevator, pushed_btn chan eio.ButtonEvent, obstr_chann chan bool, floor_sensor chan int, door_timer chan bool) {
	door_open_flag := false
	for {
		select {
		case btn := <- pushed_btn:
			el.Add_Request(e, btn.Floor, btn.Button)
			//if e.Dir == eio.MD_Stop {
				curr_dir := el.Choose_Dir(e)
				if btn.Floor == e.Floor && e.Dir == eio.MD_Stop{
					ea.Open_Door(door_timer, e)
				}
				if !door_open_flag {
					eio.SetMotorDirection(curr_dir)
				}
				
			
		case floor := <- floor_sensor:
			eio.SetFloorIndicator(floor)
			fmt.Println("Floor: " , floor)
			fmt.Println("Dir: " , e.Dir)
			e.Floor = floor
			if el.Stop_Here(e) {
				fmt.Println("Elevator stopping")
				door_open_flag = true
				ea.Open_Door(door_timer, e)
				e.Dir = eio.MD_Stop
				fmt.Println("Req: ", e.RequestMatrix[3][2])  // Nye bestillinger blir ike lagret i matrisen mens døra er åpen
				
			}
		case <-door_timer:
			door_open_flag = false
			curr_dir := el.Choose_Dir(e)
			eio.SetMotorDirection(curr_dir)

		case obstr := <- obstr_chann:
			if obstr {
				//Stop()
			}
		}
	}
}