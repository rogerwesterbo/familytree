package v1relationshipservice

import (
	"context"
	"fmt"
	"strings"

	"github.com/rogerwesterbo/familytree/pkg/interfaces"
)

// RelationshipService handles business logic for relationship operations
type RelationshipService struct {
	repo interfaces.RelationshipRepository
}

// NewRelationshipService creates a new relationship service
func NewRelationshipService(repo interfaces.RelationshipRepository) *RelationshipService {
	return &RelationshipService{
		repo: repo,
	}
}

// CreateRelationship creates a new relationship with validation
func (s *RelationshipService) CreateRelationship(ctx context.Context, req *interfaces.RelationshipCreateRequest) (*interfaces.Relationship, error) {
	// Validate required fields
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Create relationship entity
	relationship := &interfaces.Relationship{
		From:         strings.TrimSpace(req.From),
		To:           strings.TrimSpace(req.To),
		RelationType: strings.TrimSpace(req.RelationType),
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Notes:        strings.TrimSpace(req.Notes),
	}

	// Create in repository
	if err := s.repo.Create(ctx, relationship); err != nil {
		return nil, fmt.Errorf("failed to create relationship: %w", err)
	}

	return relationship, nil
}

// GetRelationship retrieves a relationship by ID
func (s *RelationshipService) GetRelationship(ctx context.Context, id string) (*interfaces.Relationship, error) {
	if id == "" {
		return nil, fmt.Errorf("relationship ID is required")
	}

	relationship, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return relationship, nil
}

// UpdateRelationship updates an existing relationship
func (s *RelationshipService) UpdateRelationship(ctx context.Context, id string, req *interfaces.RelationshipUpdateRequest) (*interfaces.Relationship, error) {
	if id == "" {
		return nil, fmt.Errorf("relationship ID is required")
	}

	// Get existing relationship
	relationship, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.From != "" {
		relationship.From = strings.TrimSpace(req.From)
	}
	if req.To != "" {
		relationship.To = strings.TrimSpace(req.To)
	}
	if req.RelationType != "" {
		relationship.RelationType = strings.TrimSpace(req.RelationType)
	}
	if !req.StartDate.IsZero() {
		relationship.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() {
		relationship.EndDate = req.EndDate
	}
	if req.Notes != "" {
		relationship.Notes = strings.TrimSpace(req.Notes)
	}

	// Update in repository
	if err := s.repo.Update(ctx, id, relationship); err != nil {
		return nil, fmt.Errorf("failed to update relationship: %w", err)
	}

	return relationship, nil
}

// DeleteRelationship deletes a relationship
func (s *RelationshipService) DeleteRelationship(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("relationship ID is required")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// ListRelationships retrieves all relationships
func (s *RelationshipService) ListRelationships(ctx context.Context) ([]interfaces.Relationship, error) {
	relationships, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return relationships, nil
}

// GetRelationshipsForPerson gets all relationships for a person
func (s *RelationshipService) GetRelationshipsForPerson(ctx context.Context, personID string) ([]interfaces.Relationship, error) {
	if personID == "" {
		return nil, fmt.Errorf("person ID is required")
	}

	relationships, err := s.repo.FindByPerson(ctx, personID)
	if err != nil {
		return nil, err
	}

	return relationships, nil
}

// GetRelationshipsByType gets relationships by type
func (s *RelationshipService) GetRelationshipsByType(ctx context.Context, relationType string) ([]interfaces.Relationship, error) {
	if relationType == "" {
		return nil, fmt.Errorf("relationship type is required")
	}

	relationships, err := s.repo.FindByType(ctx, relationType)
	if err != nil {
		return nil, err
	}

	return relationships, nil
}

// validateCreateRequest validates a relationship create request
func (s *RelationshipService) validateCreateRequest(req *interfaces.RelationshipCreateRequest) error {
	if strings.TrimSpace(req.From) == "" {
		return fmt.Errorf("from is required")
	}
	if strings.TrimSpace(req.To) == "" {
		return fmt.Errorf("to is required")
	}
	if strings.TrimSpace(req.RelationType) == "" {
		return fmt.Errorf("relationType is required")
	}

	// Validate that from and to are different
	if req.From == req.To {
		return fmt.Errorf("from and to must be different persons")
	}

	// Validate relationship type
	validTypes := map[string]bool{
		interfaces.RelationTypeParent:  true,
		interfaces.RelationTypeChild:   true,
		interfaces.RelationTypeSpouse:  true,
		interfaces.RelationTypeSibling: true,
	}

	if !validTypes[req.RelationType] {
		return fmt.Errorf("invalid relationship type: %s. Valid types are: parent, child, spouse, sibling", req.RelationType)
	}

	return nil
}
