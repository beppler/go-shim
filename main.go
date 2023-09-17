package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

func main() {
	debugging := os.Getenv("GO_SHIM_DEBUG") == "1"

	executable, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if debugging {
		command := os.Getenv("GO_SHIM_DEBUG_COMMAND")
		if command != "" {
			fmt.Fprintf(os.Stderr, "command replaced to: %s from: %s\n", command, executable)
			executable = command
		}
	}

	configFile := strings.TrimSuffix(executable, filepath.Ext(executable)) + ".ini"

	config, err := ini.ShadowLoad(configFile)
	if err != nil {
		if debugging {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		os.Exit(1)
	}

	cmdPath := strings.TrimSpace(config.Section("").Key("command").String())
	if cmdPath == "" {
		fmt.Fprintf(os.Stderr, "error: command is not defined on file %s.\n", configFile)
		os.Exit(1)
	}

	if !filepath.IsAbs(cmdPath) {
		cmdDir, _ := filepath.Split(cmdPath)
		if cmdDir != "" {
			exeDir, _ := filepath.Split(executable)
			cmdPath, err = filepath.Abs(filepath.Join(exeDir, cmdPath))
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err)
				os.Exit(1)
			}
		}
	}

	wait := true
	if config.Section("").HasKey("wait") {
		wait, err = config.Section("").Key("wait").Bool()
		if err != nil {
			if debugging {
				fmt.Fprintf(os.Stderr, "error: can not determine the value of debug %s: %v", configFile, err)
			}
			os.Exit(1)
		}
	}

	var cmdArgs []string
	if config.Section("").Key("args").Value() == "" {
		cmdArgs = os.Args[1:]
	} else {
		cmdArgs = append(config.Section("").Key("args").ValueWithShadows(), os.Args[1:]...)
	}

	if debugging {
		fmt.Fprintf(os.Stderr, "executing: %s %s\n", cmdPath, strings.Join(cmdArgs, " "))
	}

	code := 0

	err = RunProcess(cmdPath, cmdArgs, wait)
	if err != nil {
		code = 1
		if exit, ok := err.(*exec.ExitError); ok {
			code = exit.ExitCode()
		} else {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
	}

	if wait && debugging {
		fmt.Fprintf(os.Stderr, "exit code: %d", code)
	}

	os.Exit(code)
}
