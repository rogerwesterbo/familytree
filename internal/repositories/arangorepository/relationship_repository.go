package arangorepository

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/rogerwesterbo/familytree/pkg/interfaces"
)

// RelationshipRepository implements the RelationshipRepository interface using ArangoDB
type RelationshipRepository struct {
	*BaseRepository[interfaces.Relationship, *interfaces.Relationship]
}

// NewRelationshipRepository creates a new relationship repository
func NewRelationshipRepository(db arangodb.Database, collection arangodb.Collection) *RelationshipRepository {
	return &RelationshipRepository{
		BaseRepository: NewBaseRepository[interfaces.Relationship, *interfaces.Relationship](db, collection, "relationships"),
	}
}

// FindByPerson finds all relationships for a person (as either from or to)
func (r *RelationshipRepository) FindByPerson(ctx context.Context, personID string) ([]interfaces.Relationship, error) {
	query := `
		FOR rel IN relationships
		FILTER rel._from == @personID || rel._to == @personID
		RETURN rel
	`

	bindVars := map[string]any{
		"personID": personID,
	}

	cursor, err := r.db.Query(ctx, query, &arangodb.QueryOptions{BindVars: bindVars})
	if err != nil {
		return nil, fmt.Errorf("failed to query relationships by person: %w", err)
	}
	defer func() {
		_ = cursor.Close()
	}()

	var relationships []interfaces.Relationship
	for cursor.HasMore() {
		var relationship interfaces.Relationship
		_, err := cursor.ReadDocument(ctx, &relationship)
		if err != nil {
			return nil, fmt.Errorf("failed to read relationship: %w", err)
		}
		relationships = append(relationships, relationship)
	}

	return relationships, nil
}

// FindByType finds relationships by type
func (r *RelationshipRepository) FindByType(ctx context.Context, relationType string) ([]interfaces.Relationship, error) {
	query := `
		FOR rel IN relationships
		FILTER rel.relationType == @relationType
		RETURN rel
	`

	bindVars := map[string]any{
		"relationType": relationType,
	}

	cursor, err := r.db.Query(ctx, query, &arangodb.QueryOptions{BindVars: bindVars})
	if err != nil {
		return nil, fmt.Errorf("failed to query relationships by type: %w", err)
	}
	defer func() {
		_ = cursor.Close()
	}()

	var relationships []interfaces.Relationship
	for cursor.HasMore() {
		var relationship interfaces.Relationship
		_, err := cursor.ReadDocument(ctx, &relationship)
		if err != nil {
			return nil, fmt.Errorf("failed to read relationship: %w", err)
		}
		relationships = append(relationships, relationship)
	}

	return relationships, nil
}
