package controller

import (
	"context"
	"my_zhihu_backend/app/app_error"
	"my_zhihu_backend/app/config"
	"my_zhihu_backend/app/model"
	"my_zhihu_backend/app/request"
	"my_zhihu_backend/app/response"
	"my_zhihu_backend/app/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	service *service.ArticleService
	cfg     config.ReadConfigFunc
}

func NewArticleController(as *service.ArticleService, cfg config.ReadConfigFunc) *ArticleController {
	return &ArticleController{as, cfg}
}

type ctrlFunc[T any] func(c *gin.Context, ctx context.Context, userId model.UserId, req *T) (*response.Response, app_error.AppError)

func (ctrl *ArticleController) PostQuestion(c *gin.Context) {
	doWithUserId(c, ctrl.cfg().Service.Timeout, func(ctx context.Context, userId model.UserId, req *request.PostNewQuestionRequest) (*response.Response, app_error.AppError) {
		question, err := ctrl.service.PostNewQuestion(ctx, userId, req)
		if err != nil {
			return nil, err
		}
		return &response.Response{
			Code: 0,
			Ok:   true,
			Body: response.QuestionResponse{
				ID:          question.ID,
				Title:       question.Title,
				Content:     question.Content,
				AuthorId:    question.AuthorId,
				IsAvailable: question.IsAvailable,
				UpdatedAt:   question.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
			Message: "question posted",
		}, nil
	})
}

func (ctrl *ArticleController) UpdateQuestion(c *gin.Context) {
	doWithUserId(c, ctrl.cfg().Service.Timeout, func(ctx context.Context, userId model.UserId, req *request.UpdateQuestionRequest) (*response.Response, app_error.AppError) {
		id, err := getIdFromParams(c)
		if err != nil {
			return nil, ErrInvalidParameters.WithError(err)
		}

		if q, err := ctrl.service.UpdateQuestion(ctx, userId, id, req); err != nil {
			return nil, err
		} else {
			return &response.Response{
				Code:    0,
				Ok:      true,
				Message: "question updated",
				Body: response.QuestionResponse{
					ID:          q.ID,
					Title:       q.Title,
					Content:     q.Content,
					AuthorId:    q.AuthorId,
					IsAvailable: q.IsAvailable,
					UpdatedAt:   q.UpdatedAt.Format(time.DateTime),
				},
			}, nil
		}
	})

}

func (ctrl *ArticleController) DeleteQuestion(c *gin.Context) {
	doOnlyWithUserId(c, ctrl.cfg().Service.Timeout, func(ctx context.Context, userId model.UserId) (*response.Response, app_error.AppError) {
		id, err := getIdFromParams(c)
		if err != nil {
			return nil, ErrInvalidParameters.WithError(err)
		}
		if err := ctrl.service.DeleteQuestion(ctx, userId, id); err != nil {
			return nil, err
		}
		return &response.Response{
			Code:          0,
			Ok:            true,
			InternalError: false,
			Message:       "question deleted",
			Body:          nil,
		}, nil
	})
}

func (ctrl *ArticleController) GetQuestion(c *gin.Context) {
	doOnlyWithUserId(c, ctrl.cfg().Service.Timeout, func(ctx context.Context, userId model.UserId) (*response.Response, app_error.AppError) {
		id, err := getIdFromParams(c)
		if err != nil {
			return nil, ErrInvalidParameters.WithError(err)
		}

		if question, err := ctrl.service.GetAndIncrQuestion(ctx, id); err != nil {
			return nil, err
		} else {
			return &response.Response{
				Code:          0,
				Ok:            true,
				InternalError: false,
				Message:       "question got",
				Body: response.QuestionResponse{
					ID:          question.ID,
					Title:       question.Title,
					Content:     question.Content,
					AuthorId:    question.AuthorId,
					IsAvailable: question.IsAvailable,
					UpdatedAt:   question.UpdatedAt.Format("2006-01-02 15:04:05"),
				},
			}, nil
		}
	})
}

func (ctrl *ArticleController) ListQuestions(c *gin.Context) {
	var req request.ListQuestionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(ErrInvalidParameters.WithError(err))
		return
	}

	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, total, err := ctrl.service.ListQuestions(timeout, req.Page, req.Size, req.Keywords)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := response.ArticleSearchResponse{
		Total:   total,
		Page:    req.Page,
		Size:    req.Size,
		Records: results,
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    0,
		Ok:      true,
		Body:    resp,
		Message: "questions listed",
	})
}

