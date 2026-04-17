package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"env-manager/internal/models"
	"env-manager/internal/repository"
)

type ProjectHandler struct {
	repo repository.ProjectRepository
}

func NewProjectHandler(repo repository.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{repo}
}

func (h *ProjectHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	projects, err := h.repo.FindAll()
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}
	WriteJSON(w, http.StatusOK, ToResponse(true, "Projects retrieved", projects))
}

func (h *ProjectHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ToResponse(false, "invalid id", nil))
		return
	}

	project, err := h.repo.FindByID(uint(id))
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ToResponse(false, "project not found", nil))
		return
	}
	WriteJSON(w, http.StatusOK, ToResponse(true, "Project found", project))
}

func (h *ProjectHandler) GetEnvVars(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ToResponse(false, "invalid id", nil))
		return
	}

	_, err = h.repo.FindByID(uint(id))
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ToResponse(false, fmt.Sprintf("project with ID %v doesn't exists", id), nil))
		return
	}

	rawEnvVars, err := h.repo.FindEnvVarsByID(uint(id))

	if err != nil {
		WriteJSON(w, http.StatusNotFound, ToResponse(false, "env vars not found", nil))
		return
	}

	envVars, err := DecryptEnvVars(&rawEnvVars)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	WriteJSON(w, http.StatusOK, ToResponse(true, "Env vars found", envVars))
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProjectRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, ToResponse(false, err.Error(), nil))
		return
	}

	project := &models.Project{Name: req.Name, Description: req.Description}

	if err := h.repo.Create(project); err != nil {
		WriteJSON(w, http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}

	WriteJSON(w, http.StatusCreated, ToResponse(true, "Project created", project))
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ToResponse(false, "invalid id", nil))
		return
	}

	project, err := h.repo.FindByID(uint(id))
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ToResponse(false, "project not found", nil))
		return
	}

	var req models.UpdateProjectRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, ToResponse(false, err.Error(), nil))
		return
	}

	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != "" {
		project.Description = req.Description
	}

	if err := h.repo.Update(project); err != nil {
		WriteJSON(w, http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}
	WriteJSON(w, http.StatusOK, ToResponse(true, "Project updated", project))
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ToResponse(false, "invalid id", nil))
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		WriteJSON(w, http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}
	WriteJSON(w, http.StatusOK, ToResponse(true, "project deleted", nil))
}
