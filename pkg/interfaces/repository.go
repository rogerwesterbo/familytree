package interfaces

import "context"

// Repository defines a generic repository interface for CRUD operations
type Repository[T any] interface {
	// Create creates a new entity
	Create(ctx context.Context, entity *T) error

	// GetByID retrieves an entity by ID
	GetByID(ctx context.Context, id string) (*T, error)

	// Update updates an existing entity
	Update(ctx context.Context, id string, entity *T) error

	// Delete deletes an entity by ID
	Delete(ctx context.Context, id string) error

	// List retrieves all entities
	List(ctx context.Context) ([]T, error)
}
