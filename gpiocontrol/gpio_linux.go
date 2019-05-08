package gpio

import (
	"periph.io/x/periph/conn/gpio"
)

var restartPin PinIn
var poweroffPin PinIn

func InitGPIONative(gpioConfig *GPIOConfig) error {
	restartPin := gpioreg.ByName(gpioConfig.RestartPin)
	if restartPin == nil {
		return errors.New("unable to set up restartPin")
	}
	// we're using 'RisingEdge' to trigger interrupt upon release of pushed button
	restartPin.In(gpio.PullUp, gpio.RisingEdge)

	poweroffPin := gpioreg.ByName(gpioConfig.PoweroffPin)
	if poweroffPin == nil {
		return errors.New("unable to set up poweroffPin")
	}
	// we're using 'RisingEdge' to trigger interrupt upon release of pushed button
	poweroffPin.In(gpio.PullUp, gpio.RisingEdge)

	return nil
}