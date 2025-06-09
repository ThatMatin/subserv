package repo

import "context"

type UserRepository interface {
	Exists(ctx context.Context, ID uint) (bool, error)
}
