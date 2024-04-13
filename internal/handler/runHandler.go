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

type RunHandler struct {
	config *configs.Config
}

func NewRunHandler(config *configs.Config) RunHandler {
	return RunHandler{
		config: config,
	}
}

func (h RunHandler) HandleSetupProject(ctx context.Context) {
	token, err := utils.LoadToken()
	if err != nil {
		println("You are not logged in. Please login first")
		println("Command: `idk --login`")
		return
	}

	files, err := utils.ListFilesAndDirs()

	if err != nil {
		println("Something went wrong. Please try again!")
		return
	}

	readmeData, err := utils.FindReadmeData()
	if err != nil {
		println("Something went wrong. Please try again!")
		return
	}

	makefileData, err := utils.FindMakefileData()
	if err != nil {
		println("Something went wrong. Please try again!")
		return
	}

	projectFolderName, err := utils.GetCurrentDirName()
	if err != nil {
		println("Something went wrong. Please try again!")
		return
	}

	fmt.Println("Analyzing Project..")
	customCharset := []string{"-", "\\", "|", "/", "-", ".", "o", "O", "0", "@"}
	loadingSpinner := spinner.New(customCharset, 100*time.Millisecond)
	loadingSpinner.Start()

	response, err, responseStatus := clients.ProcessGetProjectInit(
		projectFolderName, files, readmeData, makefileData, runtime.GOOS, token, h.config.IdkBackendBaseUrl)

	loadingSpinner.Stop()

	if !utils.IsBrewInstalled() {
		println("Brew is not installed. Installing brew first")
		err := utils.InstallBrew()
		if err != nil {
			println("failed to install brew. Please manually install brew before continuing")
			return
		}
	}

	if h.isErrorResonse(responseStatus, err) {
		return
	}

	if len(response.Commands) == 0 {
		println("Something went wrong. Please try again!")
		return
	}

	h.executeCommandsAction(response.ProjectType, response.Commands)

}

func (h RunHandler) executeCommandsAction(projectType string, commands []clients.RunGetProjectInitCommand) {
	utils.PrintMessages([]string{
		fmt.Sprintf("`%s` found", projectType),
		"Commands will be executed in sequence to get your project setup:",
	})

	for i, command := range commands {
		if i == len(commands)-1 {
			// last command is to run the project skip it
			continue
		}

		reader := bufio.NewReader(os.Stdin)
		println(fmt.Sprintf("[Step %d / %d]", i+1, len(commands)-1))
		println(fmt.Sprintf("Command: %s", command.Command))
		println(fmt.Sprintf("Description: %s", command.Description))
		println("")
		println("Continue? (y/skip/stop)")
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response) // Trim whitespace and newline character
		if response == "y" {
			err := utils.RunCommand(command.Command)
			if err != nil {
				println("Error setting up project. Please try again!")
				return
			}
		} else if response == "skip" {
			continue
		} else {
			println("Project Setup Cancelled")
			return
		}
	}
	utils.PrintMessages([]string{
		"Project Setup Completed",
		"",
		"Run your Project with following command:",
		commands[len(commands)-1].Command,
	})
}

func (h RunHandler) isErrorResonse(responseStatus int, err error) bool {
	if responseStatus == http.StatusUnauthorized {
		utils.ClearToken()
		fmt.Println("Token expired. Please login again")
		println("Command: `idk --login`")
		return true
	}

	if responseStatus == http.StatusTooManyRequests {
		fmt.Println("Daily Quota limit reached. Plesae try again tomorrow or upgrade on https://idk-cli.github.io/")
		return true
	}

	if err != nil {
		fmt.Println("Something went wrong. Please try again!")
		return true
	}

	return false

}
