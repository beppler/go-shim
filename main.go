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
	executable, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	configFile := strings.TrimSuffix(executable, filepath.Ext(executable)) + ".ini"

	config, err := ini.ShadowLoad(configFile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	cmdPath := strings.TrimSpace(config.Section("").Key("command").String())
	if cmdPath == "" {
		fmt.Printf("Error: command is not defined on file %s.\n", configFile)
		os.Exit(1)
	}

	if !filepath.IsAbs((cmdPath)) {
		cmdDir, _ := filepath.Split(cmdPath)
		if cmdDir != "" {
			exeDir, _ := filepath.Split(executable)
			cmdPath, err = filepath.Abs(filepath.Join(exeDir, cmdPath))
			if err != nil {
				fmt.Println("Error:", err)
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

	group, err := NewProcessGroup()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer group.Dispose()

	cmd := exec.Command(cmdPath, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	group.SetupCommand(cmd)

	if err = cmd.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if err = group.AddProcess(cmd.Process); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if err = cmd.Wait(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			os.Exit(exit.ExitCode())
		}
		os.Exit(1)
	}
}
