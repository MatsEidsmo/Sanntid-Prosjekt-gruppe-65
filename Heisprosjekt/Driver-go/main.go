package main

import (
	ea "Driver-go/elev_actuator"
	ec "Driver-go/elev_config"
	el "Driver-go/elev_logic"
	eio "Driver-go/elevio"
	fsm "Driver-go/fsm"
	bcast "Driver-go/network/bcast"
	hb "Driver-go/network/heartbeat"
	"fmt"

	//"fmt"

	ip "Driver-go/network/localip"
	//"fmt"
	"time"
)


func Initialize_Elev(e *ec.Elevator, drv_floors chan int) {
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

    ea.Timer_init()

    el.Clear_RequestMatrix(e)
    
    e.Behaviour = ec.EB_Idle
    var err error
    e.ElevID, err = ip.LocalIP()
    if err != nil {
        fmt.Println("Error fetching local IP:", err)
        
    }
    
    

    

}

func main() {
    numFloors := 4
   
    
    var e ec.Elevator
    
    
    eio.Init("localhost:15657", numFloors)

	var d eio.MotorDirection = eio.MD_Down
	eio.SetMotorDirection(d)

    drv_floors := make(chan int)
	drv_buttons := make(chan eio.ButtonEvent)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
    
    go eio.PollFloorSensor(drv_floors)
	go eio.PollButtons(drv_buttons)
	go eio.PollObstructionSwitch(drv_obstr)
	go eio.PollStopButton(drv_stop)


    
    txChan := make(chan hb.Heartbeat)
	rxChan := make(chan hb.Heartbeat)
	activeElevators := make(map[string]hb.Heartbeat)

	go bcast.Transmitter(20023, txChan)
	go bcast.Receiver(20023, rxChan)

	
	

	go hb.Transmitter(e, txChan)
	go hb.Receiver(rxChan, activeElevators)
	go hb.RemoveInactiveElevators(activeElevators, 4*time.Second)
    
    
    Initialize_Elev(&e, drv_floors)


    defer fsm.Run(&e, drv_buttons, drv_obstr, drv_floors, activeElevators)

}