package utils

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var (
	Err              = Response{1, "Err", nil}
	Success          = Response{0, "Success", nil}
	Alive            = Response{0, "Alive", "pong"}
	OverFrequency    = Response{1, "OverFrequency", "pong"}
	ParamsErr        = Response{2, "ParamsErr", nil}
	UserIndex        = Response{0, "UserIndex", nil}
	UserInvalid      = Response{1, "UserInvalid", nil}
	LoginSuccess     = Response{0, "LoginSuccess", nil}
	LoginFailed      = Response{1, "LoginFailed", nil}
	PasswordNotMatch = Response{1, "PasswordNotMatch", nil}
	PasswordLenFail  = Response{1, "PasswordLenFail", nil}
	RegisterFail     = Response{1, "PasswordLenFail", nil}
	RegisterSuccess  = Response{0, "RegisterSuccess", nil}
	UserExists       = Response{1, "UserExists", nil}
	TokenInvalid     = Response{1, "TokenInvalid", nil}
)
