package arangorepository

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/rogerwesterbo/familytree/pkg/interfaces"
)

// PersonRepository implements the PersonRepository interface using ArangoDB
type PersonRepository struct {
	*BaseRepository[interfaces.Person, *interfaces.Person]
}

// NewPersonRepository creates a new person repository
func NewPersonRepository(db arangodb.Database, collection arangodb.Collection) *PersonRepository {
	return &PersonRepository{
		BaseRepository: NewBaseRepository[interfaces.Person, *interfaces.Person](db, collection, "persons"),
	}
}

// FindByName finds persons by first name and/or last name
func (r *PersonRepository) FindByName(ctx context.Context, firstName, lastName string) ([]interfaces.Person, error) {
	query := `
		FOR p IN persons
		FILTER (@firstName == "" || p.firstName == @firstName)
		   AND (@lastName == "" || p.lastName == @lastName)
		RETURN p
	`

	bindVars := map[string]any{
		"firstName": firstName,
		"lastName":  lastName,
	}

	cursor, err := r.db.Query(ctx, query, &arangodb.QueryOptions{BindVars: bindVars})
	if err != nil {
		return nil, fmt.Errorf("failed to query persons by name: %w", err)
	}
	defer func() {
		_ = cursor.Close()
	}()

	var persons []interfaces.Person
	for cursor.HasMore() {
		var person interfaces.Person
		_, err := cursor.ReadDocument(ctx, &person)
		if err != nil {
			return nil, fmt.Errorf("failed to read person: %w", err)
		}
		persons = append(persons, person)
	}

	return persons, nil
}
