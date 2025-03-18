package heartbeat

import (
	"Driver-go/network/bcast"
	"Driver-go/elev_config"
	"fmt"
	"time"
)

type Heartbeat struct {
	Elevator  elev_config.Elevator
	Timestamp time.Time
}

func Transmitter(elevator elev_config.Elevator, txChan chan Heartbeat) {
	for {
		txChan <- Heartbeat{Elevator: elevator, Timestamp: time.Now()}
		time.Sleep(500 * time.Millisecond) // Send heartbeat every 500ms
	}
}

func Receiver(rxChan chan Heartbeat, activeElevators map[string]Heartbeat) {
	for {
		heartbeat := <-rxChan
		activeElevators[heartbeat.Elevator.ElevID] = heartbeat
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
