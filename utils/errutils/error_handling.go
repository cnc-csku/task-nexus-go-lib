package errutils

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Error struct {
	Status       ErrorStatus
	Message      string
	StackErr     error
	DebugMessage string
	ErrFields    []string
}

func NewError(err error, errStatus ErrorStatus) *Error {
	return &Error{
		Status:   errStatus,
		Message:  err.Error(),
		StackErr: errors.New(err.Error()),
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) WithDebugMessage(debugMessage string) *Error {
	e.DebugMessage = debugMessage
	return e
}

func (e *Error) WithMessage(message string) *Error {
	e.Message = message
	return e
}

func (e *Error) WithFields(fields ...string) *Error {
	e.ErrFields = fields
	return e
}

type RestErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (e *Error) ToEchoError() error {
	return echo.NewHTTPError(e.Status.StatusCode(), RestErrorResponse{
		Status:  e.Status.String(),
		Message: e.Message,
	}).WithInternal(e)
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	var (
		code    = http.StatusInternalServerError
		message interface{}
	)

	switch e := err.(type) {
	case *echo.HTTPError:
		code = e.Code
		message = e.Message
		if e.Internal != nil {
			if internalErr, ok := e.Internal.(*Error); ok {
				message = formatCustomError(internalErr)
			} else {
				message = map[string]interface{}{
					"status":  http.StatusText(code),
					"message": e.Message,
					"error":   e.Internal.Error(),
				}
			}
		}
	case *Error:
		code = e.Status.StatusCode()
		message = formatCustomError(e)
	default:
		message = map[string]interface{}{
			"status":  http.StatusText(code),
			"message": err.Error(),
		}
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			c.NoContent(code)
		} else {
			c.JSON(code, message)
		}
	}
}

func formatCustomError(err *Error) map[string]interface{} {
	if err.ErrFields != nil {
		return map[string]interface{}{
			"status":  err.Status.String(),
			"message": err.Message,
			"fields":  err.ErrFields,
		}
	}
	return map[string]interface{}{
		"status":  err.Status.String(),
		"message": err.Message,
	}
}

type StackTracer interface {
	StackTrace() errors.StackTrace
}

type errorField struct {
	Type  string `json:"type"`
	Stack string `json:"stack"`
}

func GetStackField(err error) errorField {
	var stack string

	serr, ok := err.(StackTracer)
	if ok {
		// Capture stack trace using github.com/pkg/errors package
		st := serr.StackTrace()
		stack = fmt.Sprintf("%+v", st)
		if len(stack) > 0 && stack[0] == '\n' {
			stack = stack[1:]
		}
	} else {
		// Capture stack trace using runtime package
		stackBuf := make([]byte, 1024)
		stackSize := runtime.Stack(stackBuf, false)
		stack = string(stackBuf[:stackSize])
	}

	return errorField{
		Type:  reflect.TypeOf(err).String(),
		Stack: stack,
	}
}
