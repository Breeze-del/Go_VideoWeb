package model

//error 写好的一些返回error

type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type ErroResponse struct {
	HttpSc int
	Error  Err
}

var (
	// body解析错误
	ErrorRequestBodyParseFailed = ErroResponse{
		HttpSc: 400,
		Error: Err{
			Error:     "Request body is not correct",
			ErrorCode: "001",
		},
	}
	// 用户检验失败
	ErrorNotAuthUser = ErroResponse{
		HttpSc: 400,
		Error: Err{
			Error:     "User authentication failed",
			ErrorCode: "002",
		},
	}
	// 数据库操作失败
	ErrorDBError = ErroResponse{
		HttpSc: 500,
		Error: Err{
			Error:     "DB ops failed",
			ErrorCode: "003",
		},
	}
	ErrorInternalFailed = ErroResponse{
		HttpSc: 500,
		Error: Err{
			Error:     "Internal service failed",
			ErrorCode: "004",
		},
	}
)
