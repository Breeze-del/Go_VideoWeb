package model

//error

type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type ErroResponse struct {
	HttpSc int
	Error  Err
}

var (
	ErrorRequestBodyParseFailed = ErroResponse{
		HttpSc: 400,
		Error: Err{
			Error:     "Request body is not correct",
			ErrorCode: "001",
		},
	}
	ErrorNotAuthUser = ErroResponse{
		HttpSc: 400,
		Error: Err{
			Error:     "User authentication failed",
			ErrorCode: "002",
		},
	}
)
