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
