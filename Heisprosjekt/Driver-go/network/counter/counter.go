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
)



func ConfirmedQueue(wholeOrderList orders.OrderList) (confirmedOrderList orders.OrderList) {
	for _, o := range wholeOrderList {
		if o.OrderConfirmation == orders.CONFIRMED{
			confirmedOrderList = append(confirmedOrderList, o)
		}
	}
	return confirmedOrderList
}

func HandleButtonInput( e *ec.Elevator, pushed_btn chan eio.ButtonEvent, activeElevators map[string]hb.Heartbeat, recieve_chan chan orders.Order, transmitt_chan chan orders.Order) {
	
	
	

	for {
		select{
		case btn := <- pushed_btn:
			fmt.Println("Inside Button pushed")
			o := orders.NewOrder(btn, e.ElevID)

			if !orders.IsOrderInWorldview(o) {
				//orders.MyWorldView = append(orders.MyWorldView, &o)
				
				transmitt_chan <- o

			}
	



		case rec_order := <- recieve_chan:
			fmt.Println("Inside Recieved Worldview")
			// CONFIRM ORDER
			switch rec_order.OrderConfirmation {
			case orders.UNCONFIRMED:
				if !orders.IsElevConfirmed(e, rec_order) {
					rec_order.ElevsConfirmed = append(rec_order.ElevsConfirmed, e.ElevID)
					
					if len(rec_order.ElevsConfirmed) == len(activeElevators) {
						rec_order.OrderConfirmation = orders.CONFIRMED
						orders.MyWorldView = append(orders.MyWorldView, &rec_order)
					}
				}
				BroadcastOrder(rec_order, transmitt_chan)
				
			case orders.CONFIRMED:
				// Do Something else
			
			case orders.COMPLETED:
				//Delete order
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
	fmt.Println("Inside Broadcast Wv")

	transmitChan <- OrderUpdate
	fmt.Println(transmitChan)

}
