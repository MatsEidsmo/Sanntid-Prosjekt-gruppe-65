package main

import (
	ea "Driver-go/elev_actuator"
	ec "Driver-go/elev_config"
	el "Driver-go/elev_logic"
	eio "Driver-go/elevio"
	fsm "Driver-go/fsm"

	"Driver-go/orders"

	//"time"

	//"time"
	//nw "Driver-go/network/bcast"
	"fmt"
)


func Initialize_Elev_Pos(e *ec.Elevator, drv_floors chan int) {
    floornumber := <-drv_floors
    eio.SetMotorDirection(eio.MD_Down)
    for floornumber != 0 {
        floornumber := <-drv_floors
        eio.SetFloorIndicator(floornumber)
        if floornumber == 0 {
            break
        } 
    }
    eio.SetMotorDirection(eio.MD_Stop)
    e.Dir = eio.MD_Stop

    el.Clear_RequestMatrix(e)
    
    e.Behaviour = ec.EB_Idle

    

    

}

func main() {
    numFloors := 4
   
    
    var e ec.Elevator
    var e2 ec.Elevator
    var e3 ec.Elevator

	eio.Init("localhost:15657", numFloors)

	var d eio.MotorDirection = eio.MD_Down
	eio.SetMotorDirection(d)

    drv_floors := make(chan int)
    go eio.PollFloorSensor(drv_floors)
    fmt.Println("1")
    ea.Timer_init()

    Initialize_Elev_Pos(&e, drv_floors)
    ec.DefaultElev(&e2)
    ec.DefaultElev(&e3)

	drv_buttons := make(chan eio.ButtonEvent)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go eio.PollButtons(drv_buttons)
	go eio.PollObstructionSwitch(drv_obstr)
	go eio.PollStopButton(drv_stop)

    eio.SetButtonLamp(eio.BT_Cab, 0, false)
    eio.SetButtonLamp(eio.BT_Cab, 1, false)
    eio.SetButtonLamp(eio.BT_Cab, 2, false)
    eio.SetButtonLamp(eio.BT_Cab, 3, false)


    
    go orders.RecieveOrderAndState(&e2, &e3)
    
    


    defer fsm.Run(&e, drv_buttons, drv_obstr, drv_floors)

}