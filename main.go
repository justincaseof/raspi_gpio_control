package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"raspi_gpio_control/gpiocontrol"
	"raspi_gpio_control/logging"
	"raspi_gpio_control/oscontrol"
	"syscall"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var logger = logging.New("raspi_gpio_control_main", false)
const GPIO_CONFIG_FILENAME = "gpioconfig.yml"

var interruptChannel = make(chan gpiocontrol.Interrupt)
var processing bool

func main() {
	logger.Info("### STARTUP")

	// INIT
	oscontrol.Init(true)
	var cfg gpiocontrol.GPIOConfig
	readGPIOConfig(&cfg)
	err := gpiocontrol.InitGPIO(&cfg)
	if err != nil {
		logger.Error("Cannot set up GPIO", zap.Error(err))
		panic("Cannot set up GPIO")
	}
	processing = false

	// GO
	go mainLoop()
	go gpiocontrol.CheckInterruptRESTART(interruptChannel, &processing)
	go gpiocontrol.CheckInterruptPOWEROFF(interruptChannel, &processing)

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
		case <-time.After(1 * time.Second):
			{
				handleState()
			}
		case interrupt := <-interruptChannel:
			{
				processing = true
				logger.Debug("INTERRUPT!", zap.Uint8("interrupt", uint8(interrupt)))
				switch interrupt {
				case gpiocontrol.InterruptRESTART:
					if state == RESTART_WAITING_CONFIMRATION {
						// give our state handler the indicatino to actually restart
						state = RESTART_COMMAND_EXECUTE
					} else {
						state = RESTART_REQUESTED
					}
					time.Sleep(500 * time.Millisecond)	// cheap way of forcing the user to wait for another button press.
					processing = false
				case gpiocontrol.InterruptPOWEROFF:
					if state == POWEROFF_WAITING_CONFIMRATION {
						// give our state handler the indicatino to actually poweroff
						state = POWEROFF_COMMAND_EXECUTE
					} else {
						state = POWEROFF_REQUESTED
					}
					time.Sleep(500 * time.Millisecond)	// cheap way of forcing the user to wait for another button press.
					processing = false
				default:
					logger.Warn("Unknown interrupt")
				}
			}
		case <-resetChannel:
			{
				logger.Debug("Resetting State.")
				state = IDLE_RUNNING
			}
		}
	}
}

type State uint8

const (
	IDLE_RUNNING 					= State(iota)

	RESTART_REQUESTED 				= State(iota)
	RESTART_WAITING_CONFIMRATION 	= State(iota)
	RESTART_COMMAND_EXECUTE			= State(iota)

	POWEROFF_REQUESTED 				= State(iota)
	POWEROFF_WAITING_CONFIMRATION 	= State(iota)
	POWEROFF_COMMAND_EXECUTE		= State(iota)

)

var state = IDLE_RUNNING
var resetChannel = make(chan bool)
func handleState() {
	switch state {
			case IDLE_RUNNING:
			gpiocontrol.LEDpwm(20, 1)
	case RESTART_REQUESTED:
		{
			logger.Debug("RESTART requested, waiting for confirmation")
			state = RESTART_WAITING_CONFIMRATION
			gpiocontrol.LEDpwm(50, 8)
			// start timeout
			go func() {
				time.Sleep(3 * time.Second)
				resetChannel <- true
			}()
		}
	case POWEROFF_REQUESTED:
		{
			logger.Debug("POWEROFF requested, waiting for confirmation")
			state = POWEROFF_WAITING_CONFIMRATION
			gpiocontrol.LEDpwm(50, 8)
			// start timeout
			go func() {
				time.Sleep(3 * time.Second)
				resetChannel <- true
			}()
		}
	case RESTART_COMMAND_EXECUTE:
		{
			state = IDLE_RUNNING
			oscontrol.RestartOS()
		}
	case POWEROFF_COMMAND_EXECUTE:
		{
			state = IDLE_RUNNING
			oscontrol.PoweroffOS()
		}
	}
}