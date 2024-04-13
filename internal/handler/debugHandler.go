package handler

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/briandowns/spinner"

	"github.com/rishijash/idk_terminal/configs"
	"github.com/rishijash/idk_terminal/internal/clients"
	"github.com/rishijash/idk_terminal/internal/utils"
)

type DebugHandler struct {
	config *configs.Config
}

func NewDebugHandler(config *configs.Config) DebugHandler {
	return DebugHandler{
		config: config,
	}
}

func (h DebugHandler) HandleDebugMode(ctx context.Context) {
	token, err := utils.LoadToken()
	if err != nil {
		println("You are not logged in. Please login first")
		println("Command: `idk --login`")
		return
	}

	fmt.Println("IDK Debug mode enabled")
	fmt.Println("Before executing any command, IDK will give you context about it")
	fmt.Println("Help")
	fmt.Println("------------------------------------------")
	fmt.Println("• exit	:	exit debug mode")
	fmt.Println("------------------------------------------")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("IDK DEBUG MODE > ")
		if !scanner.Scan() {
			break // Handles EOF or read error
		}
		command := scanner.Text()
		if command == "exit" {
			break // User exits the proxy shell
		}

		// Execute the command via a shell
		// Here 'bash' is used, but you can replace it with 'sh' or any other shell
		cmd := exec.Command("bash", "-c", command)

		// Create buffer to capture standard output and standard error
		var outBuffer, errBuffer bytes.Buffer
		cmd.Stdout = &outBuffer
		cmd.Stderr = &errBuffer
		cmd.Stdin = os.Stdin

		// Run the command
		err := cmd.Run()

		// Check for errors and output them
		if err != nil {
			fmt.Printf(errBuffer.String())
			detailedErr :=
				fmt.Errorf("command failed with error: %s", errBuffer.String())
			h.commandDebugAction(command, detailedErr, token)
		} else {
			fmt.Printf(outBuffer.String())
		}
	}

	fmt.Println("Exiting from IDK Debug Mode")
}

func (h DebugHandler) commandDebugAction(command string, err error, token string) {
	fmt.Println("Analyzing Error..")
	customCharset := []string{"-", "\\", "|", "/", "-", ".", "o", "O", "0", "@"}
	loadingSpinner := spinner.New(customCharset, 100*time.Millisecond)
	loadingSpinner.Start()

	debugResponse, err, responseStatus := clients.ProcessDebugCommand(command, runtime.GOOS, err, token, h.config.IdkBackendBaseUrl)
	loadingSpinner.Stop()

	if responseStatus == http.StatusUnauthorized {
		utils.ClearToken()
		fmt.Println("Token expired. Please login again")
		println("Command: `idk --login`")
		return
	}

	if responseStatus == http.StatusTooManyRequests {
		fmt.Println("Daily Quota limit reached. Plesae try again tomorrow or upgrade on https://idk-cli.github.io/")
		return
	}

	if err != nil {
		fmt.Println("Something went wrong. Please try again!")
		return
	}

	utils.PrintMessage(debugResponse.Response)
}
