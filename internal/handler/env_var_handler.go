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

	envVars, total, err := h.envVarsRepo.FindAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": envVars, "total": total, "page": page, "limit": limit})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	if err := h.envVarsRepo.Create(envVar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": envVar})
}

func (h *EnvVarHandlerand) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	envVar, err := h.envVarsRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "env var not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": envVar})
}

func (h *EnvVarHandlerand) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	envVar, err := h.envVarsRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "env var not found"})
		return
	}

	var req models.UpdateEnvVarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Key != "" {
		envVar.Key = req.Key
	}
	if req.Value != "" {
		envVar.Value = req.Value
	}

	if err := h.envVarsRepo.Update(envVar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": envVar})
}

func (h *EnvVarHandlerand) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.envVarsRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "env var deleted"})
}
