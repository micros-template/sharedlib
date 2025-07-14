package utils

type ResponseError struct {
	StatusCode uint16 `json:"status_code"`
	Message    string `json:"message"`
}

type ResponseSuccess struct {
	StatusCode uint16      `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func ReturnResponseError(statusCode uint16, message string) ResponseError {
	return ResponseError{
		StatusCode: statusCode,
		Message:    message,
	}
}
func ReturnResponseSuccess(statusCode uint16, message string, data ...interface{}) ResponseSuccess {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	} else {
		responseData = nil
	}
	return ResponseSuccess{
		StatusCode: statusCode,
		Message:    message,
		Data:       responseData,
	}
}
