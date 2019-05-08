package oscontrol

import "os/exec"

func RestartOSnative() error {
	return exec.Command("/usr/bin/systemctl", "restart").Run()

}

func PoweroffOSnative() error {
	return exec.Command("/usr/bin/systemctl", "poweroff").Run()
}
