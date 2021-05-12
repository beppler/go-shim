package main

// https://gist.github.com/hallazzang/76f3970bfc949831808bbebc8ca15209#gistcomment-2948162

import (
	"fmt"
	"os"
	"os/exec"
	"unsafe"

	"golang.org/x/sys/windows"
)

func RunProcess(name string, wait bool, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	group, err := newProcessGroup()
	if err != nil {
		return fmt.Errorf("error creating process group: %v", err)
	}
	defer group.Dispose()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting process: %v", err)
	}

	if wait {
		group.AddProcess(cmd.Process)

		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("error waiting for the process: %v", err)
		}
	}

	return nil
}

// We use this struct to retreive process handle(which is unexported)
// from os.Process using unsafe operation.
type process struct {
	Pid    int
	Handle uintptr
}

type processGroup windows.Handle

func newProcessGroup() (processGroup, error) {
	handle, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return 0, err
	}

	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}
	if _, err := windows.SetInformationJobObject(
		handle,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info))); err != nil {
		return 0, err
	}

	return processGroup(handle), nil
}

func (g processGroup) Dispose() error {
	return windows.CloseHandle(windows.Handle(g))
}

func (g processGroup) AddProcess(p *os.Process) error {
	return windows.AssignProcessToJobObject(windows.Handle(g), windows.Handle((*process)(unsafe.Pointer(p)).Handle))
}
