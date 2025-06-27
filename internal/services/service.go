package services

import (
	"context"
	appcore_config "hashicorp-vault/cmd/hashicorp-vault/config"

	vault "github.com/hashicorp/vault/api"
)

type UserService interface {
	GetSecret(ctx context.Context, dest string) (*vault.KVSecret, error)
	GetConfig(ctx context.Context) (*appcore_config.Configuration, error)
}
