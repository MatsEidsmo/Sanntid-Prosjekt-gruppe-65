package sendorders

import (
	orders "Driver-go/orders"
)

func RecieveOrderList(rxOrderChan chan orders.OrderList, OrderlistChan chan orders.OrderList) {
	for {
		OrderlistChan <- <-rxOrderChan
	}
}