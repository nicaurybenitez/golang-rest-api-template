package cms

type Template struct {
    gorm.Model
    Name     string `json:"name" gorm:"not null"`
    Content  string `json:"content"`
    Type     string `json:"type"` // page, post, custom
    Fields   JSON   `json:"fields"`
}
