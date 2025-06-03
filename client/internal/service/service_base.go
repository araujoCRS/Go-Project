package service

import (
	"context"
)

type ServiceBase[T any] interface {
	Get(ctx context.Context, id int) (*T, error)
	Save(ctx context.Context, model T) (*T, error)
	Delete(ctx context.Context, id int) (bool, error)
}
