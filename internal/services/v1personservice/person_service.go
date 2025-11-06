package v1personservice

import (
	"context"
	"fmt"
	"strings"

	"github.com/rogerwesterbo/familytree/pkg/interfaces"
)

// PersonService handles business logic for person operations
type PersonService struct {
	repo interfaces.PersonRepository
}

// NewPersonService creates a new person service
func NewPersonService(repo interfaces.PersonRepository) *PersonService {
	return &PersonService{
		repo: repo,
	}
}

// CreatePerson creates a new person with validation
func (s *PersonService) CreatePerson(ctx context.Context, req *interfaces.PersonCreateRequest) (*interfaces.Person, error) {
	// Validate required fields
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Create person entity
	person := &interfaces.Person{
		FirstName: strings.TrimSpace(req.FirstName),
		LastName:  strings.TrimSpace(req.LastName),
		BirthDate: req.BirthDate,
		DeathDate: req.DeathDate,
		Gender:    strings.TrimSpace(req.Gender),
		Email:     strings.TrimSpace(req.Email),
		Phone:     strings.TrimSpace(req.Phone),
	}

	// Create in repository
	if err := s.repo.Create(ctx, person); err != nil {
		return nil, fmt.Errorf("failed to create person: %w", err)
	}

	return person, nil
}

// GetPerson retrieves a person by ID
func (s *PersonService) GetPerson(ctx context.Context, id string) (*interfaces.Person, error) {
	if id == "" {
		return nil, fmt.Errorf("person ID is required")
	}

	person, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return person, nil
}

// UpdatePerson updates an existing person
func (s *PersonService) UpdatePerson(ctx context.Context, id string, req *interfaces.PersonUpdateRequest) (*interfaces.Person, error) {
	if id == "" {
		return nil, fmt.Errorf("person ID is required")
	}

	// Get existing person
	person, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != "" {
		person.FirstName = strings.TrimSpace(req.FirstName)
	}
	if req.LastName != "" {
		person.LastName = strings.TrimSpace(req.LastName)
	}
	if !req.BirthDate.IsZero() {
		person.BirthDate = req.BirthDate
	}
	if !req.DeathDate.IsZero() {
		person.DeathDate = req.DeathDate
	}
	if req.Gender != "" {
		person.Gender = strings.TrimSpace(req.Gender)
	}
	if req.Email != "" {
		person.Email = strings.TrimSpace(req.Email)
	}
	if req.Phone != "" {
		person.Phone = strings.TrimSpace(req.Phone)
	}

	// Update in repository
	if err := s.repo.Update(ctx, id, person); err != nil {
		return nil, fmt.Errorf("failed to update person: %w", err)
	}

	return person, nil
}

// DeletePerson deletes a person
func (s *PersonService) DeletePerson(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("person ID is required")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// ListPersons retrieves all persons
func (s *PersonService) ListPersons(ctx context.Context) ([]interfaces.Person, error) {
	persons, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return persons, nil
}

// SearchPersonsByName searches persons by name
func (s *PersonService) SearchPersonsByName(ctx context.Context, firstName, lastName string) ([]interfaces.Person, error) {
	persons, err := s.repo.FindByName(ctx, strings.TrimSpace(firstName), strings.TrimSpace(lastName))
	if err != nil {
		return nil, err
	}

	return persons, nil
}

// validateCreateRequest validates a person create request
func (s *PersonService) validateCreateRequest(req *interfaces.PersonCreateRequest) error {
	if strings.TrimSpace(req.FirstName) == "" {
		return fmt.Errorf("firstName is required")
	}
	if strings.TrimSpace(req.LastName) == "" {
		return fmt.Errorf("lastName is required")
	}

	// Validate email format if provided
	if req.Email != "" && !isValidEmail(req.Email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	// Basic validation - contains @ and has characters before and after
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}
