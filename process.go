// +build !windows

package main

import (
	"os"
	"os/exec"
)

type ProcessGroup int

func NewProcessGroup() (ProcessExitGroup, error) {
	return ProcessGroup(0), nil
}

func (g ProcessGroup) Dispose() error {
	return nil
}

func (g ProcessGroup) AddProcess(p *os.Process) error {
	return nil
}

func (g ProcessGroup) SetupCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
