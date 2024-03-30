package handler

import (
	"context"

	"github.com/rishijash/idk_terminal/internal/clients"
	"github.com/rishijash/idk_terminal/internal/utils"
)

type LoginHandler struct {
	config *utils.Config
}

func NewLoginHandler(config *utils.Config) LoginHandler {
	return LoginHandler{
		config: config,
	}
}

func (h LoginHandler) HandleLogin(ctx context.Context) error {
	googleToken, err := clients.SignInWithGoogle(ctx, h.config.GoogleOAuth2ClientId, h.config.GoogleOAuth2Secret)
	if err != nil {
		return err
	}

	jwtToken, err := clients.CreateIDKToken(googleToken.AccessToken, googleToken.RefreshToken, h.config.IdkBackendBaseUrl)
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
