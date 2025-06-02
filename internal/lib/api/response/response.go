package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "Ok"
	StatusError = "Error"
)

func Ok() Response {
	return Response{Status: StatusOK}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
