package request

type PostNewQuestionRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateQuestionRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PostNewAnswerRequest struct {
	QuestionId int64  `json:"question_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

type UpdateAnswerRequest struct {
	Content string `json:"content"`
}

type PostNewCommentRequest struct {
	AnswerId int64  `json:"answer_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
	ParentId *int64 `json:"parent_id,omitempty"`
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}

type ListQuestionsRequest struct {
	Page     int    `form:"page,default=1"`
	Size     int    `form:"size,default=20"`
	Keywords string `form:"keywords"`
}

type ListAnswersRequest struct {
	QuestionId int64 `form:"question_id" binding:"required"`
	Page       int   `form:"page,default=1"`
	Size       int   `form:"size,default=20"`
}

type ListCommentsRequest struct {
	AnswerId int64 `form:"answer_id" binding:"required"`
	Page     int   `form:"page,default=1"`
	Size     int   `form:"size,default=20"`
}
