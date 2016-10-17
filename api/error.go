package api

import (
	"fmt"
)

type Error struct {
	code    uint
	message string
}

var (
	Errors            [NumErrors]*Error
	ErrorNone         *Error
	ErrorUnkown       *Error
	ErrorJsonBuilding *Error
)

const (
	ErrorCodeAny = -1 // this is for testing only
	
	ErrorCodeNone = 0

	ErrorCodeUnkown                = 2300
	ErrorCodeJsonBuilding          = 2301
	ErrorCodeParseJsonFailed       = 2302
	ErrorCodeUrlNotSupported       = 2303
	ErrorCodeDbNotInitlized        = 2304
	ErrorCodeAuthFailed            = 2305
	ErrorCodePermissionDenied      = 2306
	ErrorCodeInvalidParameters     = 2307
	ErrorCodeCreateApp             = 2308
	ErrorCodeDeleteApp             = 2309
	ErrorCodeModifyApp             = 2310
	ErrorCodeGetApp                = 2311
	ErrorCodeQueryApps             = 2312



	NumErrors = 2500 // about 20k memroy wasted
)

func init() {
	initError(ErrorCodeNone, "OK")
	initError(ErrorCodeUnkown, "unknown error")
	initError(ErrorCodeJsonBuilding, "json building error")
	initError(ErrorCodeParseJsonFailed, "parse json failed")

	initError(ErrorCodeUrlNotSupported, "unsupported url")
	initError(ErrorCodeDbNotInitlized, "db is not inited")
	initError(ErrorCodeAuthFailed, "auth failed")
	initError(ErrorCodePermissionDenied, "permission denied")
	initError(ErrorCodeInvalidParameters, "invalid parameters")

	initError(ErrorCodeCreateApp, "failed to create app")
	initError(ErrorCodeDeleteApp, "failed to delete app")
	initError(ErrorCodeModifyApp, "failed to modify app")
	initError(ErrorCodeGetApp, "failed to retrieve app")
	initError(ErrorCodeQueryApps, "failed to query apps")

	ErrorNone = GetError(ErrorCodeNone)
	ErrorUnkown = GetError(ErrorCodeUnkown)
	ErrorJsonBuilding = GetError(ErrorCodeJsonBuilding)
}

func initError(code uint, message string) {
	if code < NumErrors {
		Errors[code] = newError(code, message)
	}
}

func GetError(code uint) *Error {
	if code > NumErrors {
		return Errors[ErrorCodeUnkown]
	}

	return Errors[code]
}

func GetError2(code uint, message string) *Error {
	e := GetError(code)
	if e == nil {
		return newError(code, message)
	} else {
		return newError(code, fmt.Sprintf("%s (%s)", e.message, message))
	}
}

func newError(code uint, message string) *Error {
	return &Error{code: code, message: message}
}

func newUnknownError(message string) *Error {
	return &Error{
		code:    ErrorCodeUnkown,
		message: message,
	}
}

func newInvalidParameterError(paramName string) *Error {
	return &Error{
		code:    ErrorCodeInvalidParameters,
		message: fmt.Sprintf("%s: %s", GetError(ErrorCodeInvalidParameters).message, paramName),
	}
}
