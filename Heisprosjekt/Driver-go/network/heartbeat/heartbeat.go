package heartbeat

import (
	ec "Driver-go/elev_config"
	"fmt"
	"time"
	//"Driver-go/orders"
	// eio "Driver-go/elevio"
	// bcast "Driver-go/network/bcast"
	// counter "Driver-go/network/counter"
)

type Heartbeat struct {
	Elevator  ec.Elevator
	Timestamp time.Time
}

func Transmitter(elevator ec.Elevator, txChan chan Heartbeat) {
	for {
		txChan <- Heartbeat{Elevator: elevator, Timestamp: time.Now()}
		time.Sleep(2000 * time.Millisecond) 
		
		
	}
}

func Receiver(rxChan chan Heartbeat, activeElevators map[string]Heartbeat) {
	for {
		heartbeat := <-rxChan
		activeElevators[heartbeat.Elevator.ElevID] = heartbeat
		//fmt.Println("Recieved Heartbeat from:", string(heartbeat.Elevator.ElevID))
		//fmt.Println(activeElevators)
		//fmt.Println(len(activeElevators))
	}
}

func RemoveInactiveElevators(activeElevators map[string]Heartbeat, timeout time.Duration) {
	for {
		time.Sleep(timeout)
		now := time.Now()
		for id, hb := range activeElevators {
			if now.Sub(hb.Timestamp) > timeout {
				fmt.Printf("Elevator %s lost connection\n", id)
				delete(activeElevators, id)
			}
		}
	}
}

// func UpdateActiveElevators(recieved_update chan *ec.Elevator, ) {
// 	for {
// 		select {
// 		case e := <-recieved_update :
// 		}
// 	}
// }





