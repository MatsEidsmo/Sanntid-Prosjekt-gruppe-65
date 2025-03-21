package fsm

import (
	ea "Driver-go/elev_actuator"
	ec "Driver-go/elev_config"
	el "Driver-go/elev_logic"
	eio "Driver-go/elevio"
	hb "Driver-go/network/heartbeat"
	//bcast "Driver-go/network/bcast"
	orders "Driver-go/orders"
	"fmt"
	//"time"
)

func Run(e *ec.Elevator, pushed_btn chan eio.ButtonEvent, obstr_chann chan bool, floor_sensor chan int, active_elevs map[string]hb.Heartbeat) {
	
	for {
		select {
		case btn := <- pushed_btn:
			fmt.Println("Button recieved!")

			el.Add_Request(e, btn.Floor, btn.Button)

			

			new_order := orders.NewOrder(btn, e.ElevID)
			

			
			fmt.Println("Before ATE")
			orders.AssignOrderToElevator(&new_order, active_elevs)
			
			fmt.Println("Assigned to Elevator")

			//txChan := make(chan eio.ButtonEvent)
			//rxChan := make(chan string)
			

    
            
            //time.Sleep(2*time.Second)
        
    

    
            //Recieved_btn := <- rxBtnChan
            //eio.PrintButtonEvent(Recieved_btn)
      

			// Gi beskjed om mottat btn.
			//Vent på confirmation om btn
			// Regn ut hvilken heis som skal kjøre
			// if Heis som skal kjøre == denne heisen
			if e.ElevID == new_order.AssignedElevator {

				
				
					curr_dir := el.Choose_Dir(e)
					if btn.Floor == e.Floor && e.Behaviour != ec.EB_Moving{
						ea.Timer_start()
					}
					if e.Behaviour != 0 {
						eio.SetMotorDirection(curr_dir)
					}
			} 

				
			
		case floor := <- floor_sensor:
			eio.SetFloorIndicator(floor)
			// fmt.Println("Floor: " , floor)
			// fmt.Println("Dir: " , e.Dir)
			e.Floor = floor
			if el.Stop_Here(e) {
				// fmt.Println("Elevator stopping")
				ea.Open_Door(e)
				e.Behaviour = ec.EB_DoorOpen
				
			}
		case <- ea.DoorTimer.C:
			
			ea.Upon_Door_Timeout(e)
			
			

		case obstr := <- obstr_chann:
			e.Obstruction = obstr
			if !obstr && e.Behaviour == ec.EB_DoorOpen {
				ea.Open_Door(e)
			}
		}
	}
}