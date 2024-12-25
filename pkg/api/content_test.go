// pkg/api/content_test.go
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

func setupTest() (*gin.Engine, *MockContentService) {
    gin.SetMode(gin.TestMode)
    mockService := new(MockContentService)
    router := gin.New()
    api := NewContentAPI(mockService)
    api.RegisterRoutes(router)
    return router, mockService
}

func TestCreate(t *testing.T) {
    router, mockService := setupTest()
    content := cms.Content{
        Title: "Test Content",
        Slug:  "test-content",
    }

    mockService.On("Create", mock.Anything, &content).Return(nil)

    body, _ := json.Marshal(content)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/content", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
    mockService.AssertExpectations(t)
}

func TestGet(t *testing.T) {
    router, mockService := setupTest()
    content := &cms.Content{
        Title: "Test Content",
        Slug:  "test-content",
    }
    content.ID = 1

    mockService.On("Get", mock.Anything, uint(1)).Return(content, nil)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/content/1", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    
    var response cms.Content
    _ = json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, content.Title, response.Title)
    mockService.AssertExpectations(t)
}

func TestList(t *testing.T) {
    router, mockService := setupTest()
    contents := []cms.Content{
        {Title: "Content 1", Slug: "content-1"},
        {Title: "Content 2", Slug: "content-2"},
    }

    filter := cms.ContentFilter{Status: "published"}
    mockService.On("List", mock.Anything, filter).Return(contents, nil)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/content?status=published", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    
    var response []cms.Content
    _ = json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, 2, len(response))
    mockService.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
    router, mockService := setupTest()
    content := cms.Content{
        Title: "Updated Content",
        Slug:  "updated-content",
    }
    content.ID = 1

    mockService.On("Update", mock.Anything, &content).Return(nil)

    body, _ := json.Marshal(content)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("PUT", "/api/v1/content/1", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    mockService.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
    router, mockService := setupTest()
    mockService.On("Delete", mock.Anything, uint(1)).Return(nil)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("DELETE", "/api/v1/content/1", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNoContent, w.Code)
    mockService.AssertExpectations(t)
}

func TestPublish(t *testing.T) {
    router, mockService := setupTest()
    mockService.On("Publish", mock.Anything, uint(1)).Return(nil)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/content/1/publish", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    mockService.AssertExpectations(t)
}
