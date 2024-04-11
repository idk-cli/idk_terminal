package handler

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

func (h PromptHandler) HandlePrompt(prompt string, readme string, alias string) {
	handlePromptImpl(prompt, readme, "", alias, h)
}

func handlePromptImpl(prompt string, readme string, existingScript string, alias string, h PromptHandler) {
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

	aliasName := utils.RemoveWhiteSpaceFromString(alias)

	loadingSpinner.Stop()

	switch promptResponse.ActionType {
	case "COMMAND":
		if aliasName != "" {
			commandAliasAction(promptResponse.Response, aliasName)
		} else {
			commandAction(promptResponse.Response)
		}
	case "COMMANDFROMREADME":
		commandAction(promptResponse.Response)
	case "CD":
		cdAction(promptResponse.Response)
	case "SCRIPT":
		if aliasName != "" {
			scriptAliasAction(promptResponse.Response, aliasName, h)
		} else {
			scriptAction(promptResponse.Response, h)
		}
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
		// readme and alias is set to empty since scripts don't support readme
		handlePromptImpl(updateResponse, "", script, "", h)
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

func scriptAliasAction(script string, aliasName string, h PromptHandler) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Script:")
	fmt.Println("----------------")
	fmt.Println(script)
	fmt.Println("----------------")
	fmt.Printf("Do you want me to alias the script? (y/n/update): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response) // Trim whitespace and newline character
	var err error = nil

	currentTime := time.Now()
	timestampFromated := currentTime.Format("2006-01-02_15-04-05")
	scriptFileName := fmt.Sprintf("idk_script_%s.sh", timestampFromated)

	if strings.ToLower(response) == "y" {
		// save script to .idk/script folder
		idkFolder := utils.GetAbsoluteHomeDirectoryPath([]string{".idk", "scripts"})
		filePath := fmt.Sprintf("%s/%s", idkFolder, scriptFileName)
		dirPath := filepath.Dir(filePath)
		os.MkdirAll(dirPath, 0777)
		saveScript(script, filePath)
		command := fmt.Sprintf(". %s", filePath)
		err = aliasCommand(command, aliasName)
	} else if strings.ToLower(response) == "update" {
		fmt.Println("What do you want to change?")
		reader = bufio.NewReader(os.Stdin)
		updateResponse, _ := reader.ReadString('\n')
		// readme is set to empty since scripts don't support readme
		handlePromptImpl(updateResponse, "", script, aliasName, h)
	} else {
		fmt.Println("Script execution canceled")
	}

	if err != nil {
		fmt.Println(err.Error())
		// fmt.Println("Something went wrong. Please try again!")
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
// CD Logic
// ----------------------------------------------------------------------------------------

func cdAction(folderName string) {
	// try to find action from terminal history
	historyList, err := utils.GetTerminalHistory()
	if err != nil {
		fmt.Println("Something went wrong. Please try again!")
		return
	}

	cdHistoryList := utils.FilterByPrefix(historyList, "cd ")

	absoluteCDCommandPaths := utils.ResolveCDPathCommands(cdHistoryList)

	matchingCDPathCommand := utils.FindMostRelevantStringFromArr(absoluteCDCommandPaths, folderName)

	if matchingCDPathCommand == "" {
		fmt.Println("Could not find the directory. Please try again!")
		return
	}

	commandAction(matchingCDPathCommand)
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

func commandAliasAction(command string, aliasName string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Do you want me to alias `%s` as `%s`? (y/n): ", command, aliasName)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response) // Trim whitespace and newline character
	var err error = nil

	if strings.ToLower(response) == "y" {
		err = aliasCommand(command, aliasName)
	} else {
		fmt.Println("Command execution canceled")
	}

	if err != nil {
		fmt.Println("Something went wrong. Please try again!")
	}
}

func aliasCommand(commandStr string, aliasName string) error {
	var configFile string
	shellPath := os.Getenv("SHELL")
	shell := filepath.Base(shellPath)

	// Determine the shell's configuration file based on the shell type
	switch shell {
	case "bash":
		configFile = filepath.Join(os.Getenv("HOME"), ".bashrc")
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			configFile = filepath.Join(os.Getenv("HOME"), ".bash_profile")
		}
	case "zsh":
		configFile = filepath.Join(os.Getenv("HOME"), ".zshrc")
		// Ensure the configuration file exists; create it if it does not.
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			file, err := os.Create(configFile)
			if err != nil {
				return err
			}
			file.Close()
		}
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}

	// Construct the alias command
	aliasCmd := fmt.Sprintf("alias %s='%s'\n", aliasName, commandStr)

	// also run command so it is also applied to existing
	utils.RunCommand(aliasCmd)

	// Check if the alias already exists to avoid duplicates
	if aliasExists(configFile, aliasName) {
		// do nothing since it already exists
		return fmt.Errorf("Alias already exists")
	}

	// Open the configuration file in append mode
	file, err := os.OpenFile(configFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// create a backup before updating the file
	err = utils.BackupFile(configFile)
	if err != nil {
		return err
	}

	// Append the alias command to the file
	if _, err := file.WriteString(aliasCmd); err != nil {
		return err
	}

	fmt.Printf("Added alias '%s' to %s\n", aliasName, configFile)
	return nil
}

func aliasExists(configFile, aliasName string) bool {
	file, err := os.Open(configFile)
	if err != nil {
		return false // Cannot open the file, assume the alias doesn't exist
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	aliasPrefix := fmt.Sprintf("alias %s=", aliasName)
	for scanner.Scan() {
		// Check if the line starts with the alias definition
		if strings.HasPrefix(scanner.Text(), aliasPrefix) {
			return true // Found the alias, it already exists
		}
	}

	return false // Alias does not exist
}
