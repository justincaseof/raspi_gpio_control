package gpiocontrol

import (
	"periph.io/x/periph/conn/gpio"
	"time"
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

func HasInterruptRESTART() bool {
	// WaitForEdge is blocking
	return restartPin.WaitForEdge(time.Duration(-1))
}

func HasInterruptPOWEROFF() bool {
	// WaitForEdge is blocking
	return poweroffPin.WaitForEdge(time.Duration(-1))
}