package user

import (
	"context"
	"errors"
	"fmt"
	appcore_config "hashicorp-vault/config"
	"hashicorp-vault/services"

	vault "github.com/hashicorp/vault/api"
)

var (
	appRoleAuthPath = "auth/approle/login"
	dotenvPath      = "dotenv/"
	kvMountPoint    = "secret"
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

func (s *userService) GetToken() (*vault.Secret, error) {
	options := map[string]interface{}{
		"role_id":   s.appConf.VaultRoleID,
		"secret_id": s.appConf.VaultSecretID,
	}

	resp, err := s.clientVault.Logical().Write(appRoleAuthPath, options)
	if err != nil {
		return nil, fmt.Errorf("unable to login to AppRole: %v", err)
	}

	if resp.Auth == nil {
		return nil, errors.New("resp auth is nil")
	}

	s.appConf.VaultClientToken = resp.Auth.ClientToken
	s.clientVault.SetToken(resp.Auth.ClientToken)
	return resp, nil
}

func (s *userService) GetSecret(ctx context.Context, dest string) (*vault.KVSecret, error) {
	if s.appConf.VaultClientToken == "" {
		return nil, errors.New("no vault client token")
	}

	secretPath := fmt.Sprintf("%s%s", dotenvPath, dest)
	resp, err := s.clientVault.KVv2(kvMountPoint).Get(ctx, secretPath)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *userService) PutSecret(ctx context.Context, dest string, secrets map[string]interface{}) (*vault.KVSecret, error) {
	if s.appConf.VaultClientToken == "" {
		return nil, errors.New("no vault client token")
	}
	secretPath := fmt.Sprintf("%s%s", dotenvPath, dest)
	resp, err := s.clientVault.KVv2(kvMountPoint).Put(ctx, secretPath, secrets)

	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, errors.New("resp is nil")
	}

	return resp, nil
}

func (s *userService) PatchSecret(ctx context.Context, dest string, newSecret map[string]interface{}) (*vault.KVSecret, error) {
	if s.appConf.VaultClientToken == "" {
		return nil, errors.New("no vault client token")
	}

	secretPath := fmt.Sprintf("%s%s", dotenvPath, dest)
	resp, err := s.clientVault.KVv2(kvMountPoint).Patch(ctx, secretPath, newSecret)

	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, errors.New("resp is nil")
	}

	return resp, nil
}
