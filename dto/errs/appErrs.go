package errs

import "net/http"

type AppErrs struct {
	Code    int
	Message interface{}
}

// custom Bad Request error that give a 400 status code
func NewBadRequestError(message interface{}) AppErrs {
	return AppErrs{Message: message, Code: http.StatusBadRequest}
}

// custom Status Not Found error that gives a 404 status code
func NewNotFoundError(message interface{}) AppErrs {
	return AppErrs{Message: message, Code: http.StatusNotFound}
}

// Internal Server error that gives a 500 status code
func NewInternalServerError() AppErrs {
	return AppErrs{Message: "An unexpected error occured.", Code: http.StatusInternalServerError}
}

// custom Conflict error that give a 409 status code
func NewConflictError(message interface{}) AppErrs {
	return AppErrs{Message: message, Code: http.StatusConflict}
}
