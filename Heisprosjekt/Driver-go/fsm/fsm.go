package fsm

import (
	ea "Driver-go/elev_actuator"
	ec "Driver-go/elev_config"
	el "Driver-go/elev_logic"
	eio "Driver-go/elevio"
	"fmt"
)

func Run(e *ec.Elevator, pushed_btn chan eio.ButtonEvent, obstr_chann chan bool, floor_sensor chan int) {
	door_open_flag := false
	for {
		select {
		case btn := <- pushed_btn:
			el.Add_Request(e, btn.Floor, btn.Button)
			fmt.Println("btn pushed recieved")
			//if e.Dir == eio.MD_Stop {
				curr_dir := el.Choose_Dir(e)
				if btn.Floor == e.Floor && e.Behaviour != ec.EB_Moving{
					ea.Timer_start()
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
				ea.Open_Door()
				e.Behaviour = ec.EB_DoorOpen
				
			}
		case <- ea.DoorTimer.C:
			door_open_flag = false
			ea.Upon_Door_Timeout(e)
			// curr_dir := el.Choose_Dir(e)
			// eio.SetDoorOpenLamp(false)
			// el.Clear_Floor_Requests(e)
			// if curr_dir == eio.MD_Stop {
			// 	e.Behaviour = ec.EB_Idle
			// } else {
			// 	e.Behaviour = ec.EB_Moving
			// }

			// eio.SetMotorDirection(curr_dir)
			

		case obstr := <- obstr_chann:
			println("Obstruction!")
			e.Obstruction = obstr
			if !obstr && e.Behaviour == ec.EB_DoorOpen {
				ea.Open_Door()
			}
		}
	}
}