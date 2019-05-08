package oscontrol

func RestartOS() error {
	return RestartOSnative()
}

func PoweroffOS() error {
	return PoweroffOSnative()
}
