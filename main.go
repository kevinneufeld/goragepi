package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/log"
	"github.com/kevinneufeld/goragepi/door"
)

//Version ....
const Version = "0.1.0"

//Options ...
type Options struct {
	pin          string
	relayPin     int
	statusPin    int
	sleepTimeout int
	version      bool
}

var options Options

var serialNumber string = os.Getenv("RESIN_DEVICE_UUID")

func toggleDoor(o Options) func(int) {
	return func(targetState int) {
		nextState := "closed"
		if targetState == characteristic.TargetDoorStateOpen {
			nextState = "open"
		}

		if currentDoorState, err := door.CheckDoorSensor(o.statusPin); err != nil {
            fmt.Printf("ERROR: Could not read status pin %v\n", err)
		} else {
			if currentDoorState != nextState {
			    log.Info.Println("Toggling Relay Switch")
				door.ToggleSwitch(o.relayPin, o.sleepTimeout)
			}
		}
	}
}
// Update iOS with the correct target state.
func syncAccTargetDoorState(acc *GarageDoorOpener){
    targetState := acc.GarageDoorOpener.TargetDoorState.GetValue()
    currentState := acc.GarageDoorOpener.CurrentDoorState.GetValue()

    if targetState != currentState {
        time.Sleep(5 * time.Second)
        acc.GarageDoorOpener.TargetDoorState.SetValue(currentState)
    }
}

func pollDoorStatus(acc *GarageDoorOpener, pin int) {
    lastKnownDoorState := ""
	for {
		if currentDoorState, err := door.CheckDoorSensor(pin); err != nil {
            fmt.Printf("ERROR: Could not read currentDoorState pin %v\n", err)
		} else {
		    switch currentDoorState {
			case "closed":
				acc.GarageDoorOpener.CurrentDoorState.SetValue(characteristic.CurrentDoorStateClosed)
            default:
                acc.GarageDoorOpener.CurrentDoorState.SetValue(characteristic.CurrentDoorStateOpen)
            }
            
            if lastKnownDoorState != currentDoorState {
                if lastKnownDoorState != "" {
                    log.Info.Printf("DoorSensor: %s -> %s", lastKnownDoorState, currentDoorState)
                }else {
                    log.Info.Printf("InitSenorState: %s", currentDoorState)
                }
                lastKnownDoorState = currentDoorState
                syncAccTargetDoorState(acc)
            }
        }

		time.Sleep(time.Second)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&options.pin, "pin", "", "8-digit Pin for securing garage door")
	flag.IntVar(&options.relayPin, "relay-pin", 17, "GPIO pin of relay")
	flag.IntVar(&options.statusPin, "status-pin", 5, "GPIO pin of reed switch")
	flag.IntVar(&options.sleepTimeout, "sleep", 500, "Time in milliseconds to keep switch closed")
	flag.BoolVar(&options.version, "version", false, "print version and exit")
	flag.Parse()

	if options.version {
		fmt.Printf("garage-server-homekit v%v\n", Version)
		os.Exit(0)
	}

	if serialNumber == "" {
		println("You did not set SERIAL_NUMBER env var")
		os.Exit(1)
	}

	if options.pin == "" || len(options.pin) != 8 {
		println("Pin must be and 8 digit number")
		os.Exit(0)
	}

	info := accessory.Info{
		Name:         "Garage Door",
		Manufacturer: "Rusty Cog",
		Model:        "Raspberry Pi",
		SerialNumber: serialNumber,
	}

    log.Info.Printf("relayPin: %v \n", options.relayPin)
    log.Info.Printf("StatusPin: %v \n", options.statusPin)
    log.Info.Printf("StatusSleepInterval: %v \n", options.sleepTimeout)
	acc := NewGarageDoorOpener(info)

	acc.GarageDoorOpener.TargetDoorState.OnValueRemoteUpdate(toggleDoor(options))



	t, err := hc.NewIPTransport(hc.Config{Pin: options.pin}, acc.Accessory)
	if err != nil {
        log.Info.Panic(err)
	}

	go pollDoorStatus(acc, options.statusPin)

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()
}
