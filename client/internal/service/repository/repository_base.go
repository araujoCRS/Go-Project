package repository

import (
	"context"
)

type RepositoryBase[T any] interface {
	Get(ctx context.Context, id int) (*T, error)
	Create(ctx context.Context, entity *T) (*T, error)
	Delete(ctx context.Context, id int) (bool, error)
	Update(ctx context.Context, entity *T) (*T, error)
}
