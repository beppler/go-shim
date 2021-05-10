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
	isDebug := os.Getenv("GO_SHIM_DEBUG") == "1"

	executable, err := os.Executable()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	configFile := strings.TrimSuffix(executable, filepath.Ext(executable)) + ".ini"

	config, err := ini.ShadowLoad(configFile)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	cmdPath := strings.TrimSpace(config.Section("").Key("command").String())
	if cmdPath == "" {
		fmt.Printf("error: command is not defined on file %s.\n", configFile)
		os.Exit(1)
	}

	if !filepath.IsAbs(cmdPath) {
		cmdDir, _ := filepath.Split(cmdPath)
		if cmdDir != "" {
			exeDir, _ := filepath.Split(executable)
			cmdPath, err = filepath.Abs(filepath.Join(exeDir, cmdPath))
			if err != nil {
				fmt.Println("error:", err)
				os.Exit(1)
			}
		}
	}

	var cmdArgs []string
	if config.Section("").Key("args").Value() == "" {
		cmdArgs = os.Args[1:]
	} else {
		cmdArgs = append(config.Section("").Key("args").ValueWithShadows(), os.Args[1:]...)
	}

	if isDebug {
		fmt.Printf("executing: %s %s\n", cmdPath, strings.Join(cmdArgs, " "))
	}

	err = RunProcess(cmdPath, cmdArgs...)
	if err != nil {
		if isDebug {
			fmt.Println("error:", err)
		}
		if exit, ok := err.(*exec.ExitError); ok {
			os.Exit(exit.ExitCode())
		}
		os.Exit(1)
	}
}
