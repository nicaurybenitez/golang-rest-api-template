// pkg/api/router_mocks_test.go
package api

import (
    "context"
    "ezzygo/pkg/cache"
    "ezzygo/pkg/database"
    "ezzygo/pkg/storage"
    "github.com/stretchr/testify/mock"
    "go.mongodb.org/mongo-driver/mongo"
    "go.uber.org/zap"
)

type MockDatabase struct {
    mock.Mock
}

func (m *MockDatabase) Connect() error {
    args := m.Called()
    return args.Error(0)
}

func (m *MockDatabase) Close() error {
    args := m.Called()
    return args.Error(0)
}

func (m *MockDatabase) DB() interface{} {
    args := m.Called()
    return args.Get(0)
}

type MockCache struct {
    mock.Mock
}

func (m *MockCache) Set(key string, value interface{}, ttl time.Duration) error {
    args := m.Called(key, value, ttl)
    return args.Error(0)
}

func (m *MockCache) Get(key string) (interface{}, error) {
    args := m.Called(key)
    return args.Get(0), args.Error(1)
}

func (m *MockCache) Delete(key string) error {
    args := m.Called(key)
    return args.Error(0)
}

type MockMongoCollection struct {
    mock.Mock
}

func (m *MockMongoCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
    args := m.Called(ctx, document)
    return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func createMockDatabase() database.Database {
    mockDB := new(MockDatabase)
    mockDB.On("Connect").Return(nil)
    mockDB.On("DB").Return(nil)
    return mockDB
}

func createMockCache() cache.Cache {
    mockCache := new(MockCache)
    mockCache.On("Get", mock.Anything).Return(nil, nil)
    mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)
    return mockCache
}

func createMockMongoCollection() *mongo.Collection {
    mockMongo := new(MockMongoCollection)
    mockMongo.On("InsertOne", mock.Anything, mock.Anything).Return(&mongo.InsertOneResult{}, nil)
    return &mongo.Collection{} // Reemplazar con mock cuando sea necesario
}

func createMockS3Storage() *storage.S3Storage {
    return &storage.S3Storage{
        Bucket: "test-bucket",
        Client: nil, // Mock AWS client si es necesario
    }
}
