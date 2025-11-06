package v1personshandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rogerwesterbo/familytree/internal/httpserver/helpers"
	"github.com/rogerwesterbo/familytree/internal/services/v1personservice"
	"github.com/rogerwesterbo/familytree/pkg/interfaces"
)

// Handler handles HTTP requests for person operations
type Handler struct {
	service *v1personservice.PersonService
}

// NewHandler creates a new person handler
func NewHandler(service *v1personservice.PersonService) *Handler {
	return &Handler{
		service: service,
	}
}

// HandlePersons routes person requests based on HTTP method
// @Summary Person operations
// @Description Handle person CRUD operations
// @Tags persons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security OAuth2Password
// @Router /v1/persons [get]
// @Router /v1/persons [post]
// @Router /v1/persons/{id} [get]
// @Router /v1/persons/{id} [put]
// @Router /v1/persons/{id} [delete]
func (h *Handler) HandlePersons(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if strings.HasPrefix(r.URL.Path, "/v1/persons/") {
			personID := strings.TrimPrefix(r.URL.Path, "/v1/persons/")
			if personID != "" {
				h.GetPerson(w, r, personID)
				return
			}
		}
		h.ListPersons(w, r)
	case http.MethodPost:
		h.CreatePerson(w, r)
	case http.MethodPut:
		personID := strings.TrimPrefix(r.URL.Path, "/v1/persons/")
		if personID == "" {
			helpers.SendError(w, http.StatusBadRequest, "person ID is required")
			return
		}
		h.UpdatePerson(w, r, personID)
	case http.MethodDelete:
		personID := strings.TrimPrefix(r.URL.Path, "/v1/persons/")
		if personID == "" {
			helpers.SendError(w, http.StatusBadRequest, "person ID is required")
			return
		}
		h.DeletePerson(w, r, personID)
	default:
		helpers.SendError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// ListPersons returns all persons
// @Summary List all persons
// @Description Get a list of all persons in the family tree
// @Tags persons
// @Accept json
// @Produce json
// @Success 200 {object} interfaces.PersonsListResponse
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Security OAuth2Password
// @Router /v1/persons [get]
func (h *Handler) ListPersons(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	persons, err := h.service.ListPersons(ctx)
	if err != nil {
		helpers.SendError(w, http.StatusInternalServerError, fmt.Sprintf("failed to list persons: %v", err))
		return
	}

	response := interfaces.PersonsListResponse{
		Persons: persons,
		Count:   len(persons),
	}

	helpers.SendJSON(w, http.StatusOK, response)
}

// GetPerson returns a specific person by ID
// @Summary Get a person
// @Description Get a person by ID
// @Tags persons
// @Accept json
// @Produce json
// @Param id path string true "Person ID"
// @Success 200 {object} interfaces.PersonResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Security OAuth2Password
// @Router /v1/persons/{id} [get]
func (h *Handler) GetPerson(w http.ResponseWriter, r *http.Request, personID string) {
	ctx := r.Context()

	person, err := h.service.GetPerson(ctx, personID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			helpers.SendError(w, http.StatusNotFound, "person not found")
			return
		}
		helpers.SendError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get person: %v", err))
		return
	}

	response := interfaces.PersonResponse{
		Person: person,
	}

	helpers.SendJSON(w, http.StatusOK, response)
}

// CreatePerson creates a new person
// @Summary Create a person
// @Description Create a new person in the family tree
// @Tags persons
// @Accept json
// @Produce json
// @Param person body interfaces.PersonCreateRequest true "Person data"
// @Success 201 {object} interfaces.PersonResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Security OAuth2Password
// @Router /v1/persons [post]
func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req interfaces.PersonCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.SendError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	person, err := h.service.CreatePerson(ctx, &req)
	if err != nil {
		helpers.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := interfaces.PersonResponse{
		Person:  person,
		Message: "Person created successfully",
	}

	helpers.SendJSON(w, http.StatusCreated, response)
}

// UpdatePerson updates an existing person
// @Summary Update a person
// @Description Update an existing person in the family tree
// @Tags persons
// @Accept json
// @Produce json
// @Param id path string true "Person ID"
// @Param person body interfaces.PersonUpdateRequest true "Person data"
// @Success 200 {object} interfaces.PersonResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Security OAuth2Password
// @Router /v1/persons/{id} [put]
func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request, personID string) {
	ctx := r.Context()

	var req interfaces.PersonUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.SendError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	person, err := h.service.UpdatePerson(ctx, personID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			helpers.SendError(w, http.StatusNotFound, "person not found")
			return
		}
		helpers.SendError(w, http.StatusInternalServerError, fmt.Sprintf("failed to update person: %v", err))
		return
	}

	response := interfaces.PersonResponse{
		Person:  person,
		Message: "Person updated successfully",
	}

	helpers.SendJSON(w, http.StatusOK, response)
}

// DeletePerson deletes a person
// @Summary Delete a person
// @Description Delete a person from the family tree
// @Tags persons
// @Accept json
// @Produce json
// @Param id path string true "Person ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Security OAuth2Password
// @Router /v1/persons/{id} [delete]
func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request, personID string) {
	ctx := r.Context()

	err := h.service.DeletePerson(ctx, personID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			helpers.SendError(w, http.StatusNotFound, "person not found")
			return
		}
		helpers.SendError(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete person: %v", err))
		return
	}

	helpers.SendJSON(w, http.StatusOK, map[string]string{
		"message": "Person deleted successfully",
	})
}
