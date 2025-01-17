package models

type ValidationErrorResponseModel struct {
	ErrorCode    string                 `json:"errorCode"`
	ErrorMessage string                 `json:"errorMessage"`
	Errors       []ValidationErrorModel `json:"errors"`
} // @name Error

type ValidationErrorModel struct {
	Field        string `json:"field"`
	ErrorMessage string `json:"error"`
	Value        string `json:"value"`
}

func (e *ValidationErrorResponseModel) Error() string {
	return e.ErrorMessage
}
