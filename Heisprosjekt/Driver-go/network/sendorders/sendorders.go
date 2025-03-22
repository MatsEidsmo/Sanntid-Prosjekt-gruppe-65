package sendorders

import (
	orders "Driver-go/orders"
	"fmt"
)

func RecieveOrderList(rxOrderChan chan orders.OrderList) {
	for {
		wv_recieved := <-rxOrderChan
		fmt.Println("wv update recieved:", wv_recieved)
	}
}