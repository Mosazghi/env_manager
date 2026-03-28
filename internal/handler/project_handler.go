package handler

import (
	"net/http"
	"strconv"

	"env-manager/internal/models"
	"env-manager/internal/repository"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	repo repository.ProjectRepository
}

func NewProjectHandler(repo repository.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{repo}
}

func (h *ProjectHandler) GetAll(c *gin.Context) {
	projects, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, ToResponse(true, "Projects retrieved", projects))
}

func (h *ProjectHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ToResponse(false, "invalid id", nil))
		return
	}

	project, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ToResponse(false, "project not found", nil))
		return
	}
	c.JSON(http.StatusOK, ToResponse(true, "Project found", project))
}

func (h *ProjectHandler) GetEnvVars(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ToResponse(false, "invalid id", nil))
		return
	}

	envVars, err := h.repo.FindEnvVarsByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ToResponse(false, "env vars not found", nil))
		return
	}
	c.JSON(http.StatusOK, ToResponse(true, "Env vars found", envVars))
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ToResponse(false, err.Error(), nil))
		return
	}

	project := &models.Project{Name: req.Name, Description: req.Description}

	if err := h.repo.Create(project); err != nil {
		c.JSON(http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, ToResponse(true, "Project created", project))
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ToResponse(false, "invalid id", nil))
		return
	}

	project, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ToResponse(false, "project not found", nil))
		return
	}

	var req models.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ToResponse(false, err.Error(), nil))
		return
	}

	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != "" {
		project.Description = req.Description
	}

	if err := h.repo.Update(project); err != nil {
		c.JSON(http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, ToResponse(true, "Project updated", project))
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ToResponse(false, "invalid id", nil))
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ToResponse(false, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, ToResponse(true, "project deleted", nil))
}
