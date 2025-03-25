package order_timeout

import (
	"time"
	orders "Driver-go/orders"
	ec "Driver-go/elev_config"
	eio "Driver-go/elevio"
)

// Problem: Potentially active_elevators is not updated. So if the assigned_elevator is the only elevator then we do cannot assign order to
// another elevator. But if an elevator enters the network after this decision, yet the assigned_elevator is not free to go, no reassignments are
// made leading to the order being lost. Similarly if after the decision the assigned_elev disconnects without completing the order, the order is
// lost. 
// Solutions/Thoughts: this function has to be a goroutine to always keep track, active_elevators could be a global variable to always have the most
// accessable available
func OrderTimeout(
	recieve_chan chan orders.Order,
	order_reassigned chan bool,  
	assigned_elevator *ec.Elevator, 
	active_elevators map[string]ec.Elevator, 
	tolarance_duration int,
	) {



	for {
		select {
		case order:= <- recieve_chan:
			order_completion_duration := orders.TimeToRequestHandled(assigned_elevator, &order)
			timer := time.NewTimer(time.Duration(order_completion_duration + tolarance_duration))
			<-timer.C
			if order.OrderConfirmation != orders.COMPLETED {
				// Assign again discarding this elevator
				RemoveActiveElevator()
				if len(active_elevators) < 2 {
					AddActiveElevator()
					break
				}
				orders.AssignOrderToElevator(order, active_elevators)
				AddActiveElevator()
				order_reassigned <- true
			}

		case <-order_reassigned:
			order_completion_duration := orders.TimeToRequestHandled(assigned_elevator, &order)
			timer := time.NewTimer(time.Duration(order_completion_duration + tolarance_duration))
			<-timer.C
			if order.OrderConfirmation != orders.COMPLETED {
				// Assign again discarding this elevator
				RemoveActiveElevator()
				if len(active_elevators) < 2 {
					AddActiveElevator()
					break
				}
				orders.AssignOrderToElevator(order, active_elevators)
				AddActiveElevator()
				order_reassigned <- true
			}
		}
	}
}