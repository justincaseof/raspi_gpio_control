package gpiocontrol

// GPIOConfig cfg
type GPIOConfig struct {
	RestartPin 	string `yaml:"restart-pin"`
	PoweroffPin string `yaml:"poweroff-pin"`
}

func InitGPIO(gpioConfig *GPIOConfig) error {
	return InitGPIONative(gpioConfig)
}