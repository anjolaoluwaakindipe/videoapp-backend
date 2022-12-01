package errs

import "net/http"

type AppErrs struct {
	Code    int
	Message interface{}
}

func NewBadRequestError(message interface{}) AppErrs {
	return AppErrs{Message: message, Code: http.StatusBadRequest}
}

func NewNotFoundError(message interface{}) AppErrs {
	return AppErrs{Message: message, Code: http.StatusNotFound}
}

func NewInternalServerError() AppErrs {
	return AppErrs{Message: "An unexpected error occured.", Code: http.StatusInternalServerError}
}

func NewConflictError(message interface{}) AppErrs {
	return AppErrs{Message: message, Code: http.StatusConflict}
}
