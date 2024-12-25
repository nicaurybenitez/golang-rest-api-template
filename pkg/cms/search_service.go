// pkg/cms/search_service.go
package cms

import (
    "context"
    "strings"
    "gorm.io/gorm"
)

type SearchService struct {
    db *gorm.DB
}

type SearchResult struct {
    Type      string      `json:"type"`      // content, template
    ID        uint        `json:"id"`
    Title     string      `json:"title"`
    Content   string      `json:"content"`
    Score     float64     `json:"score"`
    Metadata  interface{} `json:"metadata"`
}

type SearchFilter struct {
    Types     []string    `json:"types"`     // filtrar por tipo
    Status    string      `json:"status"`    // draft, published
    FromDate  *time.Time  `json:"from_date"`
    ToDate    *time.Time  `json:"to_date"`
}

func NewSearchService(db *gorm.DB) *SearchService {
    return &SearchService{db: db}
}

func (s *SearchService) Search(ctx context.Context, query string, filter SearchFilter) ([]SearchResult, error) {
    var results []SearchResult

    // Búsqueda en contenido
    if len(filter.Types) == 0 || contains(filter.Types, "content") {
        var contents []Content
        contentQuery := s.db.WithContext(ctx).
            Where("LOWER(title) LIKE LOWER(?) OR LOWER(content) LIKE LOWER(?)", 
                "%"+query+"%", "%"+query+"%")

        if filter.Status != "" {
            contentQuery = contentQuery.Where("status = ?", filter.Status)
        }
        
        if filter.FromDate != nil {
            contentQuery = contentQuery.Where("created_at >= ?", filter.FromDate)
        }
        
        if filter.ToDate != nil {
            contentQuery = contentQuery.Where("created_at <= ?", filter.ToDate)
        }

        if err := contentQuery.Find(&contents).Error; err != nil {
            return nil, err
        }

        for _, content := range contents {
            score := calculateScore(query, content.Title, content.Content)
            results = append(results, SearchResult{
                Type:    "content",
                ID:      content.ID,
                Title:   content.Title,
                Content: truncateContent(content.Content),
                Score:   score,
                Metadata: map[string]interface{}{
                    "status": content.Status,
                    "author": content.AuthorID,
                },
            })
        }
    }

    // Búsqueda en plantillas
    if len(filter.Types) == 0 || contains(filter.Types, "template") {
        var templates []Template
        templateQuery := s.db.WithContext(ctx).
            Where("LOWER(name) LIKE LOWER(?) OR LOWER(content) LIKE LOWER(?)", 
                "%"+query+"%", "%"+query+"%")

        if err := templateQuery.Find(&templates).Error; err != nil {
            return nil, err
        }

        for _, template := range templates {
            score := calculateScore(query, template.Name, template.Content)
            results = append(results, SearchResult{
                Type:    "template",
                ID:      template.ID,
                Title:   template.Name,
                Content: truncateContent(template.Content),
                Score:   score,
                Metadata: map[string]interface{}{
                    "type":    template.Type,
                    "version": template.Version,
                },
            })
        }
    }

    // Ordenar resultados por score
    sort.Slice(results, func(i, j int) bool {
        return results[i].Score > results[j].Score
    })

    return results, nil
}

func (s *SearchService) Suggest(ctx context.Context, prefix string) ([]string, error) {
    var suggestions []string
    
    // Sugerencias de títulos de contenido
    var contentTitles []string
    if err := s.db.WithContext(ctx).Model(&Content{}).
        Where("LOWER(title) LIKE LOWER(?)", prefix+"%").
        Pluck("title", &contentTitles).Error; err != nil {
        return nil, err
    }
    suggestions = append(suggestions, contentTitles...)

    // Sugerencias de nombres de plantillas
    var templateNames []string
    if err := s.db.WithContext(ctx).Model(&Template{}).
        Where("LOWER(name) LIKE LOWER(?)", prefix+"%").
        Pluck("name", &templateNames).Error; err != nil {
        return nil, err
    }
    suggestions = append(suggestions, templateNames...)

    return unique(suggestions), nil
}

// Funciones auxiliares
func calculateScore(query, title, content string) float64 {
    query = strings.ToLower(query)
    title = strings.ToLower(title)
    content = strings.ToLower(content)

    titleScore := countOccurrences(query, title) * 2.0  // Mayor peso para coincidencias en título
    contentScore := countOccurrences(query, content)

    return titleScore + contentScore
}

func countOccurrences(substr, str string) float64 {
    return float64(strings.Count(str, substr))
}

func truncateContent(content string) string {
    if len(content) > 200 {
        return content[:200] + "..."
    }
    return content
}

func contains(slice []string, str string) bool {
    for _, s := range slice {
        if s == str {
            return true
        }
    }
    return false
}

func unique(slice []string) []string {
    keys := make(map[string]bool)
    var list []string
    for _, entry := range slice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
}
