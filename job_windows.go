package main

import (
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

// We use this struct to retreive process handle(which is unexported)
// from os.Process using unsafe operation.
type process struct {
	Pid    int
	Handle uintptr
}

type JobObject windows.Handle

func NewJobObject() (JobObject, error) {
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

	return JobObject(handle), nil
}

func (j JobObject) Dispose() error {
	return windows.CloseHandle(windows.Handle(j))
}

func (j JobObject) AssignProcess(p *os.Process) error {
	return windows.AssignProcessToJobObject(windows.Handle(j), windows.Handle((*process)(unsafe.Pointer(p)).Handle))
}
