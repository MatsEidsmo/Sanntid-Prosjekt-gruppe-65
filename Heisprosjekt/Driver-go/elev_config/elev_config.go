package elev_config

const N_floors int = 4
const N_buttons int = 3

type Elevator struct {
	Floor int
	Dir int
	RequestMatrix [N_floors][N_buttons]int
}