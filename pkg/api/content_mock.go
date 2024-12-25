// pkg/api/content_mock.go
package api

import (
    "context"
    "github.com/stretchr/testify/mock"
    "ezzygo/pkg/cms"
)

type MockContentService struct {
    mock.Mock
}

func (m *MockContentService) Create(ctx context.Context, content *cms.Content) error {
    args := m.Called(ctx, content)
    return args.Error(0)
}

func (m *MockContentService) Get(ctx context.Context, id uint) (*cms.Content, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*cms.Content), args.Error(1)
}

func (m *MockContentService) List(ctx context.Context, filter cms.ContentFilter) ([]cms.Content, error) {
    args := m.Called(ctx, filter)
    return args.Get(0).([]cms.Content), args.Error(1)
}

func (m *MockContentService) Update(ctx context.Context, content *cms.Content) error {
    args := m.Called(ctx, content)
    return args.Error(0)
}

func (m *MockContentService) Delete(ctx context.Context, id uint) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *MockContentService) Publish(ctx context.Context, id uint) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}
