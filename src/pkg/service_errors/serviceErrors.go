package service_errors

type ServiceErrors struct {
	EndUserMessage string `json:"msg"`
	Err            error  `json:"error"`
	Status         int    `json:"status"`
}

func (se *ServiceErrors) Error() string {
	return se.EndUserMessage
}
