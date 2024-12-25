// pkg/api/template_mock.go
package api

import (
    "context"
    "github.com/stretchr/testify/mock"
    "ezzygo/pkg/cms"
)

type MockTemplateService struct {
    mock.Mock
}

func (m *MockTemplateService) Create(ctx context.Context, template *cms.Template) error {
    args := m.Called(ctx, template)
    return args.Error(0)
}

func (m *MockTemplateService) Get(ctx context.Context, id uint) (*cms.Template, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*cms.Template), args.Error(1)
}

func (m *MockTemplateService) List(ctx context.Context, filter cms.TemplateFilter) ([]cms.Template, error) {
    args := m.Called(ctx, filter)
    return args.Get(0).([]cms.Template), args.Error(1)
}

func (m *MockTemplateService) Update(ctx context.Context, template *cms.Template) error {
    args := m.Called(ctx, template)
    return args.Error(0)
}

func (m *MockTemplateService) Delete(ctx context.Context, id uint) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}
