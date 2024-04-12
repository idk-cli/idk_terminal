package handler

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
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

func (h DebugHandler) HandleCommandDebug(ctx context.Context, command string) {
	token, err := utils.LoadToken()
	if err != nil {
		println("You are not logged in. Please login first")
		println("Command: `idk --login`")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("This will execute command `%s` and analyzing the result", command)
	println("")
	fmt.Printf("Continue? (y/n)")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response) // Trim whitespace and newline character

	if strings.ToLower(response) == "y" {
		err = utils.RunCommand(command)
	} else {
		fmt.Println("Command execution canceled")
		return
	}

	if err != nil {
		fmt.Println("Analyzing Error..")
		customCharset := []string{"-", "\\", "|", "/", "-", ".", "o", "O", "0", "@"}
		loadingSpinner := spinner.New(customCharset, 100*time.Millisecond)
		loadingSpinner.Start()
		h.commandDebugAction(command, err, token, loadingSpinner)
	} else {
		utils.PrintMessage("No errors found in the execution")
	}
}

func (h DebugHandler) commandDebugAction(command string, err error, token string, loadingSpinner *spinner.Spinner) {
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
