package response

type ArticleSearchResult struct {
	ID    int64   `json:"id"`
	Part  string  `json:"part"` // 部分内容
	Type  string  `json:"type"` // question | answer
	Score float64 `json:"score"`
}

type QuestionResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	AuthorId    int64  `json:"author_id"`
	IsAvailable bool   `json:"is_available"`
	UpdatedAt   string `json:"updated_at"`
}

type ArticleSearchResponse struct {
	Total   int                   `json:"total"`
	Page    int                   `json:"page"`
	Size    int                   `json:"size"`
	Records []ArticleSearchResult `json:"records"`
}
