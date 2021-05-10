package main

import (
	"os/exec"
	"syscall"
)

func RunProcess(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGKILL}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error starting process: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for the process: %v", err)
	}

	return nil
}
