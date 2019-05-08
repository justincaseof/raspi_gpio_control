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

const CONFIG_FILENAME = "config.yml"

var interruptChannel = make(chan gpiocontrol.Interrupt)
var processing bool

type AppConfig struct {
	IsSimulation bool                   `yaml:"is-simulation"`
	GPIOconfig   gpiocontrol.GPIOConfig `yaml:"gpio-config"`
}

var appConfig AppConfig

func main() {
	logger.Info("### STARTUP")

	// READ CONFIG
	readConfig(&appConfig)
	// INIT
	oscontrol.Init(appConfig.IsSimulation)
	err := gpiocontrol.InitGPIO(&appConfig.GPIOconfig)
	if err != nil {
		logger.Error("Cannot set up GPIO", zap.Error(err))
		panic("Cannot set up GPIO")
	}

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

func readConfig(appConfig *AppConfig) {
	var err error
	var bytes []byte
	bytes, err = ioutil.ReadFile(CONFIG_FILENAME)
	if err != nil {
		logger.Error("Cannot open config file", zap.String("filename", CONFIG_FILENAME))
		panic(err)
	}
	err = yaml.Unmarshal(bytes, appConfig)
	if err != nil {
		panic(err)
	}
	logger.Info("GPIOConfig parsed.")
}

func mainLoop() {
	processing = false
	for {
		select {
		case <-time.After(1 * time.Second):
			{
				//handleState()
			}
		case interrupt := <-interruptChannel:
			{
				processing = true
				logger.Debug("INTERRUPT!", zap.Uint8("interrupt", uint8(interrupt)))
				switch interrupt {
				case gpiocontrol.InterruptRESTART:
					if state == RESTART_WAITING_CONFIMRATION {
						// give our state handler the indication to actually restart
						state = RESTART_COMMAND_EXECUTE
					} else {
						state = RESTART_REQUESTED
					}
					handleState()
					time.Sleep(1000 * time.Millisecond) // cheap way of forcing the user to wait for another button press.
					processing = false
				case gpiocontrol.InterruptPOWEROFF:
					if state == POWEROFF_WAITING_CONFIMRATION {
						// give our state handler the indication to actually poweroff
						state = POWEROFF_COMMAND_EXECUTE
					} else {
						state = POWEROFF_REQUESTED
					}
					handleState()
					time.Sleep(1000 * time.Millisecond) // cheap way of forcing the user to wait for another button press.
					processing = false
				default:
					logger.Warn("Unknown interrupt")
				}
			}
		case <-resetChannel:
			{
				logger.Debug("Resetting State.")
				state = IDLE_RUNNING
				handleState()
			}
		}
	}
}

type State uint8

const (
	IDLE_RUNNING                  = State(iota)
	RESTART_REQUESTED             = State(iota)
	RESTART_WAITING_CONFIMRATION  = State(iota)
	RESTART_COMMAND_EXECUTE       = State(iota)
	POWEROFF_REQUESTED            = State(iota)
	POWEROFF_WAITING_CONFIMRATION = State(iota)
	POWEROFF_COMMAND_EXECUTE      = State(iota)
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
			gpiocontrol.LEDoff()
			oscontrol.RestartOS()
		}
	case POWEROFF_COMMAND_EXECUTE:
		{
			state = IDLE_RUNNING
			gpiocontrol.LEDoff()
			oscontrol.PoweroffOS()
		}
	}
}
