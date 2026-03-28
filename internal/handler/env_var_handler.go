package handler

import (
	"env-manager/internal/models"
	"env-manager/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EnvVarHandlerand struct {
	projectsRepo repository.ProjectRepository
	envVarsRepo  repository.EnvVarRepository
}

func NewEnvVarHandler(projectsRepo repository.ProjectRepository, envVarsRepo repository.EnvVarRepository) *EnvVarHandlerand {
	return &EnvVarHandlerand{projectsRepo, envVarsRepo}
}

func (h *EnvVarHandlerand) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	envVars, _, err := h.envVarsRepo.FindAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ToResponse(true, "Env vars retrieved", envVars))
}

func (h *EnvVarHandlerand) Create(c *gin.Context) {
	var req models.CreateEnvVarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	envVar := &models.EnvVar{ProjectID: req.ProjectID, Key: req.Key, Value: req.Value}

	// check if project exists
	_, err := h.projectsRepo.FindByID(uint(req.ProjectID))
	if err != nil {
		c.JSON(http.StatusNotFound, ToResponse(false, "project not found", nil))
		return
	}

	if err := h.envVarsRepo.Create(envVar); err != nil {
		c.JSON(http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}
	c.JSON(http.StatusCreated, ToResponse(true, "Env var created", envVar))
}

func (h *EnvVarHandlerand) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ToResponse(false, "invalid id", nil))
		return
	}

	envVar, err := h.envVarsRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ToResponse(false, "env var not found", nil))
		return
	}
	c.JSON(http.StatusOK, ToResponse(true, "Env var found", envVar))
}

func (h *EnvVarHandlerand) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ToResponse(false, "invalid id", nil))
		return
	}

	envVar, err := h.envVarsRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ToResponse(false, "env var not found", nil))
		return
	}

	var req models.UpdateEnvVarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ToResponse(false, err.Error(), nil))
		return
	}

	if req.Key != "" {
		envVar.Key = req.Key
	}
	if req.Value != "" {
		envVar.Value = req.Value
	}

	if err := h.envVarsRepo.Update(envVar); err != nil {
		c.JSON(http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, ToResponse(true, "Env var updated", envVar))
}

func (h *EnvVarHandlerand) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ToResponse(false, "invalid id", nil))
		return
	}

	if err := h.envVarsRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, ToResponse(true, "env var deleted", nil))
}
