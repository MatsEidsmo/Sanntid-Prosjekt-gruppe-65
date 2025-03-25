package order_timeout

import (
	"time"
	orders "Driver-go/orders"
	ec "Driver-go/elev_config"
	"sync"
)

// Problem: Potentially active_elevators is not updated. So if the assigned_elevator is the only elevator then we do cannot assign order to
// another elevator. But if an elevator enters the network after this decision, yet the assigned_elevator is not free to go, no reassignments are
// made leading to the order being lost. Similarly if after the decision the assigned_elev disconnects without completing the order, the order is
// lost. 
// Solutions/Thoughts: this function has to be a goroutine to always keep track, active_elevators could be a global variable to always have the most
// accessable available
func OrderTimeout(
	recieve_chan chan orders.Order,
	order_reassigned chan orders.Order, 
	active_elevators map[string]ec.Elevator, 
	tolarance_duration int,
	) {

	var mu sync.Mutex

	for {
		select {
		case order:= <-recieve_chan:
			assigned_elevator := orders.ElevIDToElevStruct(order.AssignedElevator, active_elevators)
			order_completion_duration := orders.TimeToRequestHandled(&assigned_elevator, &order)
			timer := time.NewTimer(time.Duration(order_completion_duration + tolarance_duration))
			<-timer.C
			if order.OrderConfirmation != orders.COMPLETED {
				// Assign again discarding this elevator
				if len(active_elevators) < 2 {
					break
				}
				mu.Lock()
				active_elevators = RemoveElevFromMap(&assigned_elevator, active_elevators)
				orders.AssignOrderToElevator(&order, active_elevators)
				active_elevators = AddElevToMap(&assigned_elevator, active_elevators)
				mu.Unlock()
				order_reassigned <- order
			}

		case order := <-order_reassigned:
			assigned_elevator := orders.ElevIDToElevStruct(order.AssignedElevator, active_elevators)
			order_completion_duration := orders.TimeToRequestHandled(&assigned_elevator, &order)
			timer := time.NewTimer(time.Duration(order_completion_duration + tolarance_duration))
			<-timer.C
			if order.OrderConfirmation != orders.COMPLETED {
				// Assign again discarding this elevator
				if len(active_elevators) < 2 {
					break
				}
				mu.Lock()
				active_elevators = RemoveElevFromMap(&assigned_elevator, active_elevators)
				orders.AssignOrderToElevator(&order, active_elevators)
				active_elevators = AddElevToMap(&assigned_elevator, active_elevators)
				mu.Unlock()
				order_reassigned <- order
			}
		}
	}
}

func RemoveElevFromMap(elev *ec.Elevator, elev_map map[string]ec.Elevator) map[string]ec.Elevator {
	delete(elev_map, elev.ElevID)
	return elev_map
}

func AddElevToMap(elev *ec.Elevator, elev_map map[string]ec.Elevator) map[string]ec.Elevator {
	elev_map[elev.ElevID] = *elev
	return elev_map
}