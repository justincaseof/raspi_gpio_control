package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"raspi_gpio_control/gpiocontrol"
	"raspi_gpio_control/logging"
	"syscall"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var logger = logging.New("raspi_gpio_control_main", false)
const GPIO_CONFIG_FILENAME = "gpioconfig.yml"

var interruptChannel = make(chan gpiocontrol.Interrupt)

func main() {
	logger.Info("### STARTUP")

	// INIT
	var cfg gpiocontrol.GPIOConfig
	readGPIOConfig(&cfg)
	err := gpiocontrol.InitGPIO(&cfg)
	if err != nil {
		logger.Error("Cannot set up GPIO", zap.Error(err))
		panic("Cannot set up GPIO")
	}

	// GO
	go mainLoop()
	go gpiocontrol.CheckInterruptRESTART(interruptChannel)
	go gpiocontrol.CheckInterruptPOWEROFF(interruptChannel)

	// wait indefinitely until external abortion
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // Ctrl + c
	<-sigs
	logger.Info("### EXIT")
}

// ==== I/O and properties ====

func readGPIOConfig(gpioconfig *gpiocontrol.GPIOConfig) {
	var err error
	var bytes []byte
	bytes, err = ioutil.ReadFile(GPIO_CONFIG_FILENAME)
	if err != nil {
		logger.Error("Cannot open config file", zap.String("filename", GPIO_CONFIG_FILENAME))
		panic(err)
	}
	err = yaml.Unmarshal(bytes, gpioconfig)
	if err != nil {
		panic(err)
	}
	logger.Info("GPIOConfig parsed.")
}

func mainLoop() {
	for {
		select {
		case <-time.After(5 * time.Second):
			logger.Debug("* Tick *")
		case interrupt := <-interruptChannel:
			logger.Debug("INTERRUPT!", zap.Uint8("interrupt", uint8(interrupt)))
		}
	}
}
