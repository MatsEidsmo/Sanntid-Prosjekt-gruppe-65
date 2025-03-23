package order_timeout

import (
	"time"
	orders "Driver-go/orders"
	hb "Driver-go/network/heartbeat"
	ec "Driver-go/elev_config"
)

// Problem: Potentially active_elevators is not updated. So if the assigned_elevator is the only elevator then we do cannot assign order to
// another elevator. But if an elevator enters the network after this decision, yet the assigned_elevator is not free to go, no reassignments are
// made leading to the order being lost. Similarly if after the decision the assigned_elev disconnects without completing the order, the order is
// lost. 
// Solutions/Thoughts: this function has to be a goroutine to always keep track, active_elevators could be a global variable to always have the most
// accessable available
func OrderTimeout(order orders.Order, assigned_elevator *ec.Elevator, active_elevators map[string]hb.Heartbeat, fire_duration_in_millisec int) {
	timer := time.NewTimer(time.Duration(fire_duration_in_millisec)*time.Millisecond)
	<-timer.C
	if order.OrderConfirmation != orders.COMPLETED {
		// Assign again discarding this elevator
		RemoveActiveElevator()
		orders.AssignOrderToElevator(&order, active_elevators)
		AddActiveElevator()
	}
}