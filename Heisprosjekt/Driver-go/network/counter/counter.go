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

func HandleButtonInput( e *ec.Elevator, pushed_btn chan eio.ButtonEvent, activeElevators map[string]hb.Heartbeat, recieve_chan chan orders.OrderList, transmitt_chan chan orders.OrderList) {
	
	
	

	for {
		select{
		case btn := <- pushed_btn:
			fmt.Println("Inside Button pushed")
			o := orders.NewOrder(btn, e.ElevID)
			orders.MyWorldView = append(orders.MyWorldView, &o)
			
			transmitt_chan <- orders.MyWorldView
	



		case wv_update := <- recieve_chan:
			fmt.Println("Inside Recieved Worldview")

			// CONFIRM ORDER
			for _, o := range wv_update {
				if o.OrderConfirmation ==  orders.UNCONFIRMED {

					if len(o.ElevsConfirmed) == len(activeElevators) {
						o.OrderConfirmation = orders.CONFIRMED
						BroadcastWorldview(wv_update, transmitt_chan)
	
					} else {
						elev_confirmed := false
						for _, id := range o.ElevsConfirmed {
							if e.ElevID == id {
								elev_confirmed = true
							}
						}
						if !elev_confirmed { 
							o.ElevsConfirmed = append(o.ElevsConfirmed, e.ElevID)
							BroadcastWorldview(wv_update, transmitt_chan)
						}
					}
				}
				fmt.Println(o)
			}
			
			
		}
	}
}

func BroadcastWorldview(WorldviewUpdate orders.OrderList, transmitChan chan orders.OrderList) {
	fmt.Println("Inside Broadcast Wv")

	transmitChan <- WorldviewUpdate
	fmt.Println(transmitChan)

}
