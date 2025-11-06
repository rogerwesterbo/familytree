package interfaces

import "context"

// PersonRepository defines the interface for person data access operations
// It embeds the generic Repository interface and adds person-specific methods
type PersonRepository interface {
	Repository[Person]

	// FindByName finds persons by first name and/or last name
	FindByName(ctx context.Context, firstName, lastName string) ([]Person, error)
}
