package orders

import (
	//ea "Driver-go/elev_actuator"
	ec "Driver-go/elev_config"
	el "Driver-go/elev_logic"
	eio "Driver-go/elevio"
	//"container/list"
	"fmt"
	//"strings"

	// 	fsm "Driver-go/fsm"
	//"time"

	bcast "Driver-go/network/bcast"
	hb "Driver-go/network/heartbeat"
	"sort"
)

type OrderState int
const (
	UNASSIGNED OrderState = iota
	ASSIGNED
	
)

type OrderConfirmation int
const (
	UNCONFIRMED OrderConfirmation = iota
	CONFIRMED
	COMPLETED
)

type OrderType int
const (
	HALL_UP OrderType = iota
	HALL_DOWN
	CAB
)

type Order struct {
	OrderType eio.ButtonType
	OrderFloor int
	OrderState OrderState
	OriginElevator string
	AssignedElevator string
	OrderConfirmation OrderConfirmation
	ElevsConfirmed []string

}

type OrderList []*Order

var MyWorldView OrderList

func NewOrder(btn_event eio.ButtonEvent, elevID string) Order {
	o:= Order{
		OrderState: 	UNASSIGNED,
		OrderType: 		btn_event.Button,
		OrderFloor: 	btn_event.Floor,
		OriginElevator: elevID,
	}
	//o.ElevsConfirmed = append(o.ElevsConfirmed, o.OriginElevator)
	if btn_event.Button == eio.BT_Cab {
		o.OrderState = ASSIGNED
		o.AssignedElevator = elevID
	}
	return o
}

type calc_struct struct {
	id string
	duration int
}

func AssignOrderToElevator(o *Order, active_elevs map[string]hb.Heartbeat) {
	

	var tth_arr []calc_struct
	
	var min_ElevID string
	for id, hb := range active_elevs {
		
		
		curr_tth := TimeToRequestHandled(&hb.Elevator, o)
		fmt.Println("elev", id, "has calculated tth:", curr_tth)
		tth_arr = append(tth_arr, calc_struct{id,curr_tth})
		
		
	}
	
	sort.Slice(tth_arr, func(i, j int) bool {
        return tth_arr[i].duration < tth_arr[j].duration
    })

	var equal_durations []string
	for _, tth_struct := range tth_arr{
		if tth_struct.duration == tth_arr[0].duration{
			equal_durations = append(equal_durations, tth_struct.id)
		}
	}

	sort.Slice(equal_durations, func(i, j int) bool {
        return equal_durations[i] < equal_durations[j]
    })

	min_ElevID = equal_durations[0]
   
	o.AssignedElevator = min_ElevID
	o.OrderState = ASSIGNED
}

	// time1 := TimeToIdle(e1)
	// time2 := TimeToIdle(e2)
	// time3 := TimeToIdle(e3)
	// if time1 <= time2 {
	// 	if time1 <= time3 {
	// 		o.AssignedElevator = e1.ElevID
	// 	}
	// }else if time2 <= time3{
	// 	o.AssignedElevator = e2.ElevID
	// }else {
	// 	o.AssignedElevator = e3.ElevID
	// }
	




func TimeToRequestHandled(e *ec.Elevator, o *Order) (duration int) {
	
	duration = 0
	e_floor_copy := e.Floor
	e_dir_copy := e.Dir
	e_rm_copy := e.RequestMatrix
	el.Add_Request(e, o.OrderFloor, o.OrderType)
	
	// e.Dir = 1
	// e.Floor = 1
	switch e.Behaviour {
	case ec.EB_Idle:
		return duration
	case ec.EB_DoorOpen:
		//fmt.Println("Door open")
		duration += int(ec.DOOR_TIMEOUT/2)
		e.Dir = el.Choose_Dir(e)
	case ec.EB_Moving:
		//fmt.Println("Moving!")
		duration += int(ec.TRAVEL_TIME/2)
		e.Floor += int(e.Dir)
	}
	// fmt.Println("Floor:", e.Floor, "Dir:", e.Dir)
	// fmt.Println(e.RequestMatrix[3][0], e.RequestMatrix[3][1], e.RequestMatrix[3][2])
	// fmt.Println(e.RequestMatrix[2][0], e.RequestMatrix[2][1], e.RequestMatrix[2][2])
	// fmt.Println(e.RequestMatrix[1][0], e.RequestMatrix[1][1], e.RequestMatrix[1][2])
	// fmt.Println(e.RequestMatrix[0][0], e.RequestMatrix[0][1], e.RequestMatrix[0][2])
	
	
	for {
		
		if el.Stop_Here(e) {
			//fmt.Println("Should stop here")
			duration += int(ec.DOOR_TIMEOUT)
			e.Dir = el.Choose_Dir(e)
			el.Clear_Floor_Requests(e, true)
			if e.Dir == eio.MD_Stop {
				e.Floor = e_floor_copy
				e.Dir = e_dir_copy
				e.RequestMatrix = e_rm_copy
				fmt.Println("Duration:", duration)
				return duration
			}
		}
		e.Floor += int(e.Dir)
		duration += int(ec.TRAVEL_TIME)
		//fmt.Println(e.Floor)
	}
}

func BroadcastOrderAndState(new_order Order, e *ec.Elevator)  {

	txOrderChan := make(chan Order)
	txElevChan := make(chan *ec.Elevator)
	

	go bcast.Transmitter(20023, txOrderChan)
	go bcast.Transmitter(20023, txElevChan)
	

	txOrderChan <- new_order
	txElevChan <- e

	
}

func RecieveOrderAndState(e2 *ec.Elevator, e3 *ec.Elevator)  {

	rxElevChan := make(chan *ec.Elevator)
	rxOrderChan := make(chan Order)

	go bcast.Receiver(20023, rxElevChan)
	go bcast.Receiver(20023, rxOrderChan)

	for{
		select{
		case new_order := <-rxOrderChan:
			MyWorldView = append(MyWorldView, &new_order)

			// Wait for confirmation on worldview here

			


		case state_update := <- rxElevChan:
			if state_update.ElevID == e2.ElevID{
				e2 = state_update
			}
			if state_update.ElevID == e3.ElevID{
				e3 = state_update
			}
		}
		//PeerElev := <-rxElevChan
	}
	
}

func IsOrderConfirmedInWorldview(o Order) bool {
	in_wv := false
	for _, wv_orders := range MyWorldView{
		if o.OrderType == wv_orders.OrderType && o.OrderFloor == wv_orders.OrderFloor && o.OriginElevator == wv_orders.OriginElevator {
			in_wv = true
		}
	}
	return in_wv
}

func IsElevConfirmed(e *ec.Elevator, order Order) bool {
	ret_val := false
	for _, id := range order.ElevsConfirmed{
		if id == e.ElevID {
			ret_val = true
		}
	}
	return ret_val
}

