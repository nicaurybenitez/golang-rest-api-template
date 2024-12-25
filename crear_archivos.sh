#!/bin/bash

# Función para crear un archivo con contenido
create_file_with_content() {
    local filepath="$1"
    local content="$2"

    mkdir -p "$(dirname "$filepath")"  # Crear directorio si no existe
    echo "$content" > "$filepath"     # Escribir contenido en el archivo
    echo "Archivo creado: $filepath"  # Confirmar creación
}

# Contenido de los archivos
content_pkg_cms_content="package cms

import (
    \"time\"
    \"gorm.io/gorm\"
)

type Content struct {
    gorm.Model
    Title       string    \`json:\"title\" gorm:\"not null\"\`
    Slug        string    \`json:\"slug\" gorm:\"unique;not null\"\`
    Content     string    \`json:\"content\"\`
    Status      string    \`json:\"status\" gorm:\"default:'draft'\"\` // draft, published, archived
    TemplateID  uint      \`json:\"template_id\"\`
    AuthorID    uint      \`json:\"author_id\"\`
    PublishedAt time.Time \`json:\"published_at\"\`
    MetaData    JSON     \`json:\"meta_data\"\`
}

type ContentService struct {
    db *gorm.DB
    storage StorageService
}

func NewContentService(db *gorm.DB, storage StorageService) *ContentService {
    return &ContentService{
        db: db,
        storage: storage,
    }
}"

content_pkg_cms_template="package cms

type Template struct {
    gorm.Model
    Name     string \`json:\"name\" gorm:\"not null\"\`
    Content  string \`json:\"content\"\`
    Type     string \`json:\"type\"\` // page, post, custom
    Fields   JSON   \`json:\"fields\"\`
}"

content_pkg_cms_media="package cms

type Media struct {
    gorm.Model
    Name      string \`json:\"name\"\`
    Type      string \`json:\"type\"\`
    URL       string \`json:\"url\"\`
    Size      int64  \`json:\"size\"\`
    Path      string \`json:\"path\"\`
    MimeType  string \`json:\"mime_type\"\`
}"

content_pkg_storage_s3="package storage

import (
    \"context\"
    \"github.com/aws/aws-sdk-go-v2/service/s3\"
)

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
}"

content_pkg_frontend_api="package frontend

type FrontendAPI struct {
    contentService *cms.ContentService
    cache         *cache.RedisClient
}"

# Crear archivos con su contenido
create_file_with_content "pkg/cms/content.go" "$content_pkg_cms_content"
create_file_with_content "pkg/cms/template.go" "$content_pkg_cms_template"
create_file_with_content "pkg/cms/media.go" "$content_pkg_cms_media"
create_file_with_content "pkg/storage/s3.go" "$content_pkg_storage_s3"
create_file_with_content "pkg/frontend/api.go" "$content_pkg_frontend_api"

echo "Todos los archivos han sido creados exitosamente."

