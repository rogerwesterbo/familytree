package interfaces

import "time"

// Relationship represents a relationship between two persons in the family tree
// This is an edge document in ArangoDB
type Relationship struct {
	Key          string    `json:"_key,omitempty"`
	ID           string    `json:"_id,omitempty"`
	Rev          string    `json:"_rev,omitempty"`
	From         string    `json:"_from" binding:"required"`
	To           string    `json:"_to" binding:"required"`
	RelationType string    `json:"relationType" binding:"required"`
	StartDate    time.Time `json:"startDate,omitempty"`
	EndDate      time.Time `json:"endDate,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// RelationshipCreateRequest represents the request body for creating a relationship
type RelationshipCreateRequest struct {
	From         string    `json:"from" binding:"required" example:"persons/123"`
	To           string    `json:"to" binding:"required" example:"persons/456"`
	RelationType string    `json:"relationType" binding:"required" example:"parent"`
	StartDate    time.Time `json:"startDate,omitempty" example:"2000-01-01T00:00:00Z"`
	EndDate      time.Time `json:"endDate,omitempty" example:"2020-12-31T00:00:00Z"`
	Notes        string    `json:"notes,omitempty" example:"Biological parent"`
}

// RelationshipUpdateRequest represents the request body for updating a relationship
type RelationshipUpdateRequest struct {
	From         string    `json:"from,omitempty" example:"persons/123"`
	To           string    `json:"to,omitempty" example:"persons/456"`
	RelationType string    `json:"relationType,omitempty" example:"parent"`
	StartDate    time.Time `json:"startDate,omitempty" example:"2000-01-01T00:00:00Z"`
	EndDate      time.Time `json:"endDate,omitempty" example:"2020-12-31T00:00:00Z"`
	Notes        string    `json:"notes,omitempty" example:"Biological parent"`
}

// RelationshipResponse represents the response body for relationship operations
type RelationshipResponse struct {
	Relationship *Relationship `json:"relationship,omitempty"`
	Message      string        `json:"message,omitempty"`
}

// RelationshipsListResponse represents the response body for listing relationships
type RelationshipsListResponse struct {
	Relationships []Relationship `json:"relationships"`
	Count         int            `json:"count"`
}

// Common relationship types
const (
	RelationTypeParent  = "parent"
	RelationTypeChild   = "child"
	RelationTypeSpouse  = "spouse"
	RelationTypeSibling = "sibling"
)

// SetMetadata sets the ArangoDB metadata fields
func (r *Relationship) SetMetadata(key, id, rev string) {
	r.Key = key
	r.ID = id
	r.Rev = rev
}

// SetTimestamps sets the created and updated timestamps
func (r *Relationship) SetTimestamps(createdAt, updatedAt time.Time) {
	if r.CreatedAt.IsZero() {
		r.CreatedAt = createdAt
	}
	r.UpdatedAt = updatedAt
}

// GetUpdatedAt returns the updated timestamp
func (r Relationship) GetUpdatedAt() time.Time {
	return r.UpdatedAt
}
