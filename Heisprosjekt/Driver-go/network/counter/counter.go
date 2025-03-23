package counter

import (
	ec "Driver-go/elev_config"
	// el "Driver-go/elev_logic"
	eio "Driver-go/elevio"
	//bcast "Driver-go/network/bcast"
	//"Driver-go/network/bcast"
	hb "Driver-go/network/heartbeat"
	orders "Driver-go/orders"

	//so "Driver-go/network/sendorders"
	"fmt"
	//"time"
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
	activeElevators map[string]hb.Heartbeat, 
	recieve_chan chan orders.Order, 
	transmitt_chan chan orders.Order,
	transmitt_hb_chan chan hb.Heartbeat,
	send_to_fsm chan eio.ButtonEvent,
	) {
	
	
	

	for {
		select{
		case btn := <- pushed_btn:
			fmt.Println("Inside Button pushed")
			o := orders.NewOrder(btn, e.ElevID)

			//if !orders.IsOrderInWorldview(o) {
				
				//orders.MyWorldView = append(orders.MyWorldView, &o)
				
				transmitt_chan <- o
				
			//}
	



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
				OrderToBeSent := rec_order

				
				
				
			case orders.COMPLETED:
				//Delete order
			}
			
			orders.AssignOrderToElevator(&rec_order, activeElevators)
			fmt.Println("My elevID:",e.ElevID)
			fmt.Println("Assigned ElevID:", rec_order.AssignedElevator)
			if e.ElevID == rec_order.AssignedElevator {
				fmt.Println("Order assigned to ME:)")
				send_to_fsm <- eio.ButtonEvent{rec_order.OrderFloor,rec_order.OrderType}
			}else{
				fmt.Println("Order Assigned to someone else")
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
			
			//}
			
			
			
		}
	}
}

func BroadcastOrder(OrderUpdate orders.Order, transmitChan chan orders.Order) {
	

	transmitChan <- OrderUpdate
	

}
