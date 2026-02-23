package dao

import (
	"context"
	"errors"
	"my_zhihu_backend/app/app_error"
	"my_zhihu_backend/app/model"
	"my_zhihu_backend/app/response"

	"gorm.io/gorm"
)

type ArticleDAO struct {
	db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) *ArticleDAO {
	return &ArticleDAO{db: db}
}

func (a *ArticleDAO) PostNewQuestion(ctx context.Context, question *model.Question) app_error.AppError {
	err := gorm.G[model.Question](a.db).Create(ctx, question)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	return nil
}

func (a *ArticleDAO) DeleteQuestion(ctx context.Context, userId int64, questionId int64) app_error.AppError {
	rowsAffected, err := gorm.G[model.Question](a.db).Where("id = ? and author_id = ?", questionId, userId).Delete(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	if rowsAffected == 0 {
		return app_error.ErrUserPermissionDenied
	}
	return nil
}

func (a *ArticleDAO) GetQuestion(ctx context.Context, questionId int64) (*model.Question, app_error.AppError) {
	question, err := gorm.G[model.Question](a.db).Where("id = ?", questionId).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.ErrQuestionNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, app_error.ErrTimeout.WithError(err)
		}
		return nil, app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	return &question, nil
}

func (a *ArticleDAO) UpdateQuestion(ctx context.Context, userId int64, questionId int64, updateData map[string]interface{}) (*model.Question, app_error.AppError) {
	res := a.db.WithContext(ctx).Model(&model.Question{}).Where("id = ? and author_id = ?", questionId, userId).Select("title", "body").Updates(updateData)
	if res.Error != nil {
		if errors.Is(res.Error, context.DeadlineExceeded) {
			return nil, app_error.ErrTimeout.WithError(res.Error)
		}
		return nil, app_error.NewInternalError(app_error.ErrCodeMysql, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, app_error.ErrUserPermissionDenied
	}
	return a.GetQuestion(ctx, questionId)
}

func (a *ArticleDAO) ListQuestions(ctx context.Context, page, size int, keywords string) ([]response.ArticleSearchResult, app_error.AppError) {
	var results []response.ArticleSearchResult

	rawSql := `
	select id, left(title, 50) as part, 'question' as type, Match(title, content) Against(? IN NATURAL LANGUAGE MODE) as score
	from questions
	where Match(title, content) Against(? IN NATURAL LANGUAGE MODE) and is_available = ?
	union all
	select id, left(content, 50) as part, 'answer' as type, Match(content) Against(? IN NATURAL LANGUAGE MODE) as score
    from answers
	where Match(content) Against(? IN NATURAL LANGUAGE MODE) and is_available = ?
	order by score desc 
	limit ? offset ?
`

	err := a.db.WithContext(ctx).Raw(rawSql, keywords, keywords, true, keywords, keywords, true, size, (page-1)*size).Scan(&results).Error

	if err != nil {
		return nil, app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}

	return results, nil
}

func (a *ArticleDAO) PostNewAnswer(ctx context.Context, answer *model.Answer) app_error.AppError {
	err := gorm.G[model.Answer](a.db).Create(ctx, answer)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	return nil
}

func (a *ArticleDAO) UpdateAnswer(ctx context.Context, userId int64, answerId int64, newContent string) app_error.AppError {
	rowsAffected, err := gorm.G[model.Answer](a.db).Where("id = ? and author_id = ?", answerId, userId).Update(ctx, "content", newContent)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	if rowsAffected == 0 {
		return app_error.ErrUserPermissionDenied
	}
	return nil
}

func (a *ArticleDAO) DeleteAnswer(ctx context.Context, userId int64, answerId int64) app_error.AppError {
	rowsAffected, err := gorm.G[model.Answer](a.db).Where("id = ? and author_id = ?", answerId, userId).Delete(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	if rowsAffected == 0 {
		return app_error.ErrUserPermissionDenied
	}
	return nil
}

func (a *ArticleDAO) GetAnswer(ctx context.Context, answerId int64) (*model.Answer, app_error.AppError) {
	answer, err := gorm.G[model.Answer](a.db).Where("id = ?", answerId).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.ErrAnswerNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, app_error.ErrTimeout.WithError(err)
		}
		return nil, app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	return &answer, nil
}

func (a *ArticleDAO) ListAnswers(ctx context.Context, questionId int64, page, size int) ([]model.Answer, int64, app_error.AppError) {
	var answers []model.Answer
	var total int64

	query := gorm.G[model.Answer](a.db).Where("question_id = ? AND is_available = ?", questionId, true)

	total, err := query.Count(ctx, "id")
	if err != nil {
		return nil, 0, app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}

	answers, err = query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(ctx)
	if err != nil {
		return nil, 0, app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}

	return answers, total, nil
}

func (a *ArticleDAO) PostNewComment(ctx context.Context, comment *model.Comment) app_error.AppError {
	err := gorm.G[model.Comment](a.db).Create(ctx, comment)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	return nil
}

func (a *ArticleDAO) UpdateComment(ctx context.Context, userId int64, commentId int64, newContent string) app_error.AppError {
	rowsAffected, err := gorm.G[model.Comment](a.db).Where("id = ? and author_id = ?", commentId, userId).Update(ctx, "content", newContent)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	if rowsAffected == 0 {
		return app_error.ErrUserPermissionDenied
	}
	return nil
}

func (a *ArticleDAO) DeleteComment(ctx context.Context, userId int64, commentId int64) app_error.AppError {
	rowsAffected, err := gorm.G[model.Comment](a.db).Where("id = ? and author_id = ?", commentId, userId).Delete(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return app_error.ErrTimeout.WithError(err)
		}
		return app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	if rowsAffected == 0 {
		return app_error.ErrUserPermissionDenied
	}
	return nil
}

func (a *ArticleDAO) GetComment(ctx context.Context, commentId int64) (*model.Comment, app_error.AppError) {
	comment, err := gorm.G[model.Comment](a.db).Where("id = ?", commentId).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.ErrCommentNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, app_error.ErrTimeout.WithError(err)
		}
		return nil, app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}
	return &comment, nil
}

func (a *ArticleDAO) ListComments(ctx context.Context, answerId int64, page, size int) ([]model.Comment, int64, app_error.AppError) {
	var comments []model.Comment

	query := gorm.G[model.Comment](a.db).Where("answer_id = ?", answerId)

	total, err := query.Count(ctx, "id")
	if err != nil {
		return nil, 0, app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}

	comments, err = query.Offset((page - 1) * size).Limit(size).Order("created_at ASC").Find(ctx)
	if err != nil {
		return nil, 0, app_error.NewInternalError(app_error.ErrCodeMysql, err)
	}

	return comments, total, nil
}
