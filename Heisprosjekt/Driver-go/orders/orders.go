package orders

import (
	//ea "Driver-go/elev_actuator"
	eio "Driver-go/elevio"
	ec "Driver-go/elev_config"
	el "Driver-go/elev_logic"
// 	fsm "Driver-go/fsm"
	//"time"


	nw "Driver-go/network/bcast"

)

type OrderState int
const (
	UNASSIGNED OrderState = iota
	ASSIGNED
	COMPLETED
)

type OrderConfirmation int
const (
	UNCONFIRMED OrderConfirmation = iota
	CONFIRMED
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
	if btn_event.Button == eio.BT_Cab {
		o.OrderState = ASSIGNED
		o.AssignedElevator = elevID
	}
	return o
}

func AssignOrderToElevator(o *Order, e1 *ec.Elevator, e2 *ec.Elevator, e3 *ec.Elevator) {
	if o.OrderType == eio.BT_Cab {
		o.AssignedElevator = o.OriginElevator
		o.OrderState = ASSIGNED
		return
	}
	time1 := TimeToIdle(e1)
	time2 := TimeToIdle(e2)
	time3 := TimeToIdle(e3)
	if time1 <= time2 {
		if time1 <= time3 {
			o.AssignedElevator = e1.ElevID
		}else{
			o.AssignedElevator = e3.ElevID
		}
	}else{
		if time2 <= time3
	}
	


}

func TimeToIdle(e *ec.Elevator) (duration int) {
	duration = 0
	e_copy := e
	switch e.Behaviour {
	case ec.EB_Idle:
		return duration
	case ec.EB_DoorOpen:
		duration += int(ec.DOOR_TIMEOUT/2)
	case ec.EB_Moving:
		duration += int(ec.TRAVEL_TIME/2)
		e.Floor += int(e.Dir)
	}

	for {
 		if el.Stop_Here(e) {
			duration += int(ec.DOOR_TIMEOUT)
			e.Dir = el.Choose_Dir(e)
			el.Clear_Floor_Requests(e)
			if e.Dir == eio.MD_Stop {
				e = e_copy;
				return duration
			}
		}
		e.Floor += int(e.Dir)
		duration += int(ec.TRAVEL_TIME)
	}
}

func BroadcastOrderAndState(new_order Order, e *ec.Elevator)  {

	txOrderChan := make(chan Order)
	txElevChan := make(chan *ec.Elevator)
	

	go nw.Transmitter(20024, txOrderChan)
	go nw.Transmitter(20024, txElevChan)
	

	txOrderChan <- new_order
	txElevChan <- e

	
}

func RecieveOrderAndState(e2 *ec.Elevator, e3 *ec.Elevator)  {

	rxElevChan := make(chan *ec.Elevator)
	rxOrderChan := make(chan Order)

	go nw.Receiver(20023, rxElevChan)
	go nw.Receiver(20023, rxOrderChan)

	for{
		select{
		case new_order := <-rxOrderChan:
			MyWorldView = append(MyWorldView, &new_order)

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