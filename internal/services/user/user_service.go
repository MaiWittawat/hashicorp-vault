package user

import (
	"context"
	"errors"
	"fmt"

	appcore_config "hashicorp-vault/cmd/hashicorp-vault/config"
	"hashicorp-vault/internal/services"

	vault "github.com/hashicorp/vault/api"
)

var (
	dotenvPath   = "dotenv/"
	kvMountPoint = "secret"
)

type userService struct {
	appConf     *appcore_config.Configuration
	clientVault *vault.Client
}

func NewUserService(appConf *appcore_config.Configuration, clientVault *vault.Client) services.UserService {
	return &userService{
		appConf:     appConf,
		clientVault: clientVault,
	}
}

func (s *userService) GetSecret(ctx context.Context, dest string) (*vault.KVSecret, error) {
	if !s.appConf.VaultTokenSet {
		return nil, errors.New("no vault client token")
	}

	secretPath := fmt.Sprintf("%s%s", dotenvPath, dest)
	resp, err := s.clientVault.KVv2(kvMountPoint).Get(ctx, secretPath)
	if err != nil {
		return nil, err
	}

	return resp, nil
}


func (s *userService) GetConfig(ctx context.Context) (*appcore_config.Configuration, error) {
	if !s.appConf.VaultTokenSet {
		return nil, errors.New("no vault client token")
	}

	return appcore_config.Config, nil
}
