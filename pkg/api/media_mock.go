// pkg/api/media_mock.go
package api

import (
    "context"
    "io"
    "github.com/stretchr/testify/mock"
    "ezzygo/pkg/cms"
)

type MockMediaService struct {
    mock.Mock
}

func (m *MockMediaService) Upload(ctx context.Context, media *cms.Media, file io.Reader) error {
    args := m.Called(ctx, media, file)
    return args.Error(0)
}

func (m *MockMediaService) Get(ctx context.Context, id uint) (*cms.Media, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*cms.Media), args.Error(1)
}

func (m *MockMediaService) List(ctx context.Context, filter cms.MediaFilter) ([]cms.Media, error) {
    args := m.Called(ctx, filter)
    return args.Get(0).([]cms.Media), args.Error(1)
}

func (m *MockMediaService) Delete(ctx context.Context, id uint) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}
