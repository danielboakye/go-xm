package helpers

import (
	"errors"
	"net/http"
)

var (
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInvalidParameters = errors.New("invalid parameters")
	ErrNoRecordFound     = errors.New("no record found")
	ErrDuplicateRecord   = errors.New("duplicate record")
	ErrProcessingFailed  = errors.New("request could not be processed")
	ErrInvalidToken      = errors.New("token is invalid")
	ErrExpiredToken      = errors.New("token is expired")
)

func GetHttpStatusByErr(err error) (status int) {
	switch err {
	case ErrUnauthorized:
		status = http.StatusUnauthorized
	case ErrInvalidToken:
		status = http.StatusUnauthorized
	case ErrExpiredToken:
		status = http.StatusUnauthorized
	case ErrInvalidParameters:
		status = http.StatusUnprocessableEntity
	case ErrNoRecordFound:
		status = http.StatusNotFound
	case ErrDuplicateRecord:
		status = http.StatusConflict
	case ErrProcessingFailed:
		status = http.StatusInternalServerError
	default:
		status = http.StatusInternalServerError
	}
	return
}
