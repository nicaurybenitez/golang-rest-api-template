// pkg/api/content.go
package api

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "ezzygo/pkg/cms"
)

type ContentAPI struct {
    service *cms.ContentService
}

func NewContentAPI(service *cms.ContentService) *ContentAPI {
    return &ContentAPI{service: service}
}

func (api *ContentAPI) RegisterRoutes(router *gin.Engine) {
    content := router.Group("/api/v1/content")
    {
        content.POST("/", api.Create)
        content.GET("/", api.List)
        content.GET("/:id", api.Get)
        content.PUT("/:id", api.Update)
        content.DELETE("/:id", api.Delete)
        content.POST("/:id/publish", api.Publish)
    }
}

func (api *ContentAPI) Create(c *gin.Context) {
    var content cms.Content
    if err := c.ShouldBindJSON(&content); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := api.service.Create(c.Request.Context(), &content); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, content)
}

func (api *ContentAPI) Get(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    content, err := api.service.Get(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
        return
    }

    c.JSON(http.StatusOK, content)
}

func (api *ContentAPI) List(c *gin.Context) {
    filter := cms.ContentFilter{
        Status: c.Query("status"),
    }

    contents, err := api.service.List(c.Request.Context(), filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, contents)
}

func (api *ContentAPI) Update(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    var content cms.Content
    if err := c.ShouldBindJSON(&content); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    content.ID = uint(id)

    if err := api.service.Update(c.Request.Context(), &content); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, content)
}

func (api *ContentAPI) Delete(c *gin.Context) {
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

func (api *ContentAPI) Publish(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    if err := api.service.Publish(c.Request.Context(), uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusOK)
}
