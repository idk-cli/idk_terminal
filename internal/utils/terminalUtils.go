package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

func GetTerminalHistory() ([]string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	homeDir := currentUser.HomeDir

	// Attempt to read Bash history
	bashHistoryPath := homeDir + "/.bash_history"
	bashHistory, err := readHistoryFile(bashHistoryPath)
	if err == nil {
		return bashHistory, nil
	}

	// Attempt to read Zsh history
	zshHistoryPath := homeDir + "/.zsh_history"
	zshHistory, err := readHistoryFile(zshHistoryPath)
	if err != nil {
		return nil, err
	}

	return zshHistory, nil
}

func readHistoryFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var commands []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// For Zsh, you might need to process lines to remove timestamps if enabled in history
		line := scanner.Text()
		if strings.HasPrefix(line, ":") { // Zsh history with metadata
			if split := strings.SplitN(line, ";", 2); len(split) == 2 {
				line = split[1]
			}
		}
		commands = append(commands, line)
	}

	return commands, scanner.Err()
}

func ResolveCDPathCommands(commands []string) []string {
	var result []string
	currentPath := ""

	for _, command := range commands {
		parts := strings.Fields(command)
		if len(parts) != 2 {
			continue // Skip invalid commands
		}
		directory := parts[1]

		if strings.HasPrefix(directory, "/") {
			// Absolute path
			currentPath = directory
			result = append(result, command)
		} else if directory == ".." {
			// Navigate up one directory
			if currentPath != "" {
				lastIndex := strings.LastIndex(currentPath, "/")
				if lastIndex > 0 {
					currentPath = currentPath[:lastIndex]
				} else {
					// Back to root
					currentPath = ""
				}
			}
			if currentPath != "" {
				result = append(result, "cd "+currentPath)
			}
		} else {
			// Relative path
			if currentPath == "" || currentPath == "/" {
				currentPath = "/" + directory
			} else {
				currentPath = currentPath + "/" + directory
			}
			result = append(result, "cd "+currentPath)
		}
	}

	// filter out paths that don't exist
	var validPathCommands []string
	for _, resultCommand := range result {
		path := strings.Replace(resultCommand, "cd ", "", 1)
		if _, err := os.Stat(path); err == nil {
			// The path exists
			validPathCommands = append(validPathCommands, fmt.Sprintf("cd %s", path))
		}
	}

	return validPathCommands
}

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