//func (ctrl *ArticleController) PostAnswer(c *gin.Context) {
//	doWithUserId(c, ctrl.cfg().Service.Timeout, func(ctx context.Context, userId model.UserId, req *request.PostNewAnswerRequest) (*response.Response, app_error.AppError) {
//		answer, err := ctrl.service.PostNewAnswer(ctx, userId, req)
//		if err != nil {
//			return nil, err
//		}
//		return &response.Response{
//			Code: 0,
//			Ok:   true,
//			Body: response{
//				ID:         answer.ID,
//				QuestionId: answer.QuestionId,
//				Content:    answer.Content,
//				AuthorId:   answer.AuthorId,
//				LikeCount:  answer.LikeCount,
//				CreatedAt:  answer.CreatedAt.Format("2006-01-02 15:04:05"),
//				UpdatedAt:  answer.UpdatedAt.Format("2006-01-02 15:04:05"),
//			},
//			Message: "answer posted",
//		}, nil
//	})
//}
//
//func (ctrl *ArticleController) UpdateAnswer(c *gin.Context) {
//	doWithBody(c, func(c *gin.Context, ctx context.Context, userId model.UserId, req *request.UpdateAnswerRequest) (*response.Response, app_error.AppError) {
//		id := c.Param("id")
//		idNum, err := strconv.Atoi(id)
//		if err != nil {
//			return nil, ErrInvalidParameters.WithError(err)
//		}
//		err = ctrl.service.UpdateAnswer(ctx, userId, uint(idNum), req)
//		if err != nil {
//			return nil, err
//		}
//		return &response.Response{
//			Code:    0,
//			Ok:      true,
//			Message: "answer updated",
//		}, nil
//	})
//}
//
//func (ctrl *ArticleController) DeleteAnswer(c *gin.Context) {
//	doWithBody(c, func(c *gin.Context, ctx context.Context, userId model.UserId, req *request.BaseRequest) (*response.Response, app_error.AppError) {
//		id := c.Param("id")
//		idNum, err := strconv.Atoi(id)
//		if err != nil {
//			return nil, ErrInvalidParameters.WithError(err)
//		}
//		if err := ctrl.service.DeleteAnswer(ctx, userId, uint(idNum)); err != nil {
//			return nil, err
//		}
//		return &response.Response{
//			Code:          0,
//			Ok:            true,
//			InternalError: false,
//			Message:       "answer deleted",
//			Body:          nil,
//		}, nil
//	})
//}
//
//func (ctrl *ArticleController) ListAnswers(c *gin.Context) {
//	var req request.ListAnswersRequest
//	if err := c.ShouldBindQuery(&req); err != nil {
//		_ = c.Error(ErrInvalidParameters.WithError(err))
//		return
//	}
//
//	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	answers, total, err := ctrl.service.ListAnswers(timeout, req.QuestionId, req.Page, req.Size)
//	if err != nil {
//		_ = c.Error(err)
//		return
//	}
//
//	var records []response.AnswerResponse
//	for _, a := range answers {
//		records = append(records, response.AnswerResponse{
//			ID:         a.ID,
//			QuestionId: a.QuestionId,
//			Content:    a.Content,
//			AuthorId:   a.AuthorId,
//			LikeCount:  a.LikeCount,
//			CreatedAt:  a.CreatedAt.Format("2006-01-02 15:04:05"),
//			UpdatedAt:  a.UpdatedAt.Format("2006-01-02 15:04:05"),
//		})
//	}
//
//	resp := response.ListAnswersResponse{
//		Total:   total,
//		Page:    req.Page,
//		Size:    req.Size,
//		Records: records,
//	}
//
//	c.JSON(http.StatusOK, response.Response{
//		Code:    0,
//		Ok:      true,
//		Body:    resp,
//		Message: "answers listed",
//	})
//}
//
//func (ctrl *ArticleController) PostComment(c *gin.Context) {
//	doWithBody(c, func(c *gin.Context, ctx context.Context, userId model.UserId, req *request.PostNewCommentRequest) (*response.Response, app_error.AppError) {
//		comment, err := ctrl.service.PostNewComment(ctx, userId, req)
//		if err != nil {
//			return nil, err
//		}
//		return &response.Response{
//			Code: 0,
//			Ok:   true,
//			Body: response.CommentResponse{
//				ID:        comment.ID,
//				AnswerId:  comment.AnswerId,
//				AuthorId:  comment.AuthorId,
//				Content:   comment.Content,
//				ParentId:  comment.ParentId,
//				CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
//				UpdatedAt: comment.UpdatedAt.Format("2006-01-02 15:04:05"),
//			},
//			Message: "comment posted",
//		}, nil
//	})
//}
//
//func (ctrl *ArticleController) UpdateComment(c *gin.Context) {
//	doWithBody(c, func(c *gin.Context, ctx context.Context, userId model.UserId, req *request.UpdateCommentRequest) (*response.Response, app_error.AppError) {
//		id := c.Param("id")
//		idNum, err := strconv.Atoi(id)
//		if err != nil {
//			return nil, ErrInvalidParameters.WithError(err)
//		}
//		err = ctrl.service.UpdateComment(ctx, userId, uint(idNum), req)
//		if err != nil {
//			return nil, err
//		}
//		return &response.Response{
//			Code:    0,
//			Ok:      true,
//			Message: "comment updated",
//		}, nil
//	})
//}
//
//func (ctrl *ArticleController) DeleteComment(c *gin.Context) {
//	doWithBody(c, func(c *gin.Context, ctx context.Context, userId model.UserId, req *request.BaseRequest) (*response.Response, app_error.AppError) {
//		id := c.Param("id")
//		idNum, err := strconv.Atoi(id)
//		if err != nil {
//			return nil, ErrInvalidParameters.WithError(err)
//		}
//		if err := ctrl.service.DeleteComment(ctx, userId, uint(idNum)); err != nil {
//			return nil, err
//		}
//		return &response.Response{
//			Code:          0,
//			Ok:            true,
//			InternalError: false,
//			Message:       "comment deleted",
//			Body:          nil,
//		}, nil
//	})
//}
//
//func (ctrl *ArticleController) ListComments(c *gin.Context) {
//	var req request.ListCommentsRequest
//	if err := c.ShouldBindQuery(&req); err != nil {
//		_ = c.Error(ErrInvalidParameters.WithError(err))
//		return
//	}
//
//	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	comments, total, err := ctrl.service.ListComments(timeout, req.AnswerId, req.Page, req.Size)
//	if err != nil {
//		_ = c.Error(err)
//		return
//	}
//
//	var records []response.CommentResponse
//	for _, c := range comments {
//		records = append(records, response.CommentResponse{
//			ID:        c.ID,
//			AnswerId:  c.AnswerId,
//			AuthorId:  c.AuthorId,
//			Content:   c.Content,
//			ParentId:  c.ParentId,
//			CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
//			UpdatedAt: c.UpdatedAt.Format("2006-01-02 15:04:05"),
//		})
//	}
//
//	resp := response.ListCommentsResponse{
//		Total:   total,
//		Page:    req.Page,
//		Size:    req.Size,
//		Records: records,
//	}
//
//	c.JSON(http.StatusOK, response.Response{
//		Code:    0,
//		Ok:      true,
//		Body:    resp,
//		Message: "comments listed",
//	})
//}
