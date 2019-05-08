package oscontrol

import (
	"go.uber.org/zap"
	"raspi_gpio_control/logging"
)

var logger = logging.New("raspi_gpio_control_oscontrol", false)

var isSimulation bool = false

func Init(simulateOScontrol bool) {
	isSimulation = simulateOScontrol
}

func RestartOS() error {
	logger.Debug("  --> executing RESTART command...")
	if isSimulation {
		logger.Debug("  <<RESTART-SIMULATION>> ")
		return nil
	}
	err := restartOSnative()
	if err != nil {
		logger.Error("Error executing restart command", zap.Error(err))
	}
	logger.Debug("...done.")
	return err
}

func PoweroffOS() error {
	logger.Debug("  --> executing POWEROFF command...")
	if isSimulation {
		logger.Debug("  <<POWEROFF-SIMULATION>> ")
		return nil
	}
	err := poweroffOSnative()
	if err != nil {
		logger.Error("Error executing poweroff command", zap.Error(err))
	}
	logger.Debug("...done.")
	return err
}
