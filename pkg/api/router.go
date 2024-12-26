package api

import (
    "context"
    "time"
    "ezzygo/pkg/cache"
    "ezzygo/pkg/database"
    "ezzygo/pkg/middleware"
    "ezzygo/pkg/storage"
    docs "ezzygo/docs"
    "github.com/gin-gonic/gin"
    swaggerfiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    "go.mongodb.org/mongo-driver/mongo"
    "go.uber.org/zap"
    "golang.org/x/time/rate"
)

func NewRouter(
    logger *zap.Logger,
    mongoCollection *mongo.Collection,
    db database.Database,
    redisClient cache.Cache,
    s3Storage *storage.S3Storage,
    ctx *context.Context,
) *gin.Engine {
    // Repositorios existentes
    bookRepository := NewBookRepository(db, redisClient, ctx)
    userRepository := NewUserRepository(db, ctx)

    // Nuevos servicios CMS
    contentAPI := NewContentAPI(db, redisClient)
    templateAPI := NewTemplateAPI(db, redisClient)
    mediaAPI := NewMediaAPI(db, s3Storage)

    r := gin.Default()

    // Middleware existente
    r.Use(ContextMiddleware(bookRepository))
    r.Use(middleware.Logger(logger, mongoCollection))
    
    if gin.Mode() == gin.ReleaseMode {
        r.Use(middleware.Security())
        r.Use(middleware.Xss())
    }
    
    r.Use(middleware.Cors())
    r.Use(middleware.RateLimiter(rate.Every(1*time.Minute), 60))

    // Swagger
    docs.SwaggerInfo.BasePath = "/api/v1"

    // API v1
    v1 := r.Group("/api/v1")
    {
        // Rutas existentes
        v1.GET("/", bookRepository.Healthcheck)
        v1.GET("/books", middleware.APIKeyAuth(), bookRepository.FindBooks)
        v1.POST("/books", middleware.APIKeyAuth(), middleware.JWTAuth(), bookRepository.CreateBook)
        v1.GET("/books/:id", middleware.APIKeyAuth(), bookRepository.FindBook)
        v1.PUT("/books/:id", middleware.APIKeyAuth(), bookRepository.UpdateBook)
        v1.DELETE("/books/:id", middleware.APIKeyAuth(), bookRepository.DeleteBook)
        v1.POST("/login", middleware.APIKeyAuth(), userRepository.LoginHandler)
        v1.POST("/register", middleware.APIKeyAuth(), userRepository.RegisterHandler)

        // Nuevas rutas CMS
        cms := v1.Group("/cms", middleware.APIKeyAuth(), middleware.JWTAuth())
        {
            // Content routes
            content := cms.Group("/content")
            {
                content.POST("/", contentAPI.Create)
                content.GET("/", contentAPI.List)
                content.GET("/:id", contentAPI.Get)
                content.PUT("/:id", contentAPI.Update)
                content.DELETE("/:id", contentAPI.Delete)
                content.POST("/:id/publish", contentAPI.Publish)
            }

            // Template routes
            templates := cms.Group("/templates")
            {
                templates.POST("/", templateAPI.Create)
                templates.GET("/", templateAPI.List)
                templates.GET("/:id", templateAPI.Get)
                templates.PUT("/:id", templateAPI.Update)
                templates.DELETE("/:id", templateAPI.Delete)
            }

            // Media routes
            media := cms.Group("/media")
            {
                media.POST("/upload", mediaAPI.Upload)
                media.GET("/", mediaAPI.List)
                media.GET("/:id", mediaAPI.Get)
                media.DELETE("/:id", mediaAPI.Delete)
            }
        }
    }

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
    
    return r
}

func ContextMiddleware(bookRepository BookRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("appCtx", bookRepository)
        c.Next()
    }
}
