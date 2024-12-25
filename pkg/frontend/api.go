package frontend

type FrontendAPI struct {
    contentService *cms.ContentService
    cache         *cache.RedisClient
}
