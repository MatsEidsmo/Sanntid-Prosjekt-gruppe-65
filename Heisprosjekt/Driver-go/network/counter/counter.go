package counter

import (

	 ec "Driver-go/elev_config"
	// el "Driver-go/elev_logic"
	 eio "Driver-go/elevio"
	//bcast "Driver-go/network/bcast"
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

func HandleButtonInput( e *ec.Elevator, pushed_btn chan eio.ButtonEvent, activeElevators map[string]hb.Heartbeat, recieved_info chan orders.OrderList) {
	
	TransmitWorldviewChan := make(chan orders.OrderList)
	//RecieveWorldviewChan := make(chan orders.OrderList)
	

	for {
		select{
		case btn := <- pushed_btn:
			fmt.Println("Inside Button pushed")
			o := orders.NewOrder(btn, e.ElevID)
			orders.MyWorldView = append(orders.MyWorldView, &o)
			BroadcastWorldview(orders.MyWorldView, TransmitWorldviewChan)
			fmt.Println("End of Button pushed")



		case wv_update := <- recieved_info:
			fmt.Println("Inside Recieved Worldview")

			for _, o := range wv_update{
				if o.OrderConfirmation ==  orders.UNCONFIRMED {

					if len(o.ElevsConfirmed) == len(activeElevators) {
						o.OrderConfirmation = orders.CONFIRMED
						BroadcastWorldview(wv_update, TransmitWorldviewChan)
	
					} else {
						elev_confirmed := false
						for _, id := range o.ElevsConfirmed{
							if e.ElevID == id {
								elev_confirmed = true
							}
						}
						if !elev_confirmed {
							o.ElevsConfirmed = append(o.ElevsConfirmed, e.ElevID)
							BroadcastWorldview(wv_update, TransmitWorldviewChan)
						}
					}
				}
			}
			fmt.Println("upsi")
		}
	}
}

func BroadcastWorldview(WorldviewUpdate orders.OrderList, transmitChan chan orders.OrderList) {
	fmt.Println("Inside Broadcast Wv")

	transmitChan <- WorldviewUpdate
	fmt.Println("End of broadcast")

}
