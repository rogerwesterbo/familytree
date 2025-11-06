package interfaces

import "time"

// Person represents a person in the family tree
type Person struct {
	Key       string    `json:"_key,omitempty"`
	ID        string    `json:"_id,omitempty"`
	Rev       string    `json:"_rev,omitempty"`
	FirstName string    `json:"firstName" binding:"required"`
	LastName  string    `json:"lastName" binding:"required"`
	BirthDate time.Time `json:"birthDate,omitempty"`
	DeathDate time.Time `json:"deathDate,omitempty"`
	Gender    string    `json:"gender,omitempty"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// PersonCreateRequest represents the request body for creating a person
type PersonCreateRequest struct {
	FirstName string    `json:"firstName" binding:"required" example:"John"`
	LastName  string    `json:"lastName" binding:"required" example:"Doe"`
	BirthDate time.Time `json:"birthDate,omitempty" example:"1980-01-15T00:00:00Z"`
	DeathDate time.Time `json:"deathDate,omitempty" example:"2050-05-20T00:00:00Z"`
	Gender    string    `json:"gender,omitempty" example:"male"`
	Email     string    `json:"email,omitempty" example:"john.doe@example.com"`
	Phone     string    `json:"phone,omitempty" example:"+1234567890"`
}

// PersonUpdateRequest represents the request body for updating a person
type PersonUpdateRequest struct {
	FirstName string    `json:"firstName,omitempty" example:"John"`
	LastName  string    `json:"lastName,omitempty" example:"Doe"`
	BirthDate time.Time `json:"birthDate,omitempty" example:"1980-01-15T00:00:00Z"`
	DeathDate time.Time `json:"deathDate,omitempty" example:"2050-05-20T00:00:00Z"`
	Gender    string    `json:"gender,omitempty" example:"male"`
	Email     string    `json:"email,omitempty" example:"john.doe@example.com"`
	Phone     string    `json:"phone,omitempty" example:"+1234567890"`
}

// PersonResponse represents the response body for person operations
type PersonResponse struct {
	Person  *Person `json:"person,omitempty"`
	Message string  `json:"message,omitempty"`
}

// PersonsListResponse represents the response body for listing persons
type PersonsListResponse struct {
	Persons []Person `json:"persons"`
	Count   int      `json:"count"`
}

// SetMetadata sets the ArangoDB metadata fields
func (p *Person) SetMetadata(key, id, rev string) {
	p.Key = key
	p.ID = id
	p.Rev = rev
}

// SetTimestamps sets the created and updated timestamps
func (p *Person) SetTimestamps(createdAt, updatedAt time.Time) {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = createdAt
	}
	p.UpdatedAt = updatedAt
}

// GetUpdatedAt returns the updated timestamp
func (p Person) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}
