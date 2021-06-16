package service

import (
	"context"
)

// Store is an interface for handling the storing of favorites.
type Store interface {
	GetAll(ctx context.Context)
	Create(ctx context.Context)
	Delete(ctx context.Context)
}