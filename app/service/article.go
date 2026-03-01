package service

import (
	"context"
	"my_zhihu_backend/app/app_error"
	"my_zhihu_backend/app/dao"
	"my_zhihu_backend/app/model"
	"my_zhihu_backend/app/request"
	"my_zhihu_backend/app/response"
	"my_zhihu_backend/app/util"

	"gorm.io/gorm"
)

type ArticleService struct {
	dao  *dao.ArticleDAO
	util *util.Util
}

func NewArticleService(db *gorm.DB) *ArticleService {
	aDAO := dao.NewArticleDAO(db)
	u := new(util.Util)
	return &ArticleService{aDAO, u}
}

func (a *ArticleService) PostNewQuestion(ctx context.Context, userId model.UserId, req *request.PostNewQuestionRequest) (*model.Question, app_error.AppError) {
	question := &model.Question{
		ID:          a.util.GenerateSnowflakeID(),
		Title:       req.Title,
		Content:     req.Content,
		AuthorId:    int64(userId),
		IsAvailable: true,
	}
	err := a.dao.PostNewQuestion(ctx, question)
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (a *ArticleService) UpdateQuestion(ctx context.Context, userId model.UserId, questionId int64, req *request.UpdateQuestionRequest) (*model.Question, app_error.AppError) {
	updateData := map[string]interface{}{}
	if req.Title != "" {
		updateData["title"] = req.Title
	}
	if req.Body != "" {
		updateData["body"] = req.Body
	}

	return a.dao.UpdateQuestion(ctx, int64(userId), questionId, updateData)
}

func (a *ArticleService) DeleteQuestion(ctx context.Context, userId model.UserId, questionId int64) app_error.AppError {
	return a.dao.DeleteQuestion(ctx, int64(userId), questionId)
}

func (a *ArticleService) GetQuestion(ctx context.Context, questionId int64) (*model.Question, app_error.AppError) {
	q, err := a.dao.GetQuestion(ctx, questionId)
	if err != nil {
		return nil, err
	}
	if !q.IsAvailable {
		return nil, app_error.ErrUserPermissionDenied
	}
	return q, nil
}

func (a *ArticleService) GetAndIncrQuestion(ctx context.Context, questionId int64) (*model.Question, app_error.AppError) {
	q, err := a.dao.GetQuestion(ctx, questionId)
	if err != nil {
		return nil, err
	}
	if !q.IsAvailable {
		return nil, app_error.ErrUserPermissionDenied
	}
	return q, nil
}

func (a *ArticleService) ListQuestions(ctx context.Context, page, size int, keywords string) ([]response.ArticleSearchResult, int, app_error.AppError) {
	results, err := a.dao.ListQuestions(ctx, page, size, keywords)
	if err != nil {
		return nil, 0, err
	}
	return results, len(results), nil
}

func (a *ArticleService) PostNewAnswer(ctx context.Context, userId model.UserId, req *request.PostNewAnswerRequest) (*model.Answer, app_error.AppError) {
	// 检查问题是否存在
	_, err := a.dao.GetQuestion(ctx, req.QuestionId)
	if err != nil {
		return nil, err
	}

	answer := &model.Answer{
		ID:          a.util.GenerateSnowflakeID(),
		QuestionId:  req.QuestionId,
		AuthorId:    int64(userId),
		Content:     req.Content,
		LikeCount:   0,
		IsAvailable: true,
	}
	err = a.dao.PostNewAnswer(ctx, answer)
	if err != nil {
		return nil, err
	}
	return answer, nil
}

func (a *ArticleService) UpdateAnswer(ctx context.Context, userId model.UserId, answerId int64, req *request.UpdateAnswerRequest) app_error.AppError {
	return a.dao.UpdateAnswer(ctx, int64(userId), answerId, req.Content)
}

func (a *ArticleService) DeleteAnswer(ctx context.Context, userId model.UserId, answerId int64) app_error.AppError {
	return a.dao.DeleteAnswer(ctx, int64(userId), answerId)
}

func (a *ArticleService) ListAnswers(ctx context.Context, questionId int64, page, size int) ([]model.Answer, int64, app_error.AppError) {
	return a.dao.ListAnswers(ctx, questionId, page, size)
}

func (a *ArticleService) PostNewComment(ctx context.Context, userId model.UserId, req *request.PostNewCommentRequest) (*model.Comment, app_error.AppError) {
	// 检查答案是否存在
	_, err := a.dao.GetAnswer(ctx, req.AnswerId)
	if err != nil {
		return nil, err
	}

	// 如果是回复评论，检查父评论是否存在
	if req.ParentId != nil {
		_, err := a.dao.GetComment(ctx, *req.ParentId)
		if err != nil {
			return nil, err
		}
	}

	comment := &model.Comment{
		ID:       a.util.GenerateSnowflakeID(),
		AnswerId: req.AnswerId,
		AuthorId: int64(userId),
		Content:  req.Content,
		ParentId: req.ParentId,
	}
	err = a.dao.PostNewComment(ctx, comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (a *ArticleService) UpdateComment(ctx context.Context, userId model.UserId, commentId int64, req *request.UpdateCommentRequest) app_error.AppError {
	// 检查评论是否存在且属于当前用户
	comment, err := a.dao.GetComment(ctx, commentId)
	if err != nil {
		return err
	}
	if comment.AuthorId != int64(userId) {
		return app_error.ErrUserPermissionDenied
	}

	return a.dao.UpdateComment(ctx, int64(userId), commentId, req.Content)
}

func (a *ArticleService) DeleteComment(ctx context.Context, userId model.UserId, commentId int64) app_error.AppError {
	return a.dao.DeleteComment(ctx, int64(userId), commentId)
}

func (a *ArticleService) ListComments(ctx context.Context, answerId int64, page, size int) ([]model.Comment, int64, app_error.AppError) {
	return a.dao.ListComments(ctx, answerId, page, size)
}
