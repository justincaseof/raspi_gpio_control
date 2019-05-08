package gpiocontrol

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	gpioperiph "periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/physic"
	iohost "periph.io/x/periph/host"
)

var restartPin gpioperiph.PinIn
var poweroffPin gpioperiph.PinIn
var ledPin gpioperiph.PinOut

func initGPIONative(gpioConfig *GPIOConfig) error {
	logger.Info("* initializing gpio lib *")
	if _, err := iohost.Init(); err != nil {
		logger.Error("error initializing gpio lib", zap.Error(err))
		return errors.New("error initializing gpio lib")
	}
	logger.Info("done.")

	logger.Info("* setting up GPIO pins *")

	// ### RESTART PIN ###
	logger.Info("\t--> initializing pin:", zap.String("restartPin", gpioConfig.RestartPin))
	restartPin = gpioreg.ByName(gpioConfig.RestartPin)
	if restartPin == nil {
		return errors.New("unable to set up restartPin")
	}
	// we're using 'RisingEdge' to trigger interrupt upon release of pushed button
	restartPin.In(gpioperiph.PullUp, gpioperiph.RisingEdge)
	logger.Info("\t--> done.")

	// ### POWEROFF PIN ###
	logger.Info("\t--> initializing pin:", zap.String("poweroffPin", gpioConfig.PoweroffPin))
	poweroffPin = gpioreg.ByName(gpioConfig.PoweroffPin)
	if poweroffPin == nil {
		return errors.New("unable to set up poweroffPin")
	}
	// we're using 'RisingEdge' to trigger interrupt upon release of pushed button
	poweroffPin.In(gpioperiph.PullUp, gpioperiph.RisingEdge)
	logger.Info("\t--> done.")

	// ### LED PIN ###
	logger.Info("\t--> initializing pin:", zap.String("ledPin", gpioConfig.LEDPin))
	ledPin = gpioreg.ByName(gpioConfig.LEDPin)
	if ledPin == nil {
		return errors.New("unable to set up ledPin")
	}
	// we're using 'RisingEdge' to trigger interrupt upon release of pushed button
	//ledPin.PWM(gpioperiph.DutyHalf, 10 * physic.Hertz)
	//ledPin.Out(gpioperiph.High)
	//ledPin.Out(gpioperiph.Low)
	logger.Info("\t--> done.")

	return nil
}

func hasInterruptRESTART() bool {
	// WaitForEdge is blocking
	return restartPin.WaitForEdge(-1)
}

func hasInterruptPOWEROFF() bool {
	// WaitForEdge is blocking
	return poweroffPin.WaitForEdge(-1)
}

var isOn bool = true

func toggleLEDnative() {
	if isOn {
		isOn = false
		logger.Info("  LED: OFF")
		ledPin.Out(gpioperiph.Low)
	} else {
		isOn = true
		logger.Info("  LED: ON")
		ledPin.Out(gpioperiph.High)
	}
}

func pinPWMnative(dutyPercentage uint32, hertz uint64) {
	ledPin.PWM(gpioperiph.Duty(uint32(gpioperiph.DutyMax)*dutyPercentage/100), physic.Frequency(hertz*uint64(physic.Hertz)))
}
func pinHIGHnative() {
	ledPin.Out(gpioperiph.High)
}
func pinLOWnative() {
	ledPin.Out(gpioperiph.Low)
}
