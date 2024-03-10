package service_errors

type ServiceError struct {
	EndUserMessage string `json:"msg"`
	Err            error  `json:"error"`
}

func (se *ServiceError) Error() string {
	return se.EndUserMessage
}
