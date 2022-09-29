package v1

type ErrCode int

const (
	ErrCodeNone ErrCode = iota
	ErrCodeInvalidArgument
	ErrCodeInternal
	ErrCodeUnauthenticated
	ErrCodeNoAccess
)
