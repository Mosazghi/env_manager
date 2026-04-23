package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"env-manager/internal/models"
	"env-manager/internal/repository"
)

type EnvVarHandlerand struct {
	projectsRepo repository.ProjectRepository
	envVarsRepo  repository.EnvVarRepository
}

func NewEnvVarHandler(projectsRepo repository.ProjectRepository, envVarsRepo repository.EnvVarRepository) *EnvVarHandlerand {
	return &EnvVarHandlerand{projectsRepo, envVarsRepo}
}

func (h *EnvVarHandlerand) GetAll(w http.ResponseWriter, r *http.Request) {
	pageRaw := r.URL.Query().Get("page")
	if pageRaw == "" {
		pageRaw = "1"
	}
	limitRaw := r.URL.Query().Get("limit")
	if limitRaw == "" {
		limitRaw = "10"
	}

	page, _ := strconv.Atoi(pageRaw)
	limit, _ := strconv.Atoi(limitRaw)

	envVars, _, err := h.envVarsRepo.FindAll(page, limit)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	envVars, err = DecryptEnvVars(&envVars)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	WriteJSON(w, http.StatusOK, ToResponse(true, "Env vars retrieved", envVars))
}

func (h *EnvVarHandlerand) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEnvVarRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	envVar := &models.EnvVar{ProjectID: req.ProjectID, Key: req.Key, EncryptedVal: req.Value}

	// check if project exists
	_, err := h.projectsRepo.FindByID(uint(req.ProjectID))
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ToResponse(false, "project not found", nil))
		return
	}

	enc, err := EncryptValue(req.Value)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}

	envVar.EncryptedVal = enc

	if err := h.envVarsRepo.Create(envVar); err != nil {
		WriteJSON(w, http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}
	WriteJSON(w, http.StatusCreated, ToResponse(true, "Env var created", envVar))
}

func (h *EnvVarHandlerand) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ToResponse(false, "invalid id", nil))
		return
	}

	envVar, err := h.envVarsRepo.FindByID(uint(id))
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ToResponse(false, "env var not found", nil))
		return
	}

	dec, err := DecryptValue(envVar.EncryptedVal)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ToResponse(false, err.Error(), nil))
		return
	}
	envVar.Value = string(dec)
	WriteJSON(w, http.StatusOK, ToResponse(true, "Env var found", envVar))
}

func (h *EnvVarHandlerand) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ToResponse(false, "invalid id", nil))
		return
	}

	envVar, err := h.envVarsRepo.FindByID(uint(id))
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ToResponse(false, "env var not found", nil))
		return
	}

	var req models.UpdateEnvVarRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, ToResponse(false, err.Error(), nil))
		return
	}

	if req.Key != "" {
		envVar.Key = req.Key
	}
	if req.Value != "" {
		enc, err := EncryptValue(req.Value)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
			return
		}
		envVar.EncryptedVal = enc
	}

	if err := h.envVarsRepo.Update(envVar); err != nil {
		WriteJSON(w, http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}
	WriteJSON(w, http.StatusOK, ToResponse(true, "Env var updated", envVar))
}

func (h *EnvVarHandlerand) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ToResponse(false, "invalid id", nil))
		return
	}

	_, err = h.envVarsRepo.FindByID(uint(id))
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ToResponse(false, fmt.Sprintf("env var with ID %v doesn't exists", id), nil))
		return
	}

	if err := h.envVarsRepo.Delete(uint(id)); err != nil {
		WriteJSON(w, http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}
	WriteJSON(w, http.StatusOK, ToResponse(true, "env var deleted", nil))
}
