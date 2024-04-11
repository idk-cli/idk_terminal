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
	case "windows":
		cmd = exec.Command("cmd", "/c", commandStr)
	default:
		return fmt.Errorf("unsupported platform")
	}

	cmd.Stdin = os.Stdin   // Connect the command's standard input to the os Stdin
	cmd.Stdout = os.Stdout // Connect the command's standard output to the os Stdout
	cmd.Stderr = os.Stderr // Connect the command's standard error to the os Stderr

	// Start the command and wait for it to finish
	if err := cmd.Start(); err != nil {
		return err
	}

	err := cmd.Wait()

	return err
}
