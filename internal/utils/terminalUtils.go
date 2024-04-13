package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
)

func RunCommand(commandStr string) error {
	var cmd *exec.Cmd
	var stdoutBuf, stderrBuf bytes.Buffer

	// Check the operating system
	switch runtime.GOOS {
	case "linux", "darwin": // darwin is macOS
		cmd = exec.Command("/bin/sh", "-c", commandStr)
	case "windows":
		cmd = exec.Command("cmd", "/c", commandStr)
	default:
		return fmt.Errorf("unsupported platform")
	}

	// Capture the standard output and error
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Start the command and wait for it to finish
	if err := cmd.Start(); err != nil {
		return err
	}

	err := cmd.Wait()

	if err != nil {
		return fmt.Errorf("command failed with error: %v\nstdout: %s\nstderr: %s", err, stdoutBuf.String(), stderrBuf.String())
	}

	fmt.Println(stdoutBuf.String())

	return nil
}

func IsBrewInstalled() bool {
	cmd := exec.Command("brew", "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func InstallBrew() error {
	var cmd *exec.Cmd

	// Determine the OS and set up the appropriate installation command
	switch runtime.GOOS {
	case "darwin":
		// macOS installation command
		cmd = exec.Command("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)")
	case "linux":
		// Linux installation command
		cmd = exec.Command("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)")
	default:
		return fmt.Errorf("Homebrew installation is not supported on your OS: %s", runtime.GOOS)
	}

	// Run the installation command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install Homebrew: %s, %s", err, output)
	}

	fmt.Println("Homebrew installed successfully")
	return nil
}
