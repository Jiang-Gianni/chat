package web

type ctxKey string

var (
	UsernameCtxKey = ctxKey("usernameCtxKey")
	UserIDCtxKey   = ctxKey("userIDCtxKey")
)
