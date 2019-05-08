package gpiocontrol

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
        iohost "periph.io/x/periph/host"
	gpioperiph "periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
        "raspi_gpio_control/logging"
)

var logger = logging.New("raspi_gpio_control_linux", false)
var restartPin gpioperiph.PinIn
var poweroffPin gpioperiph.PinIn

func InitGPIONative(gpioConfig *GPIOConfig) error {
    logger.Info("* initializing gpio lib *")
    if _, err := iohost.Init(); err != nil {
        logger.Error("error initializing gpio lib", zap.Error(err))
        return errors.New("error initializing gpio lib")
    }
    logger.Info("done.")

	logger.Info("* setting up GPIO pins *")

	logger.Info("\t--> initializing pin:", zap.String("restartPin", gpioConfig.RestartPin))
	restartPin = gpioreg.ByName(gpioConfig.RestartPin)
	if restartPin == nil {
		return errors.New("unable to set up restartPin")
	}
	// we're using 'RisingEdge' to trigger interrupt upon release of pushed button
	restartPin.In(gpioperiph.PullUp, gpioperiph.RisingEdge)
	logger.Info("\t--> done.")

	logger.Info("\t--> initializing pin:", zap.String("poweroffPin", gpioConfig.PoweroffPin))
	poweroffPin = gpioreg.ByName(gpioConfig.PoweroffPin)
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
	return restartPin.WaitForEdge(-1)
}

func HasInterruptPOWEROFF() bool {
	// WaitForEdge is blocking
	return poweroffPin.WaitForEdge(-1)
}
