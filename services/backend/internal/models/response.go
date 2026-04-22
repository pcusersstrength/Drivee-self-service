package models

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	statusOK    = "OK"
	statusError = "Error"
)

// OK returns a successful response with status "OK" and no error.
func OK() Response {
	return Response{
		Status: statusOK,
	}
}

// Error returns a failed response with status "Error" and provided message.
func Error(msg string) Response {
	return Response{
		Status: statusError,
		Error:  msg,
	}
}

type RegisterResponse struct {
	UserID  int64  `json:"user_id"`
	Message string `json:"message,omitempty"`
}