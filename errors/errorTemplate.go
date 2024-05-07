package errors

import (
	"fmt"
)

// ErrorTemplate is a template for constructing instances of a specific error
//
//	Fields
//		ID          - unique identifier for error
//		ErrorType   - The type of the error (Error, Warning, Status)
//		MessageFmt  - The format message for the diagnostic message
//		IsTransient - True if the error is transient
//		HTTPStatus  - HTTP status code for error (default http.StatusInternalServerError)
type ErrorTemplate struct {
	ID         string    `json:"-"`
	Type       ErrorType `json:"-"`
	MessageFmt string    `json:"-"`
	Transient  bool      `json:"-"`
	HTTPStatus int       `json:"-"`
}

// NewErrorTemplate creates an instance of ErrorTemplate representing an error
func NewErrorTemplate(id string, messageFmt string, httpStatus int, isTransient bool) *ErrorTemplate {
	return &ErrorTemplate{
		Type:       EtError,
		MessageFmt: messageFmt,
		ID:         id,
		HTTPStatus: httpStatus,
		Transient:  isTransient,
	}
}

// NewWarningTemplate creates an instance of ErrorTemplate representing a warrning
func NewWarningTemplate(id string, messageFmt string, httpStatus int) *ErrorTemplate {
	return &ErrorTemplate{
		Type:       EtWarning,
		MessageFmt: messageFmt,
		ID:         id,
		HTTPStatus: httpStatus,
	}
}

// NewStatusTemplate creates an instance of ErrorTemplate representing a status
func NewStatusTemplate(id string, messageFmt string, httpStatus int) *ErrorTemplate {
	return &ErrorTemplate{
		Type:       EtStatus,
		MessageFmt: messageFmt,
		ID:         id,
		HTTPStatus: httpStatus,
	}
}

// Instancef creates an instance of Error from an ErrorTemplate with a
// custom message format string and replacement args
func (template *ErrorTemplate) Instancef(fmtString string, args ...interface{}) *Error {
	err := &Error{
		ID:         template.ID,
		Type:       template.Type,
		Transient:  template.Transient,
		HTTPStatus: template.HTTPStatus,
	}

	if len(args) > 0 {
		err.WithMessage(fmt.Sprintf(fmtString, args...))
	} else {
		err.WithMessage(fmtString)
	}

	return err
}

// Instance creates an instance of Error from an ErrorTemplate using the
// default message format string and replacement args
func (template *ErrorTemplate) Instance(args ...interface{}) *Error {
	return template.Instancef(template.MessageFmt, args...)
}
