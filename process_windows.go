package main

// https://gist.github.com/hallazzang/76f3970bfc949831808bbebc8ca15209#gistcomment-2948162

import (
	"fmt"
	"os"
	"os/exec"
)

func RunProcess(name string, args []string, wait bool) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if wait {
		group, err := NewJobObject()
		if err != nil {
			return fmt.Errorf("error creating process group: %w", err)
		}
		defer group.Dispose()

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("error starting process: %w", err)
		}

		if err := group.AssignProcess(cmd.Process); err != nil {
			return fmt.Errorf("error assigning process group: %w", err)
		}

		if err := cmd.Wait(); err != nil {
			if exit, ok := err.(*exec.ExitError); ok {
				return exit
			}
			return fmt.Errorf("error waiting for process: %w", err)
		}

		return nil
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting process: %w", err)
	}

	return nil
}
