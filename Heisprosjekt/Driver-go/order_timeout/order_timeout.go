package order_timeout

import (
	"time"
	orders "Driver-go/orders"
	ec "Driver-go/elev_config"
	//"sync"
	"fmt"
	ea "Driver-go/elev_actuator"
	
)

// Problem: Potentially active_elevators is not updated. So if the assigned_elevator is the only elevator then we do cannot assign order to
// another elevator. But if an elevator enters the network after this decision, yet the assigned_elevator is not free to go, no reassignments are
// made leading to the order being lost. Similarly if after the decision the assigned_elev disconnects without completing the order, the order is
// lost. 
// Solutions/Thoughts: this function has to be a goroutine to always keep track, active_elevators could be a global variable to always have the most
// accessable available
func OrderTimeout(
	recieve_chan chan orders.Order,
	TransmitOrderChan chan orders.Order, 
	active_elevators map[string]ec.Elevator, 
	Order_timeout time.Duration,
	) {

	//var mu sync.Mutex

	
	
	for {
		select {
		case <- ea.OrderTimer.C:

			fmt.Println("Timer Closed!")



			// 	for _, o := range orders.MyWorldView {
			// 		if o == nil {
			// 			return
			// 		}
			// 		if time.Since(o.OrderID)> ec.ORDER_TIMEOUT {
			// 			fmt.Println("OTO Assigning to someone else")
			// 			//mu.Lock()
			// 			// fmt.Println("1")
			// 			// assigned_elevator := orders.ElevIDToElevStruct(o.AssignedElevator, active_elevators)
			// 			// fmt.Println("2")
			// 			// active_elevators = RemoveElevFromMap(&assigned_elevator, active_elevators)
			// 			// fmt.Println("3")
			// 			orders.AssignOrderToElevator(o, active_elevators)
			// 			// fmt.Println("4")
			// 			// active_elevators = AddElevToMap(&assigned_elevator, active_elevators)
			// 			// //mu.Unlock()
			// 			// fmt.Println("5")
			// 		}
			// }


		case order := <-recieve_chan:
				fmt.Println("Inside Ordertimeout", order)
		}
			// timer := time.NewTimer(fire_duration)
			// <-timer.C
			// fmt.Println("Before if")
			// if order.OrderConfirmation != orders.COMPLETED {
			// 	fmt.Println("IFF")
			// 	// Assign again discarding this elevator
			// 	// if len(active_elevators) < 2 {
			// 	// 	break
			// 	// }
			// 	fmt.Println("OTO Assigning to someone else")
			// 	//mu.Lock()
			// 	fmt.Println("1")
			// 	assigned_elevator := orders.ElevIDToElevStruct(order.AssignedElevator, active_elevators)
			// 	fmt.Println("2")
			// 	active_elevators = RemoveElevFromMap(&assigned_elevator, active_elevators)
			// 	fmt.Println("3")
			// 	orders.AssignOrderToElevator(&order, active_elevators)
			// 	fmt.Println("4")
			// 	//active_elevators = AddElevToMap(&assigned_elevator, active_elevators)
			// 	//mu.Unlock()
			// 	fmt.Println("5")
			// 	//TransmitOrderChan <- order
			// }

		// case order := <-order_reassigned:
		// 	timer := time.NewTimer(time.Duration(fire_duration))
		// 	<-timer.C
		// 	fmt.Println("Before if2")
		// 	if order.OrderConfirmation != orders.COMPLETED {
		// 		fmt.Println("IFF2")
		// 		// Assign again discarding this elevator
		// 		// if len(active_elevators) < 2 {
		// 		// 	break
		// 		// }
		// 		fmt.Println("OTO Assigning to someone else2")
		// 		mu.Lock()
		// 		assigned_elevator := orders.ElevIDToElevStruct(order.AssignedElevator, active_elevators)
		// 		active_elevators = RemoveElevFromMap(&assigned_elevator, active_elevators)
		// 		orders.AssignOrderToElevator(&order, active_elevators)
		// 		active_elevators = AddElevToMap(&assigned_elevator, active_elevators)
		// 		mu.Unlock()
		// 		order_reassigned <- order
		// 	}
		
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