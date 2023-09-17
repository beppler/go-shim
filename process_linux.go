package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func RunProcess(name string, args []string, wait bool) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if wait {
		cmd.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGKILL}

		if err := cmd.Run(); err != nil {
			if exit, ok := err.(*exec.ExitError); ok {
				return exit
			}
			return fmt.Errorf("error running process: %w", err)
		}
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting process: %w", err)
	}

	return nil
}
