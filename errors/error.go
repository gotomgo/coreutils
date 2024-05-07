package errors

import (
	"fmt"
	"net/http"
)

// ErrorType defines the type of an error and is an alias for uint
type ErrorType uint

const (
	// EtError indicates an error
	EtError ErrorType = 0x80000000 // 1 << 31
	// EtWarning indicates a warning
	EtWarning ErrorType = 0x40000000 // 1 << 30
	// EtReserved2 is reserved for future use
	EtReserved2 ErrorType = 0x20000000 // 1 << 29
	// EtReserved1 is reserved for future use
	EtReserved1 ErrorType = 0x10000000 // 1 << 28
	// EtStatus indicates a status code
	EtStatus ErrorType = 0
	// ErrorTypeMask can be used to isolate the error type bits
	ErrorTypeMask ErrorType = 0xF0000000 // bits 28-31
)

// IsNotError determines if an err value is nil
func IsNotError(err error) bool {
	return err == nil
}

// NoError is an alias for IsNotError
func NoError(err error) bool {
	return err == nil
}

// IsError determines if an error value is non-nil
func IsError(err error) bool {
	return err != nil
}

const (
	// StatusOKID is used when an error is needed for a non-error situation
	StatusOKID = "api/ok"
	// ErrGenericAPIErrorID is used for errors that originate from a code
	// that doesn't use Error (external dependencies generally) but we need
	// an instance of Error
	ErrGenericAPIErrorID = "api/error"
)

var (
	// OK is the status returned when the operation was successful but an error
	// value is required (usually for bulk operations)
	OK = Error{ID: StatusOKID, Type: EtStatus, Message: "The operation was successful", HTTPStatus: http.StatusOK}
)

// ErrorContext represents key/value pairs stored with an error
type ErrorContext map[string]interface{}

// Error is a context rich error value used by Overtone server APIs
//
//		Fields
//	   ID         - Unique identifier for error
//			Type       - the type: {Error, Warning, Status}
//			Message    - A diagnostic message associated with the error
//	   HTTPStatus - Optional HTTP status code for error (default == http.StatusInternalServerError)
//			InnerError - An optional error value present when an error is deriative
//				of another error
//			Context    - Any context key/value pairs stored with the error
//			Code       - represents any externally defined error code associatd with
//									 the error
type Error struct {
	ID         string       `json:"id"`
	Type       ErrorType    `json:"type"`
	Message    string       `json:"message,omitempty"`
	HTTPStatus int          `json:"-"`
	Transient  bool         `json:"transient,omitempty"`
	InnerError *Error       `json:"innerError,omitempty"`
	Context    ErrorContext `json:"context,omitempty"`
	Code       int          `json:"code,omitempty"`
}

// NewErrorFromError creates an instance of Error from an error
//
//	Notes
//		If err is already of type *Error then err is returned, otherwise
//		a new Error is created with APIErrorCode
func NewErrorFromError(err error) (newError *Error) {
	var ok bool

	if err != nil {
		if newError, ok = err.(*Error); !ok {
			newError = &Error{ID: ErrGenericAPIErrorID, Type: EtError, Message: err.Error(), HTTPStatus: http.StatusInternalServerError}
		}
	}

	return
}

// NewErrorWithID creates an instance of Error from an error with a custom id
func NewErrorWithID(errID string, err error) (newError *Error) {
	var ok bool

	if err != nil {
		if newError, ok = err.(*Error); !ok {
			newError = &Error{ID: errID, Type: EtError, Message: err.Error(), HTTPStatus: http.StatusInternalServerError}
		}
	}

	return
}

// NewErrorFormatted creates an error with a formatted error message
func NewErrorFormatted(errID string, msgFormat string, args ...interface{}) *Error {
	return &Error{
		ID:         errID,
		Type:       EtError,
		Message:    fmt.Sprintf(msgFormat, args...),
		HTTPStatus: http.StatusInternalServerError,
	}
}

// Error returns the string representation of the error
func (err *Error) Error() string {
	if err == nil {
		return ""
	}

	errMsg := fmt.Sprintf("[%s]: %s", err.ID, err.Message)

	if err.InnerError != nil {
		errMsg = fmt.Sprintf("%s\n\t-> %s", errMsg, err.InnerError)
	}

	return errMsg
}

// WithMessage sets the message of an error
func (err *Error) WithMessage(message string) *Error {
	if err == nil {
		return nil
	}

	err.Message = message
	return err
}

// WithMessageF sets the message of an error
func (err *Error) WithMessageF(fmtMsg string, args ...interface{}) *Error {
	if err == nil {
		return nil
	}

	err.Message = fmt.Sprintf(fmtMsg, args...)
	return err
}

