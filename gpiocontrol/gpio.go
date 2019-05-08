package gpiocontrol

import (
	"raspi_gpio_control/logging"
)

// GPIOConfig cfg
type GPIOConfig struct {
	RestartPin  string `yaml:"restart-pin"`
	PoweroffPin string `yaml:"poweroff-pin"`
	LEDPin      string `yaml:"LED-pin"`
}

type Interrupt uint8

const (
	InterruptNONE     Interrupt = 0
	InterruptRESTART  Interrupt = 1
	InterruptPOWEROFF Interrupt = 2
)

var logger = logging.New("raspi_gpio_control_base", false)

func InitGPIO(gpioConfig *GPIOConfig) error {
	return initGPIONative(gpioConfig)
}

func CheckInterruptRESTART(interruptChannel chan Interrupt, processing *bool) {
	for {
		if hasInterruptRESTART() {
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
		if hasInterruptPOWEROFF() {
			if *processing {
				logger.Debug("discarding POWEROFF interrupt")
			} else {
				interruptChannel <- InterruptPOWEROFF
			}
		}
	}
}

func LEDpwm(dutyPercentage uint, hertz uint) {
	pinPWMnative(dutyPercentage, hertz)
}

func LEDon() {
	pinHIGHnative()
}

func LEDoff() {
	pinLOWnative()
}
