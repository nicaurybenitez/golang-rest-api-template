// pkg/api/template_test.go
package api

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "ezzygo/pkg/cms"
)

func setupTemplateTest() (*gin.Engine, *MockTemplateService) {
    gin.SetMode(gin.TestMode)
    mockService := new(MockTemplateService)
    router := gin.New()
    api := NewTemplateAPI(mockService)
    api.RegisterRoutes(router)
    return router, mockService
}

func TestTemplateCreate(t *testing.T) {
    router, mockService := setupTemplateTest()
    template := cms.Template{
        Name: "Blog Post",
        Type: "post",
    }

    mockService.On("Create", mock.Anything, &template).Return(nil)

    body, _ := json.Marshal(template)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/templates", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
    mockService.AssertExpectations(t)
}

func TestTemplateGet(t *testing.T) {
    router, mockService := setupTemplateTest()
    template := &cms.Template{
        Name: "Blog Post",
        Type: "post",
    }
    template.ID = 1

    mockService.On("Get", mock.Anything, uint(1)).Return(template, nil)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/templates/1", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    
    var response cms.Template
    _ = json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, template.Name, response.Name)
    mockService.AssertExpectations(t)
}

func TestTemplateList(t *testing.T) {
    router, mockService := setupTemplateTest()
    templates := []cms.Template{
        {Name: "Blog Post", Type: "post"},
        {Name: "Page", Type: "page"},
    }

    filter := cms.TemplateFilter{Type: "post"}
    mockService.On("List", mock.Anything, filter).Return(templates, nil)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/templates?type=post", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    
    var response []cms.Template
    _ = json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, 2, len(response))
    mockService.AssertExpectations(t)
}

func TestTemplateUpdate(t *testing.T) {
    router, mockService := setupTemplateTest()
    template := cms.Template{
        Name: "Updated Template",
        Type: "post",
    }
    template.ID = 1

    mockService.On("Update", mock.Anything, &template).Return(nil)

    body, _ := json.Marshal(template)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("PUT", "/api/v1/templates/1", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    mockService.AssertExpectations(t)
}

func TestTemplateDelete(t *testing.T) {
    router, mockService := setupTemplateTest()
    mockService.On("Delete", mock.Anything, uint(1)).Return(nil)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("DELETE", "/api/v1/templates/1", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNoContent, w.Code)
    mockService.AssertExpectations(t)
}
