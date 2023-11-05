package errors

import (
	"errors"
	"net/http"
)

const (
	DefaultCode          int = 1000
	ValidationCode       int = 1001
	MongoCode            int = 1002
	UnauthorizedCode     int = 1003
	PermissionDeniedCode int = 1004

	UserNotFoundCode int = 1010

	TicketNotFoundCode        int = 1020
	TicketHistoryNotFoundCode int = 1021
	TicketCommentNotFoundCode int = 1022
)

var (
	messageErr = map[int]string{
		DefaultCode:               "something went wrong",
		ValidationCode:            "validation failed",
		MongoCode:                 "mongo error",
		UnauthorizedCode:          "unauthorized",
		PermissionDeniedCode:      "permission denied",
		UserNotFoundCode:          "user not found",
		TicketNotFoundCode:        "ticket not found",
		TicketHistoryNotFoundCode: "ticket history not found",
		TicketCommentNotFoundCode: "ticket comment not found",
	}
)

var (
	ErrDefault             = NewInternalErr(http.StatusInternalServerError, DefaultCode)
	ErrValidation          = NewInternalErr(http.StatusBadRequest, ValidationCode)
	ErrMongo               = NewInternalErr(http.StatusInternalServerError, MongoCode)
	ErrUnauthorized        = NewInternalErr(http.StatusUnauthorized, UnauthorizedCode)
	ErrUnauthorizedContext = NewInternalErr(http.StatusUnauthorized, UnauthorizedCode).SetError(errors.New("auth context not found"))
	ErrPermissionDenied    = NewInternalErr(http.StatusForbidden, PermissionDeniedCode)

	ErrUserNotFound          = NewInternalErr(http.StatusNotFound, UserNotFoundCode)
	ErrTicketNotFound        = NewInternalErr(http.StatusNotFound, TicketNotFoundCode)
	ErrTicketHistoryNotFound = NewInternalErr(http.StatusNotFound, TicketHistoryNotFoundCode)
	ErrTicketCommentNotFound = NewInternalErr(http.StatusNotFound, TicketCommentNotFoundCode)
)
