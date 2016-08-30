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

	ErrorCodeUnkown                = 5300
	ErrorCodeJsonBuilding          = 5301
	ErrorCodeUrlNotSupported       = 5302
	ErrorCodeDbNotInitlized        = 5303
	ErrorCodeAuthFailed            = 5304
	ErrorCodePermissionDenied      = 5305
	ErrorCodeInvalidParameters     = 5306
	ErrorCodeGetDataItem           = 5307
	ErrorCodeCreateMarket            = 5308
	ErrorCodeGetMarket               = 5309
	ErrorCodeCancelMarket            = 5310
	ErrorCodeQueryMarket             = 5311
	ErrorCodeGetStatistics         = 5312
	ErrorCodeParseJsonFailed       = 5313
	ErrorCodeFailedToConnectRemote = 5314
	ErrorCodeNotOkRemoteResponse   = 5315
	ErrorCodeInvalidRemoteResponse = 5316

	NumErrors = 5999 // about 50k memroy wasted
)

func init() {
	initError(ErrorCodeNone, "OK")
	initError(ErrorCodeUnkown, "unknown error")
	initError(ErrorCodeJsonBuilding, "json building error")

	initError(ErrorCodeUrlNotSupported, "unsupported url")
	initError(ErrorCodeDbNotInitlized, "db is not inited")
	initError(ErrorCodeAuthFailed, "auth failed")
	initError(ErrorCodePermissionDenied, "permission denied")
	initError(ErrorCodeInvalidParameters, "invalid parameters")
	initError(ErrorCodeGetDataItem, "failed to get data item")
	initError(ErrorCodeCreateMarket, "failed to create star")
	initError(ErrorCodeGetMarket, "failed to get star")
	initError(ErrorCodeCancelMarket, "failed to cancel star")
	initError(ErrorCodeQueryMarket, "failed to  query stars")
	initError(ErrorCodeGetStatistics, "failed to get statistics")
	initError(ErrorCodeParseJsonFailed, "parse json failed")
	initError(ErrorCodeFailedToConnectRemote, "failed to connect remote")
	initError(ErrorCodeNotOkRemoteResponse, "remote response is not ok")
	initError(ErrorCodeInvalidRemoteResponse, "remote response error")

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
