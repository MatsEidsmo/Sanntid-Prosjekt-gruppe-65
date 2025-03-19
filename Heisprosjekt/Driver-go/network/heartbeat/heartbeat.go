package heartbeat

import (
	ec "Driver-go/elev_config"
	"fmt"
	"time"
)

type Heartbeat struct {
	Elevator  ec.Elevator
	Timestamp time.Time
}

func Transmitter(elevator ec.Elevator, txChan chan Heartbeat) {
	for {
		txChan <- Heartbeat{Elevator: elevator, Timestamp: time.Now()}
		time.Sleep(500 * time.Millisecond) 
		
	}
}

func Receiver(rxChan chan Heartbeat, activeElevators map[string]Heartbeat) {
	for {
		heartbeat := <-rxChan
		activeElevators[heartbeat.Elevator.ElevID] = heartbeat
		//fmt.Println("Recieved Heartbeat from:", string(heartbeat.Elevator.ElevID))
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
