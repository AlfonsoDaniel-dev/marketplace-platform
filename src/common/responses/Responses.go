package responses

type Response struct {
	Status  string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Responses struct {
	Status  string `json:"code"`
	Message string `json:"message"`
	Data    []Response
}

func NewResponse(code, message string, data interface{}) Response {
	return Response{
		Status:  code,
		Message: message,
		Data:    data,
	}
}

func GenerateResponses(code, message string, Content []any) Responses {

	ResponseContent := []Response{}
	for _, v := range Content {
		response := NewResponse(code, message, v)

		ResponseContent = append(ResponseContent, response)
	}

	return Responses{
		Status:  code,
		Message: message,
		Data:    ResponseContent,
	}
}
