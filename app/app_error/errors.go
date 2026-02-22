package app_error

const ErrCodeOk ErrCode = 0
const ErrCodeUnknown ErrCode = 1

const (
	ErrCodeUserNotExists ErrCode = 10001 + iota
	ErrCodeUserAlreadyExists

	ErrCodeInvalidParameters
	ErrCodeInvalidAuthorizationHeader

	ErrCodeTimeout

	ErrCodeUserPermissionDenied
	ErrCodeUserNotAuthorized
	ErrCodeUserWrongPassword
	ErrCodeUserInvalidToken
	ErrCodeUserWrongTokenType

	ErrCodeQuestionNotFound
	ErrCodeAnswerNotFound
	ErrCodeCommentNotFound
)

const (
	ErrCodeMysql ErrCode = 20001 + iota
	ErrCodeRedis
	ErrCodeUserToken
	ErrCodeEncryption
)

var (
	ErrUserNotExists        = NewInputError("user not exists", ErrCodeUserNotExists, nil)
	ErrUserAlreadyExists    = NewInputError("user already exists", ErrCodeUserAlreadyExists, nil)
	ErrUserNotAuthorized    = NewInputError("user not authorized", ErrCodeUserNotAuthorized, nil)
	ErrUserPermissionDenied = NewInputError("user permission denied", ErrCodeUserPermissionDenied, nil)
	ErrUserInvalidToken     = NewInputError("user invalid token", ErrCodeUserInvalidToken, nil)
	ErrUserWrongTokenType   = NewInputError("user wrong token type", ErrCodeUserWrongTokenType, nil)
	ErrUserWrongPassword    = NewInputError("user wrong password", ErrCodeUserWrongPassword, nil)
	ErrTimeout              = NewInputError("timeout", ErrCodeTimeout, nil)
	ErrQuestionNotFound     = NewInputError("question not found", ErrCodeQuestionNotFound, nil)
	ErrAnswerNotFound       = NewInputError("answer not found", ErrCodeAnswerNotFound, nil)
	ErrCommentNotFound      = NewInputError("comment not found", ErrCodeCommentNotFound, nil)
)
