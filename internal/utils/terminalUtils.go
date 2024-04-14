package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func RunCommand(commandStr string) error {
	var cmd *exec.Cmd

	// Check the operating system
	switch runtime.GOOS {
	case "linux", "darwin": // darwin is macOS
		cmd = exec.Command("/bin/sh", "-c", commandStr)
	default:
		return fmt.Errorf("Unsupported platform")
	}

	cmd.Stdin = os.Stdin   // Connect the command's standard input to the os Stdin
	cmd.Stdout = os.Stdout // Connect the command's standard output to the os Stdout
	cmd.Stderr = os.Stderr // Connect the command's standard error to the os Stderr

	// Start the command and wait for it to finish
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
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
