package elev_config

import (
	eio "Driver-go/elevio"
	"time"
)

const N_floors int = 4
const N_buttons int = 3


const DOOR_TIMEOUT = 3 * time.Second

type Elevator struct {
	Floor int
	Dir eio.MotorDirection // except MD_Stop
	RequestMatrix [N_floors][N_buttons]int
}