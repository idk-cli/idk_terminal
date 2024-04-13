package handler

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/briandowns/spinner"

	"github.com/rishijash/idk_terminal/configs"
	"github.com/rishijash/idk_terminal/internal/clients"
	"github.com/rishijash/idk_terminal/internal/utils"
)

type PromptHandler struct {
	config *configs.Config
}

func NewPromptHandler(config *configs.Config) PromptHandler {
	return PromptHandler{
		config: config,
	}
}

func (h PromptHandler) HandlePrompt(prompt string, readme string) {
	handlePromptImpl(prompt, readme, "", h)
}

func handlePromptImpl(prompt string, readme string, existingScript string, h PromptHandler) {
	if prompt == "" {
		println("Your prompt can not be empty")
		println("Learn more :`idk -h`")
		return
	}

	token, err := utils.LoadToken()
	if err != nil {
		println("You are not logged in. Please login first")
		println("Command: `idk --login`")
		return
	}

	readmeData := ""
	if readme != "" {
		readmeDataBytes, err := os.ReadFile(readme)
		if err != nil {
			fmt.Println("Error fetching README file. Please try again!")
			return
		}
		readmeData = string(readmeDataBytes)
	}

	customCharset := []string{"-", "\\", "|", "/", "-", ".", "o", "O", "0", "@"}
	loadingSpinner := spinner.New(customCharset, 100*time.Millisecond)
	loadingSpinner.Start()

	pwd, err := os.Getwd()
	if err != nil {
		pwd = ""
	}

	promptResponse, err, responseStatus := clients.ProcessPrompt(prompt, runtime.GOOS, readmeData, existingScript, pwd, token, h.config.IdkBackendBaseUrl)
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

	loadingSpinner.Stop()

	switch promptResponse.ActionType {
	case "COMMAND":
		commandAction(promptResponse.Response)
	case "COMMANDFROMREADME":
		commandAction(promptResponse.Response)
	case "SCRIPT":
		scriptAction(promptResponse.Response, h)
	default:
		println(promptResponse.Response)
	}
}

// ----------------------------------------------------------------------------------------
// Script Logic
// ----------------------------------------------------------------------------------------
func scriptAction(script string, h PromptHandler) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Script:")
	fmt.Println("----------------")
	fmt.Println(script)
	fmt.Println("----------------")
	fmt.Printf("Do you want me to execute the script? (y/n/update/save): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response) // Trim whitespace and newline character
	var err error = nil

	currentTime := time.Now()
	timestampFromated := currentTime.Format("2006-01-02_15-04-05")
	scriptFileName := fmt.Sprintf("idk_script_%s.sh", timestampFromated)

	if strings.ToLower(response) == "y" {
		err = runScript(script, scriptFileName)
		fmt.Println("Script execution completed")
	} else if strings.ToLower(response) == "update" {
		fmt.Println("What do you want to change?")
		reader = bufio.NewReader(os.Stdin)
		updateResponse, _ := reader.ReadString('\n')
		// readme is set to empty since scripts don't support readme
		handlePromptImpl(updateResponse, "", script, h)
	} else if strings.ToLower(response) == "save" {
		err = saveScript(script, scriptFileName)
		fmt.Printf("Script saved as %s", scriptFileName)
	} else {
		fmt.Println("Script execution canceled")
	}

	if err != nil {
		fmt.Println("Something went wrong. Please try again!")
	}
}

func saveScript(script string, filePath string) error {
	scriptBytes := []byte(script)
	err := os.WriteFile(filePath, scriptBytes, 0600)
	return err
}

func runScript(script string, fileName string) error {
	// save file
	scriptBytes := []byte(script)
	err := os.WriteFile(fileName, scriptBytes, 0600)
	if err != nil {
		return err
	}

	utils.RunCommand(fmt.Sprintf(". %s", fileName))

	err = os.Remove(fileName)
	if err != nil {
		return err
	}

	return nil
}

// ----------------------------------------------------------------------------------------
// Command Logic
// ----------------------------------------------------------------------------------------

func commandAction(command string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Do you want me to execute `%s`? (y/n/copy): ", command)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response) // Trim whitespace and newline character
	var err error = nil

	if strings.ToLower(response) == "y" {
		err = utils.RunCommand(command)
	} else if strings.ToLower(response) == "copy" {
		err := clipboard.WriteAll(command)
		if err != nil {
			fmt.Println("Failed to copy command to clipboard")
			return
		}
		fmt.Println("Command copied to clipboard")
	} else {
		fmt.Println("Command execution canceled")
	}

	if err != nil {
		fmt.Println("Something went wrong. Please try again!")
	}
}
