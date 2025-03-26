package counter

import (
	ec "Driver-go/elev_config"
	//"time"

	// el "Driver-go/elev_logic"
	eio "Driver-go/elevio"
	//bcast "Driver-go/network/bcast"
	//"Driver-go/network/bcast"
	hb "Driver-go/network/heartbeat"
	orders "Driver-go/orders"

	//so "Driver-go/network/sendorders"
	"fmt"
	
)



func ConfirmedQueue(wholeOrderList orders.OrderList) (confirmedOrderList orders.OrderList) {
	for _, o := range wholeOrderList {
		if o.OrderConfirmation == orders.CONFIRMED{
			confirmedOrderList = append(confirmedOrderList, o)
		}
	}
	return confirmedOrderList
}

func HandleButtonInput( 
	e *ec.Elevator, 
	pushed_btn chan eio.ButtonEvent, 
	recieve_chan chan orders.Order, 
	transmitt_chan chan orders.Order,
	transmitt_state_chan chan ec.Elevator,
	activeElevators map[string]hb.Heartbeat, 
	elevatorStates map[string]ec.Elevator,
	send_to_fsm chan eio.ButtonEvent,
	) {
	
	
	

	for {
		select{
		case btn := <- pushed_btn:
			o := orders.NewOrder(btn, e.ElevID, send_to_fsm)
			
			for _, e := range elevatorStates {
				if e.Floor == o.OrderFloor && e.Behaviour != ec.EB_Moving {
					o.AssignedElevator = e.ElevID
					

				}
			}
			
				
			fmt.Println("Sending Button WÆÆÆÆÆÆÆ")
			//orders.MyWorldView = append(orders.MyWorldView, &o)
			
			transmitt_chan <- o
			
			
			
			



		case rec_order := <- recieve_chan:
			
			// CONFIRM ORDER
			switch rec_order.OrderConfirmation {
			case orders.UNCONFIRMED:
				//fmt.Println("Order is unconfirmed")
				if !orders.IsElevConfirmed(e, rec_order) {
					rec_order.ElevsConfirmed = append(rec_order.ElevsConfirmed, e.ElevID)
					
					if len(rec_order.ElevsConfirmed) == len(activeElevators) {
						rec_order.OrderConfirmation = orders.CONFIRMED
						orders.MyWorldView = append(orders.MyWorldView, &rec_order)
						
						
					}
					BroadcastOrder(rec_order, transmitt_chan)
					//fmt.Println("Broadcasted order:", rec_order)
				}
				
			case orders.CONFIRMED:
				fmt.Println("Num Active Elevs:",len(activeElevators))
				fmt.Println("Num Elevators in elevstates:", len(elevatorStates))
				if rec_order.OrderState != orders.ASSIGNED {
					orders.AssignOrderToElevator(&rec_order, elevatorStates)
					eio.SetButtonLamp(rec_order.OrderType, rec_order.OrderFloor, true)
					
					if rec_order.AssignedElevator == e.ElevID{

						fmt.Println("Order assigned to ME:)")
						send_to_fsm <- eio.ButtonEvent{rec_order.OrderFloor,rec_order.OrderType}
					}

				}
				
				
				
				
				
			case orders.COMPLETED:
				eio.SetButtonLamp(rec_order.OrderType, rec_order.OrderFloor, false)
				
				filtered_wv := orders.MyWorldView[:0]
				for _, order := range orders.MyWorldView {
					if order.OrderID != rec_order.OrderID {
						filtered_wv = append(filtered_wv, order)
					}else{
						BroadcastOrder(rec_order, transmitt_chan)
					}
				}
				orders.MyWorldView = filtered_wv
			}
			



			// for _, o := range orders.MyWorldView {
			// 	if o.OrderConfirmation ==  orders.UNCONFIRMED {

			// 		if len(o.ElevsConfirmed) == len(activeElevators) {
			// 			o.OrderConfirmation = orders.CONFIRMED
			// 			BroadcastOrder(wv_update, transmitt_chan)
	
			// 		} else {
			// 			elev_confirmed := false
			// 			for _, id := range o.ElevsConfirmed {
			// 				if e.ElevID == id {
			// 					elev_confirmed = true
			// 				}
			// 			}
			// 			if !elev_confirmed { 
			// 				o.ElevsConfirmed = append(o.ElevsConfirmed, e.ElevID)
			// 				BroadcastOrder(wv_update, transmitt_chan)
			// 			}
			// 		}
			// 	}
			// 	fmt.Println(o)
			// }
			
			
		}
	}
}

func BroadcastOrder(OrderUpdate orders.Order, transmitChan chan orders.Order) {
	

	transmitChan <- OrderUpdate
	

}
