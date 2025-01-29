package Elevator

const n_floors int = 4
const n_buttons int = 3

type Elevator struct {
	floor int
	dir int
	requestMatrix [n_floors][n_buttons]int
}