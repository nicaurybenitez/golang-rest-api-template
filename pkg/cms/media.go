package cms

type Media struct {
    gorm.Model
    Name      string `json:"name"`
    Type      string `json:"type"`
    URL       string `json:"url"`
    Size      int64  `json:"size"`
    Path      string `json:"path"`
    MimeType  string `json:"mime_type"`
}
