package main

import (
	ea "Driver-go/elev_actuator"
	ec "Driver-go/elev_config"
	el "Driver-go/elev_logic"
	eio "Driver-go/elevio"
	fsm "Driver-go/fsm"
	bcast "Driver-go/network/bcast"
	hb "Driver-go/network/heartbeat"
	//so "Driver-go/network/sendorders"
	"Driver-go/orders"
	counter "Driver-go/network/counter"

	//"Driver-go/orders"

	"fmt"

	ip "Driver-go/network/localip"

	"flag"
	"os"
	"strconv"
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
    //e.ElevID = "Elevator1"
    
    

    

}

func main() {
	buff_size := 16*1024

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()
	
	if id == "" {
		localIP, err := ip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}
	
	port := 15657
	id_int, _ := strconv.Atoi(id)
	
	elev := ec.InitElev(id)
	e := &elev
	eio.Init("localhost:"+strconv.Itoa(port+id_int), ec.N_floors)
	//PeerList := make([]string, 0)
    //numFloors := 4
   
    
    //var e ec.Elevator
    
    
    //eio.Init("localhost:15657", numFloors)

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


    
    txhbChan := make(chan hb.Heartbeat)
	rxhbChan := make(chan hb.Heartbeat)
	
	RecieveWorldviewChan := make(chan orders.OrderList, buff_size)
	TransmitWorldviewChan := make(chan orders.OrderList, buff_size)
	

	activeElevators := make(map[string]hb.Heartbeat)

	go bcast.Transmitter(20023, txhbChan)
	go bcast.Receiver(20023, rxhbChan)
	go bcast.Transmitter(20023, TransmitWorldviewChan)
	go bcast.Receiver(20023, RecieveWorldviewChan)
	
	

	go hb.Transmitter(*e, txhbChan)
	go hb.Receiver(rxhbChan, activeElevators)
	go hb.RemoveInactiveElevators(activeElevators, 4*time.Second)
	//go counter.BroadcastWorldview(orders.MyWorldView,txOrderListChan)
	// go func() {
	// 	for recievedWorldview := range RecieveWorldviewChan {
	// 		fmt.Println("Recieved Worldview:", recievedWorldview)
	// 		orders.MyWorldView = recievedWorldview
	// 	}
	// }()
    
 	test_channel := make(chan eio.ButtonEvent)
	//block_chan := make(chan orders.OrderList)

    Initialize_Elev(e, drv_floors)

	go counter.HandleButtonInput(e, drv_buttons, activeElevators, RecieveWorldviewChan, TransmitWorldviewChan)


    defer fsm.Run(e, test_channel, drv_obstr, drv_floors, activeElevators)

}