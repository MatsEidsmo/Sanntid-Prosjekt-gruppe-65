package counter

import (

	 ec "Driver-go/elev_config"
	// el "Driver-go/elev_logic"
	 eio "Driver-go/elevio"
	bcast "Driver-go/network/bcast"
	hb "Driver-go/network/heartbeat"
	orders "Driver-go/orders"
	so "Driver-go/network/sendorders"
)

func ConfirmedQueue(wholeOrderList orders.OrderList) (confirmedOrderList orders.OrderList) {
	for _, o := range wholeOrderList {
		if o.OrderConfirmation == orders.CONFIRMED{
			confirmedOrderList = append(confirmedOrderList, o)
		}
	}
	return confirmedOrderList
}

func HandleButtonInput( e *ec.Elevator, pushed_btn chan eio.ButtonEvent, activeElevators map[string]hb.Heartbeat) (Assigned_btn chan eio.ButtonEvent) {
	
	RecieveWorldviewChan := make(chan orders.OrderList)
	go so.RecieveOrder(RecieveWorldviewChan)

	for {
		select{
		case btn := <- pushed_btn:
			o := orders.NewOrder(btn, e.ElevID)
			orders.MyWorldView = append(orders.MyWorldView, &o)
			BroadcastWorldview(orders.MyWorldView)


		case wv_update := <- RecieveWorldviewChan:
			for _, o := range wv_update{
				if o.OrderConfirmation ==  orders.UNCONFIRMED {

					if len(o.ElevsConfirmed) == len(activeElevators) {
						o.OrderConfirmation = orders.CONFIRMED
						BroadcastWorldview(wv_update)
	
					} else {
						elev_confirmed := false
						for _, id := range o.ElevsConfirmed{
							if e.ElevID == id {
								elev_confirmed = true
							}
						}
						if !elev_confirmed {
							o.ElevsConfirmed = append(o.ElevsConfirmed, e.ElevID)
							BroadcastWorldview(wv_update)
						}
					}
				}
			}
		}
	}
	return 
}

func BroadcastWorldview(WorldviewUpdate orders.OrderList) {
	TransmitWorldviewChan := make(chan orders.OrderList)
	TransmitWorldviewChan <- WorldviewUpdate
}
