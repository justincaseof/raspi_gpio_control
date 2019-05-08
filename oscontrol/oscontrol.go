package oscontrol

import (
	"go.uber.org/zap"
	"raspi_gpio_control/logging"
)

var logger = logging.New("raspi_gpio_control_oscontrol", false)

func RestartOS() error {
	logger.Debug("  --> executing RESTART command...")
	err := RestartOSnative()
	if err != nil {
		logger.Error("Error executing restart command", zap.Error(err))
	}
	logger.Debug("...done.")
	return err
}

func PoweroffOS() error {
	logger.Debug("  --> executing POWEROFF command...")
	err := PoweroffOSnative()
	if err != nil {
		logger.Error("Error executing poweroff command", zap.Error(err))
	}
	logger.Debug("...done.")
	return err
}
