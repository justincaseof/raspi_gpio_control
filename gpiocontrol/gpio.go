package gpiocontrol

import (
	"raspi_gpio_control/logging"
)

// GPIOConfig cfg
type GPIOConfig struct {
	RestartPin 	string `yaml:"restart-pin"`
	PoweroffPin string `yaml:"poweroff-pin"`
}

type Interrupt uint8
const (
	InterruptNONE 		Interrupt = 0
	InterruptRESTART 	Interrupt = 1
	InterruptPOWEROFF 	Interrupt = 2
)

var logger = logging.New("raspi_gpio_control_base", false)

func InitGPIO(gpioConfig *GPIOConfig) error {
	return InitGPIONative(gpioConfig)
}

func CheckInterruptRESTART(interruptChannel chan Interrupt, processing *bool) {
	for {
		if HasInterruptRESTART() {
			if *processing {
				logger.Debug("discarding RESTART interrupt")
			} else {
				interruptChannel <- InterruptRESTART
			}
		}
	}
}

func CheckInterruptPOWEROFF(interruptChannel chan Interrupt, processing *bool) {
	for {
		if HasInterruptPOWEROFF() {
			if *processing {
				logger.Debug("discarding POWEROFF interrupt")
			} else {
				interruptChannel <- InterruptPOWEROFF
			}
		}
	}
}