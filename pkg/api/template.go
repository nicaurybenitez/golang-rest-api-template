// pkg/api/template.go
package api

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "ezzygo/pkg/cms"
)

type TemplateAPI struct {
    service *cms.TemplateService
}

func NewTemplateAPI(service *cms.TemplateService) *TemplateAPI {
    return &TemplateAPI{service: service}
}

func (api *TemplateAPI) RegisterRoutes(router *gin.Engine) {
    templates := router.Group("/api/v1/templates")
    {
        templates.POST("/", api.Create)
        templates.GET("/", api.List)
        templates.GET("/:id", api.Get)
        templates.PUT("/:id", api.Update)
        templates.DELETE("/:id", api.Delete)
    }
}

func (api *TemplateAPI) Create(c *gin.Context) {
    var template cms.Template
    if err := c.ShouldBindJSON(&template); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := api.service.Create(c.Request.Context(), &template); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, template)
}

func (api *TemplateAPI) Get(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    template, err := api.service.Get(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
        return
    }

    c.JSON(http.StatusOK, template)
}

func (api *TemplateAPI) List(c *gin.Context) {
    filter := cms.TemplateFilter{
        Type: c.Query("type"),
    }

    templates, err := api.service.List(c.Request.Context(), filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, templates)
}

func (api *TemplateAPI) Update(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    var template cms.Template
    if err := c.ShouldBindJSON(&template); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    template.ID = uint(id)

    if err := api.service.Update(c.Request.Context(), &template); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, template)
}

func (api *TemplateAPI) Delete(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    if err := api.service.Delete(c.Request.Context(), uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}
