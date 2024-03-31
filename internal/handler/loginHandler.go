package handler

import (
	"context"

	"github.com/rishijash/idk_terminal/configs"
	"github.com/rishijash/idk_terminal/internal/clients"
	"github.com/rishijash/idk_terminal/internal/utils"
)

type LoginHandler struct {
	config *configs.Config
}

func NewLoginHandler(config *configs.Config) LoginHandler {
	return LoginHandler{
		config: config,
	}
}

func (h LoginHandler) HandleLogin(ctx context.Context) error {
	state := utils.GenerateRandomString(10)
	authUrl, err := clients.CreateGoogleAuthCodeURL(state, h.config.IdkBackendBaseUrl)
	if err != nil {
		return err
	}

	googleAuthCode, err := clients.StartAuthFlow(ctx, authUrl, state)
	if err != nil {
		return err
	}

	jwtToken, err := clients.CreateIDKToken(googleAuthCode, h.config.IdkBackendBaseUrl)
	if err != nil {
		return err
	}

	err = utils.SaveToken(jwtToken)
	if err != nil {
		return err
	}

	return nil
}

func (h LoginHandler) HandleLoginVerification(ctx context.Context) error {
	_, err := utils.LoadToken()
	if err != nil {
		return err
	}

	return nil
}

func (h LoginHandler) HandleLogout(ctx context.Context) error {
	return utils.ClearToken()
}
