package gpiocontrol

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	gpioperiph "periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"time"
        "raspi_gpio_control/logging"
)
var logger = logging.New("raspi_gpio_control_linux", false)
var restartPin gpioperiph.PinIn
var poweroffPin gpioperiph.PinIn

func InitGPIONative(gpioConfig *GPIOConfig) error {
	logger.Info("* setting up GPIO *")

	logger.Info("\t--> initializing pin:", zap.String("restartPin", gpioConfig.RestartPin))
	restartPin := gpioreg.ByName(gpioConfig.RestartPin)
	if restartPin == nil {
		return errors.New("unable to set up restartPin")
	}
	// we're using 'RisingEdge' to trigger interrupt upon release of pushed button
	restartPin.In(gpioperiph.PullUp, gpioperiph.RisingEdge)
	logger.Info("\t--> done.")

	logger.Info("\t--> initializing pin:", zap.String("poweroffPin", gpioConfig.PoweroffPin))
	poweroffPin := gpioreg.ByName(gpioConfig.PoweroffPin)
	if poweroffPin == nil {
		return errors.New("unable to set up poweroffPin")
	}
	// we're using 'RisingEdge' to trigger interrupt upon release of pushed button
	poweroffPin.In(gpioperiph.PullUp, gpioperiph.RisingEdge)
	logger.Info("\t--> done.")

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
