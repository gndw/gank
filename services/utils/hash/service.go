package hash

import "context"

type Service interface {
	IsPasswordMatch(ctx context.Context, password string, hash string) (err error)
	GetPasswordHash(ctx context.Context, password string) (hash string, err error)
}
