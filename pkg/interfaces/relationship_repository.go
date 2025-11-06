package interfaces

import "context"

// RelationshipRepository defines the interface for relationship data access operations
// It embeds the generic Repository interface and adds relationship-specific methods
type RelationshipRepository interface {
	Repository[Relationship]

	// FindByPerson finds all relationships for a person (as either from or to)
	FindByPerson(ctx context.Context, personID string) ([]Relationship, error)

	// FindByType finds relationships by type
	FindByType(ctx context.Context, relationType string) ([]Relationship, error)
}
