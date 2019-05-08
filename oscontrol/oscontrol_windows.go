package oscontrol

import (
	"os/exec"
	"context"
	"fmt"
	"time"
)

func RestartOSnative() error {
	runCommand("shutdown", "/r", "/t 0")
	return nil
}

func PoweroffOSnative() error {
	runCommand("shutdown", "/s", "/t 0")
	return nil
}

func runCommand(command string, arg1 string, arg2 string) {
	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Create the command with our context
	cmd := exec.CommandContext(ctx, command, arg1, arg2)

	// This time we can simply use Output() to get the result.
	out, err := cmd.Output()

	// We want to check the context error to see if the timeout was executed.
	// The error returned by cmd.Output() will be OS specific based on what
	// happens when a process is killed.
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Command timed out")
		return
	}

	// If there's no context error, we know the command completed (or errored).
	fmt.Println("Output:", string(out))
	if err != nil {
		fmt.Println("Non-zero exit code:", err)
	}

}