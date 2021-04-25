package main

// https://gist.github.com/hallazzang/76f3970bfc949831808bbebc8ca15209#gistcomment-2948162

import (
	"os"
	"os/exec"
	"unsafe"

	"golang.org/x/sys/windows"
)

// We use this struct to retreive process handle(which is unexported)
// from os.Process using unsafe operation.
type process struct {
	Pid    int
	Handle uintptr
}

type ProcessGroup windows.Handle

func NewProcessGroup() (ProcessGroup, error) {
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

	return ProcessGroup(handle), nil
}

func (g ProcessGroup) Dispose() error {
	return windows.CloseHandle(windows.Handle(g))
}

func (g ProcessGroup) AddProcess(p *os.Process) error {
	return windows.AssignProcessToJobObject(windows.Handle(g), windows.Handle((*process)(unsafe.Pointer(p)).Handle))
}

func (g ProcessGroup) SetupCommand(cmd *exec.Cmd) {

}
