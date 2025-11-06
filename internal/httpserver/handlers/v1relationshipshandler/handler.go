package v1relationshipshandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rogerwesterbo/familytree/internal/httpserver/helpers"
	"github.com/rogerwesterbo/familytree/internal/services/v1relationshipservice"
	"github.com/rogerwesterbo/familytree/pkg/interfaces"
)

// Handler handles HTTP requests for relationship operations
type Handler struct {
	service *v1relationshipservice.RelationshipService
}

// NewHandler creates a new relationship handler
func NewHandler(service *v1relationshipservice.RelationshipService) *Handler {
	return &Handler{
		service: service,
	}
}

// HandleRelationships routes relationship requests based on HTTP method
// @Summary Relationship operations
// @Description Handle relationship CRUD operations
// @Tags relationships
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security OAuth2Password
// @Router /v1/relationships [get]
// @Router /v1/relationships [post]
// @Router /v1/relationships/{id} [get]
// @Router /v1/relationships/{id} [put]
// @Router /v1/relationships/{id} [delete]
func (h *Handler) HandleRelationships(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if strings.HasPrefix(r.URL.Path, "/v1/relationships/") {
			relationshipID := strings.TrimPrefix(r.URL.Path, "/v1/relationships/")
			if relationshipID != "" {
				h.GetRelationship(w, r, relationshipID)
				return
			}
		}
		h.ListRelationships(w, r)
	case http.MethodPost:
		h.CreateRelationship(w, r)
	case http.MethodPut:
		relationshipID := strings.TrimPrefix(r.URL.Path, "/v1/relationships/")
		if relationshipID == "" {
			helpers.SendError(w, http.StatusBadRequest, "relationship ID is required")
			return
		}
		h.UpdateRelationship(w, r, relationshipID)
	case http.MethodDelete:
		relationshipID := strings.TrimPrefix(r.URL.Path, "/v1/relationships/")
		if relationshipID == "" {
			helpers.SendError(w, http.StatusBadRequest, "relationship ID is required")
			return
		}
		h.DeleteRelationship(w, r, relationshipID)
	default:
		helpers.SendError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// ListRelationships returns all relationships
// @Summary List all relationships
// @Description Get a list of all relationships in the family tree
// @Tags relationships
// @Accept json
// @Produce json
// @Success 200 {object} interfaces.RelationshipsListResponse
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Security OAuth2Password
// @Router /v1/relationships [get]
func (h *Handler) ListRelationships(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	relationships, err := h.service.ListRelationships(ctx)
	if err != nil {
		helpers.SendError(w, http.StatusInternalServerError, fmt.Sprintf("failed to list relationships: %v", err))
		return
	}

	response := interfaces.RelationshipsListResponse{
		Relationships: relationships,
		Count:         len(relationships),
	}

	helpers.SendJSON(w, http.StatusOK, response)
}

// GetRelationship returns a specific relationship by ID
// @Summary Get a relationship
// @Description Get a relationship by ID
// @Tags relationships
// @Accept json
// @Produce json
// @Param id path string true "Relationship ID"
// @Success 200 {object} interfaces.RelationshipResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Security OAuth2Password
// @Router /v1/relationships/{id} [get]
func (h *Handler) GetRelationship(w http.ResponseWriter, r *http.Request, relationshipID string) {
	ctx := r.Context()

	relationship, err := h.service.GetRelationship(ctx, relationshipID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			helpers.SendError(w, http.StatusNotFound, "relationship not found")
			return
		}
		helpers.SendError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get relationship: %v", err))
		return
	}

	response := interfaces.RelationshipResponse{
		Relationship: relationship,
	}

	helpers.SendJSON(w, http.StatusOK, response)
}

// CreateRelationship creates a new relationship
// @Summary Create a relationship
// @Description Create a new relationship in the family tree
// @Tags relationships
// @Accept json
// @Produce json
// @Param relationship body interfaces.RelationshipCreateRequest true "Relationship data"
// @Success 201 {object} interfaces.RelationshipResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Security OAuth2Password
// @Router /v1/relationships [post]
func (h *Handler) CreateRelationship(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req interfaces.RelationshipCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.SendError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	relationship, err := h.service.CreateRelationship(ctx, &req)
	if err != nil {
		helpers.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := interfaces.RelationshipResponse{
		Relationship: relationship,
		Message:      "Relationship created successfully",
	}

	helpers.SendJSON(w, http.StatusCreated, response)
}

// UpdateRelationship updates an existing relationship
// @Summary Update a relationship
// @Description Update an existing relationship in the family tree
// @Tags relationships
// @Accept json
// @Produce json
// @Param id path string true "Relationship ID"
// @Param relationship body interfaces.RelationshipUpdateRequest true "Relationship data"
// @Success 200 {object} interfaces.RelationshipResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Security OAuth2Password
// @Router /v1/relationships/{id} [put]
func (h *Handler) UpdateRelationship(w http.ResponseWriter, r *http.Request, relationshipID string) {
	ctx := r.Context()

	var req interfaces.RelationshipUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.SendError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	relationship, err := h.service.UpdateRelationship(ctx, relationshipID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			helpers.SendError(w, http.StatusNotFound, "relationship not found")
			return
		}
		helpers.SendError(w, http.StatusInternalServerError, fmt.Sprintf("failed to update relationship: %v", err))
		return
	}

	response := interfaces.RelationshipResponse{
		Relationship: relationship,
		Message:      "Relationship updated successfully",
	}

	helpers.SendJSON(w, http.StatusOK, response)
}

// DeleteRelationship deletes a relationship
// @Summary Delete a relationship
// @Description Delete a relationship from the family tree
// @Tags relationships
// @Accept json
// @Produce json
// @Param id path string true "Relationship ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Security OAuth2Password
// @Router /v1/relationships/{id} [delete]
func (h *Handler) DeleteRelationship(w http.ResponseWriter, r *http.Request, relationshipID string) {
	ctx := r.Context()

	err := h.service.DeleteRelationship(ctx, relationshipID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			helpers.SendError(w, http.StatusNotFound, "relationship not found")
			return
		}
		helpers.SendError(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete relationship: %v", err))
		return
	}

	helpers.SendJSON(w, http.StatusOK, map[string]string{
		"message": "Relationship deleted successfully",
	})
}
