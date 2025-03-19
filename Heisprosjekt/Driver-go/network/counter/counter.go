package counter

import (

	// ec "Driver-go/elev_config"
	// el "Driver-go/elev_logic"
	// eio "Driver-go/elevio"
	// nw "Driver-go/network/bcast"
	// hb "Driver-go/network/heartbeat"
	orders "Driver-go/orders"
)

func ConfirmedQueue(wholeOrderList orders.OrderList) (confirmedOrderList orders.OrderList) {
	for _, o := range wholeOrderList {
		if o.OrderConfirmation == orders.CONFIRMED{
			confirmedOrderList = append(confirmedOrderList, o)
		}
	}
	return confirmedOrderList
}