// WithInner sets the inner error of an error
//
//	Notes
//		if inner is not an instance of *Error, it is converted to one
//		with a base code of APIErrorCode
//
//		If inner == nil, nothing happens
func (err *Error) WithInner(inner error) *Error {
	if err == nil {
		return nil
	}

	if inner != nil {
		err.InnerError = NewErrorFromError(inner)
	}

	return err
}

// WithContext adds a collection of key/value pairs to the error context
func (err *Error) WithContext(context ErrorContext) *Error {
	if err == nil {
		return nil
	}

	if err.Context == nil {
		err.Context = ErrorContext{}
	}

	for k, v := range context {
		err.Context[k] = v
	}

	return err
}

// WithContextValue adds a key/value pair to the error context
func (err *Error) WithContextValue(key string, value interface{}) *Error {
	if err == nil {
		return nil
	}

	if err.Context == nil {
		err.Context = ErrorContext{}
	}
	err.Context[key] = value
	return err
}

// WithTransient marks the error as transient
func (err *Error) WithTransient(isTransient bool) *Error {
	if err == nil {
		return nil
	}

	err.Transient = isTransient
	return err
}

// WithHTTPStatus sets the suggested HTTP status code for the error
func (err *Error) WithHTTPStatus(httpStatus int) *Error {
	if err == nil {
		return nil
	}

	err.HTTPStatus = httpStatus
	return err
}

// IsError determines if the error represents an error
func (err *Error) IsError() bool {
	if err == nil {
		return false
	}

	return err.Type == EtError
}

// IsWarning determines if the error represents a warning
func (err *Error) IsWarning() bool {
	if err == nil {
		return false
	}

	return err.Type == EtWarning
}

// IsStatus determines if an error represents a status
func (err *Error) IsStatus() bool {
	if err == nil {
		return false
	}

	return err.Type == EtStatus
}

// IsTransient determines if the error represents a transient condition
func (err *Error) IsTransient() bool {
	if err == nil {
		return false
	}

	return err.Transient
}

// IsEqual determines if an error is a specific code
func (err *Error) IsEqual(errID string) bool {
	if err == nil {
		return false
	}

	return err.ID == errID
}

// GetHTTPStatusCode determines the appropriate HTTP status code for an error
//
//	Notes
//		If the error does not specify an HTTP status code then
//		http.StatusInternalServerError is returned for errors, and http.StatusOK
//		is returned for warnings and statuses
func (err *Error) GetHTTPStatusCode() (result int) {
	if err == nil {
		return http.StatusOK
	}

	result = err.HTTPStatus

	if result == 0 {
		if err.IsError() {
			result = http.StatusInternalServerError
		} else {
			result = http.StatusOK
		}
	}

	return
}

// GetHTTPStatusCode determines the appropriate HTTP status code for an error
//
//	Notes
//		If the error does not specify an HTTP status code then
//		http.StatusInternalServerError is returned for errors, and http.StatusOK
//		is returned for warnings and statuses
func GetHTTPStatusCode(err error) (result int) {
	if err == nil {
		return http.StatusOK
	}

	err2, ok := err.(*Error)
	if ok {
		result = err2.HTTPStatus

		if result == 0 {
			if err2.IsError() {
				result = http.StatusInternalServerError
			} else {
				// status errors return 200 OK by default
				result = http.StatusOK
			}
		}
	} else {
		result = http.StatusInternalServerError
	}

	return
}

// GetErrorID returns the unique identifier for the error
func GetErrorID(err error) (result string) {
	if err != nil {
		ourError, ok := err.(*Error)
		if ok {
			result = ourError.ID
		}
	}

	return
}

// IsErrorID determines if an error matches an error id
//
//	Notes
//		err may be nil, in which case false is returned
//
//		If err is NOT an instance of *Error then this method returns false
func IsErrorID(err error, errID string) bool {
	if err == nil {
		return false
	}

	return GetErrorID(err) == errID
}

// IsAPIError determines if an error is equivalent to APIGenericErrorID
//
//	Notes
//		err may be nil, in which case false is returned
//
//		When err is NOT an instance of *Error then this method returns true
//		otherwise err.ID == APIGenericErrorID
func IsAPIError(err error) bool {
	if err == nil {
		return false
	}

	return IsErrorID(err, ErrGenericAPIErrorID)
}

// IsStatusError determines if an error is any Status error code
func IsStatusError(err error) bool {
	if err == nil {
		return false
	}

	ourErr := GetError(err)
	if ourErr == nil {
		return false
	}

	return ourErr.IsStatus()
}

// GetError returns *Error if err is that type, otherwise nil
func GetError(err error) *Error {
	if err2, ok := err.(*Error); ok {
		return err2
	}

	return nil
}
