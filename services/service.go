package services

import (
	"context"

	vault "github.com/hashicorp/vault/api"
)

type UserService interface {
	GetToken() (*vault.Secret, error)
	GetSecret(ctx context.Context, dest string) (*vault.KVSecret, error)
	PutSecret(ctx context.Context, dest string, secret map[string]interface{}) (*vault.KVSecret, error)
	PatchSecret(ctx context.Context, dest string, newSecret map[string]interface{}) (*vault.KVSecret, error)
}
