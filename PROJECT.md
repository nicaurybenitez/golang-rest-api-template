# Proyecto CMS Go

## Estructura Base
Basado en: https://github.com/araujo88/ezzygo

## Pasos de Implementación

### 1. Setup Inicial
```bash
# Clonar base
git clone https://github.com/araujo88/ezzygo
cd ezzygo

# Crear nueva estructura CMS
mkdir -p pkg/{cms,storage,frontend}
```

### 2. Estructura de Archivos
```
pkg/
├── cms/
│   ├── content.go     # Gestión contenido
│   ├── template.go    # Sistema plantillas
│   └── media.go       # Gestión media
├── storage/
│   └── s3.go          # AWS S3 integration
└── frontend/          # API para Next.js
```

### 3. Modelos Base

#### Content Model
```go
type Content struct {
    gorm.Model
    Title       string    `json:"title" gorm:"not null"`
    Slug        string    `json:"slug" gorm:"unique;not null"`
    Content     string    `json:"content"`
    Status      string    `json:"status" gorm:"default:'draft'"`
    TemplateID  uint      `json:"template_id"`
    AuthorID    uint      `json:"author_id"`
    PublishedAt time.Time `json:"published_at"`
    MetaData    JSON      `json:"meta_data"`
}
```

#### Template Model
```go
type Template struct {
    gorm.Model
    Name     string `json:"name" gorm:"not null"`
    Content  string `json:"content"`
    Type     string `json:"type"`
    Fields   JSON   `json:"fields"`
}
```

#### Media Model
```go
type Media struct {
    gorm.Model
    Name      string `json:"name"`
    Type      string `json:"type"`
    URL       string `json:"url"`
    Size      int64  `json:"size"`
    Path      string `json:"path"`
    MimeType  string `json:"mime_type"`
}
```

### 4. Servicios

#### ContentService
```go
type ContentService struct {
    db      *gorm.DB
    storage StorageService
    cache   CacheService
}

// Métodos principales:
- Create(ctx context.Context, content *Content) error
- Get(ctx context.Context, id uint) (*Content, error)
- List(ctx context.Context, filter ContentFilter) ([]Content, error)
- Update(ctx context.Context, content *Content) error
- Delete(ctx context.Context, id uint) error
- Publish(ctx context.Context, id uint) error
```

### 5. Configuración S3
```go
type S3Client struct {
    client *s3.Client
    bucket string
}

func NewS3Client(bucket string) (*S3Client, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        return nil, err
    }
    return &S3Client{
        client: s3.NewFromConfig(cfg),
        bucket: bucket,
    }, nil
}
```

## Pasos Siguientes

1. Implementar handlers API para cada servicio
2. Configurar rutas en el router
3. Implementar middleware específico CMS
4. Configurar tests
5. Implementar frontend con Next.js

## Dependencias Necesarias
```go
require (
    gorm.io/gorm
    github.com/gin-gonic/gin
    github.com/aws/aws-sdk-go-v2
    github.com/redis/go-redis/v9
)
```

## Variables de Entorno
```env
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_REGION=
S3_BUCKET=
POSTGRES_HOST=
POSTGRES_DB=
POSTGRES_USER=
POSTGRES_PASSWORD=
REDIS_URL=
```

## Comandos Útiles
```bash
# Setup inicial
make setup

# Ejecutar migraciones
make migrate

# Ejecutar tests
make test

# Ejecutar servidor
make run

# Build
make build
```

## Docker
```yaml
version: '3.8'
services:
  app:
    build: .
    environment:
      - AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
      - AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=${POSTGRES_DB}
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

  redis:
    image: redis:alpine
```

## Tests
1. Implementar tests unitarios para cada servicio
2. Implementar tests de integración
3. Configurar CI/CD con GitHub Actions

## Documentación API
Usar Swagger para documentar endpoints:
- /api/v1/content
- /api/v1/templates
- /api/v1/media

## Frontend (Next.js)
Crear repositorio separado para el frontend con:
1. Panel de administración
2. Editor de contenido
3. Preview en tiempo real
4. SEO optimizado
