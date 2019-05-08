package gpiocontrol

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

func InitGPIO(gpioConfig *GPIOConfig) error {
	return InitGPIONative(gpioConfig)
}

func CheckInterruptRESTART(interruptChannel chan Interrupt) {
	for {
		if HasInterruptRESTART() {
			interruptChannel <- InterruptRESTART
		}
	}
}

func CheckInterruptPOWEROFF(interruptChannel chan Interrupt) {
	for {
		if HasInterruptPOWEROFF() {
			interruptChannel <- InterruptPOWEROFF
		}
	}
}