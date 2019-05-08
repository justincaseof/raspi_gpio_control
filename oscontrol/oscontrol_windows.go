package oscontrol

import "os/exec"

func RestartOSnative() error {
	// FIXME: never tested. check arguments syntax
	return exec.Command("shutdown", "/r /t 0").Run()

}

func PoweroffOSnative() error {
	// FIXME: never tested. check arguments syntax
	return exec.Command("shutdown", "/s /t 0").Run()
}
