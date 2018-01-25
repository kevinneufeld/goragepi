package door

import (
	"fmt"
	"time"
	"github.com/stianeikeland/go-rpio"
)

func CheckDoorStatus(pinNumber int) (state string, err error) {
	fmt.Println("CheckDoorStatus: Start")
	err = rpio.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rpio.Close()

	pin := rpio.Pin(pinNumber)

	status := "open"
	if pin.Read() == 0 {
		status = "closed"
	}
	fmt.Println("CheckDoorStatus: %s", status)
	return status, err
}

func ToggleSwitch(pinNumber int, sleepTimeout int) (err error) {
	err = rpio.Open()
	if err != nil {
		return err
	}
	pin := rpio.Pin(pinNumber)
	pin.Output()

	pin.Low()
	rpio.Close()

	snooze := time.Duration(sleepTimeout) * time.Millisecond
	time.Sleep(snooze)

	err = rpio.Open()
	if err != nil {
		return err
	}
	pin.High()
	rpio.Close()

	return nil
}
