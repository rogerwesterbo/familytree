package arangorepository

import (
	"context"
	"fmt"
	"time"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/arangodb/shared"
)

// Entity is an interface that all entities must implement to work with the base repository
type Entity interface {
	SetMetadata(key, id, rev string)
	SetTimestamps(createdAt, updatedAt time.Time)
	GetUpdatedAt() time.Time
}

// BaseRepository provides common CRUD operations for any entity type
// T must be a pointer to a type that implements Entity
type BaseRepository[T any, PT interface {
	*T
	Entity
}] struct {
	collection     arangodb.Collection
	db             arangodb.Database
	collectionName string
}

// NewBaseRepository creates a new base repository
func NewBaseRepository[T any, PT interface {
	*T
	Entity
}](db arangodb.Database, collection arangodb.Collection, collectionName string) *BaseRepository[T, PT] {
	return &BaseRepository[T, PT]{
		db:             db,
		collection:     collection,
		collectionName: collectionName,
	}
}

// Create creates a new entity
func (r *BaseRepository[T, PT]) Create(ctx context.Context, entity PT) error {
	now := time.Now()
	entity.SetTimestamps(now, now)

	meta, err := r.collection.CreateDocument(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create entity: %w", err)
	}

	entity.SetMetadata(meta.Key, string(meta.ID), meta.Rev)

	return nil
}

// GetByID retrieves an entity by ID
func (r *BaseRepository[T, PT]) GetByID(ctx context.Context, id string) (PT, error) {
	var entity T
	meta, err := r.collection.ReadDocument(ctx, id, &entity)
	if err != nil {
		if shared.IsNotFound(err) {
			return nil, fmt.Errorf("entity not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get entity: %w", err)
	}

	ptr := PT(&entity)
	ptr.SetMetadata(meta.Key, string(meta.ID), meta.Rev)

	return ptr, nil
}

// Update updates an existing entity
func (r *BaseRepository[T, PT]) Update(ctx context.Context, id string, entity PT) error {
	now := time.Now()
	entity.SetTimestamps(entity.GetUpdatedAt(), now)

	meta, err := r.collection.UpdateDocument(ctx, id, entity)
	if err != nil {
		if shared.IsNotFound(err) {
			return fmt.Errorf("entity not found: %w", err)
		}
		return fmt.Errorf("failed to update entity: %w", err)
	}

	entity.SetMetadata(meta.Key, string(meta.ID), meta.Rev)

	return nil
}

// Delete deletes an entity by ID
func (r *BaseRepository[T, PT]) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteDocument(ctx, id)
	if err != nil {
		if shared.IsNotFound(err) {
			return fmt.Errorf("entity not found: %w", err)
		}
		return fmt.Errorf("failed to delete entity: %w", err)
	}

	return nil
}

// List retrieves all entities
func (r *BaseRepository[T, PT]) List(ctx context.Context) ([]T, error) {
	query := fmt.Sprintf("FOR doc IN %s RETURN doc", r.collectionName)
	cursor, err := r.db.Query(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query entities: %w", err)
	}
	defer func() {
		_ = cursor.Close()
	}()

	var entities []T
	for cursor.HasMore() {
		var entity T
		_, err := cursor.ReadDocument(ctx, &entity)
		if err != nil {
			return nil, fmt.Errorf("failed to read entity: %w", err)
		}
		entities = append(entities, entity)
	}

	return entities, nil
}
