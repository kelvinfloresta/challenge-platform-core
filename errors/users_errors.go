package core_errors

var (
	PasswordNotMatch         Conflict   = "PASSWORD_NOT_MATCH"
	UserNotFound             NotFound   = "USER_NOT_FOUND"
	ErrDuplicatedEmail       Conflict   = "DUPLICATED_EMAIL"
	ErrDuplicatedLogin       Conflict   = "DUPLICATED_LOGIN"
	ErrEmailAndDocumentEmpty BadRequest = "EMAIL_AND_DOCUMENT_EMPTY"
	ErrEmailAndPhoneEmpty    BadRequest = "EMAIL_AND_PHONE_EMPTY"
)
