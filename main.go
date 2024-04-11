package main

import (
	"context"
	"os"
	"strings"

	"github.com/alexflint/go-arg"

	"github.com/rishijash/idk_terminal/configs"
	"github.com/rishijash/idk_terminal/internal/handler"
	"github.com/rishijash/idk_terminal/internal/utils"
)

func main() {
	ctx := context.Background()

	var args struct {
		Prompt []string `arg:"positional" help:"prompt in plain english to execute terminal commands or scripts"`
		Login  bool     `arg:"--login" help:"login to idk cli"`
		Logout bool     `arg:"--logout" help:"logout from idk cli"`
		Readme string   `arg:"--readme" help:"path of your script's readme file to use with prompt"`
		Alias  string   `arg:"--alias" help:"set alias for your terminal commands or scripts"`
		Update bool     `arg:"--update" help:"update idk to the latest version"`
	}
	arg.MustParse(&args)

	appConfigs, err := configs.LoadConfig()
	if err != nil {
		println("Error running the script. Please try again!")
		return
	}

	loginHandler := handler.NewLoginHandler(appConfigs)
	promptHandler := handler.NewPromptHandler(appConfigs)

	prompt := strings.Join(args.Prompt, " ")

	if args.Update {
		utils.RunCommand("curl -o- https://idk-cli.github.io/scripts/install.sh | bash")
		return
	}

	if args.Login {
		err := loginHandler.HandleLogin(ctx)
		if err != nil {
			println("Failed to Sign In With Google. Please try again!")
			return
		}
		println("Login Successful")
		println("Try: `idk <your prompt>`")
		println("Learn more :`idk -h`")
		return
	}

	if args.Logout {
		_ = loginHandler.HandleLogout(ctx)
		println("Logout Successful")
		return
	}

	if args.Readme != "" {
		_, err := os.Stat(args.Readme)
		if err != nil {
			println("Invalid README file path")
			return
		}
	}

	err = loginHandler.HandleLoginVerification(ctx)
	if err != nil {
		println("You are not logged in. Please login first")
		println("Command: `idk --login`")
		return
	}

	promptHandler.HandlePrompt(prompt, args.Readme, args.Alias)
}
