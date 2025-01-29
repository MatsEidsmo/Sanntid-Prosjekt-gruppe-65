package main

import "Driver-go/elevio"
//import "fmt"

func main() {
	numFloors := 4

    elevio.Init("localhost:15657", numFloors)
    
    var d elevio.MotorDirection = elevio.MD_Down
    //var b elevio.ButtonType = elevio.BT_Cab
    elevio.SetMotorDirection(d)

    drv_floors := make(chan int)
    drv_buttons := make(chan elevio.ButtonEvent)

    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollButtons(drv_buttons)

    elevio.SetButtonLamp(elevio.BT_Cab, 0, false)
    elevio.SetButtonLamp(elevio.BT_Cab, 1, false)
    elevio.SetButtonLamp(elevio.BT_Cab, 2, false)
    elevio.SetButtonLamp(elevio.BT_Cab, 3, false)

    for {
        select{
        case floornumber := <-drv_floors:
            elevio.SetFloorIndicator(floornumber)

            if floornumber == 2 {
                elevio.SetMotorDirection(elevio.MD_Stop)
            }

        case button := <-drv_buttons:
            elevio.SetButtonLamp(button.Button, button.Floor, false)
        }

    }

}