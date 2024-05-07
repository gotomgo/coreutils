package errors

import "net/http"

const (
	// ErrInternalErrorID is the error id used for internal errors
	ErrInternalErrorID = "api/internal"
	// ErrNotImplementedID is the error id used when a feature is not
	// implemented
	ErrNotImplementedID = "api/not-implemented"

	// ErrArgumentExceptionID is the error id used for an invalid argument
	ErrArgumentExceptionID = "api/arg/exception"
	// ErrArgumentNullID is the error id used for a null argument
	ErrArgumentNullID = "api/arg/null"
	// ErrArgumentEmptyID is the error id used for a empty arguments
	ErrArgumentEmptyID = "api/arg/empty"
	// ErrArgumentOutOfRangeID is the error id used for arguments that are
	// out of range
	ErrArgumentOutOfRangeID = "api/arg/range"

	// ErrRequiredValueMissingID is the error id used for validation failures
	// where a required value is not specified
	ErrRequiredValueMissingID = "api/data/required"
	// ErrDataIntegrityID is the error id used for validation failures
	// where data values are inconsistent or invalid
	ErrDataIntegrityID = "api/data/bad"
	// ErrConversionFailedID is the error id used when a conversion does not
	// work as expected
	ErrConversionFailedID = "api/data/failed/conversion"
	// ErrDataValidationFailedID is the error id used when a data validation
	// fails (the error should always wrap a more specific reason)
	ErrDataValidationFailedID = "api/data/failed/validation"

	// ErrAccessDeniedID is the error id used for access denied
	ErrAccessDeniedID = "api/resource/access-denied"
	// ErrResourceNotFoundID is the error id used when something (generically,
	// a resource) cannot be found
	ErrResourceNotFoundID = "api/resource/not-found"
)

var (
	// ErrInternalError is used to create error instances for internal
	// Errors
	//
	//	Notes
	//		The error message has the form:
	//      "an unexpected internal error occurred"
	//
	//		Feel free to override the instance message via WithMessageF, and
	//		provide context on what operation was being performed, including
	//		WithInner(err) for the underlying error
	//
	ErrInternalError = NewErrorTemplate(
		ErrInternalErrorID,
		"an unexpected internal error occurred",
		http.StatusInternalServerError,
		false)

	// ErrArgumentException is used for custom argument errors
	//
	//	Notes
	//		The message has the form:
	//      "the argument, '%s', cannot be %s"
	//
	ErrArgumentException = NewErrorTemplate(
		ErrArgumentExceptionID,
		"the argument, '%s', cannot be %s",
		http.StatusBadRequest,
		false)

	// ErrArgumentNull is used for custom nil / null argument errors
	//
	//	Notes
	//		The message has the form:
	//      "the argument, %s, cannot be null"
	//
	ErrArgumentNull = NewErrorTemplate(
		ErrArgumentNullID,
		"the argument, %s, cannot be null",
		http.StatusBadRequest,
		false)

	// ErrArgumentEmpty is used for custom empty value argument errors
	//
	//	Notes
	//		The message has the form:
	//      "the argument, %s, cannot be null"
	//
	ErrArgumentEmpty = NewErrorTemplate(
		ErrArgumentEmptyID,
		"the argument, %s, is empty",
		http.StatusBadRequest,
		false)

	// ErrArgumentOutOfRange is used for argument errors where a
	// parameter is out-of-range
	//
	//	Notes
	//		The message has the form:
	//      "the argument, '%s', is out of range"
	//
	ErrArgumentOutOfRange = NewErrorTemplate(
		ErrArgumentOutOfRangeID,
		"the argument, '%s', is out of range",
		http.StatusBadRequest,
		false)

	// ErrAccessDeniedTemplate is an error template for ErrCodeAccessDenied
	//
	//	Notes
	//		The error message has the form:
	//      "access to the requested resource is denied"
	//
	ErrAccessDenied = NewErrorTemplate(
		ErrAccessDeniedID,
		"access to a requested resource is denied",
		http.StatusForbidden,
		false)

	// ErrRequiredValueMissing is an error template for
	// ErrRequiredValueMissingID
	//
	//	Notes
	//		The error message has the form:
	//      "the required field '%s' is missing"
	//
	ErrRequiredValueMissing = NewErrorTemplate(
		ErrRequiredValueMissingID,
		"the required field '%s' is missing",
		http.StatusBadRequest,
		false)

	// ErrDataIntegrity is an error template for
	// ErrDataIntegrityID
	//
	//	Notes
	//		The error message has the form:
	//      "the value for '%s' is unexpected"
	//
	ErrDataIntegrity = NewErrorTemplate(
		ErrDataIntegrityID,
		"the value for '%s' is unexpected",
		http.StatusInternalServerError,
		false)

	// ErrResourceNotFound is an error template used when something
	// (generically referred to as a resource) cannot be found
	//
	//	Notes
	//		The error message has the form:
	//      "the resource '%s' was not found"
	//
	ErrResourceNotFound = NewErrorTemplate(
		ErrResourceNotFoundID,
		"the resource '%s' was not found",
		http.StatusNotFound,
		false)

	// ErrNotImplemented is an error template used when a feature is not
	// implemented
	//
	//	Notes
	//		The error message has the form:
	//      "'%s' is not implemented"
	//
	ErrNotImplemented = NewErrorTemplate(
		ErrNotImplementedID,
		"'%s' is not implemented",
		http.StatusNotImplemented,
		false)

	// ErrConversionFailed is an error template used when something cannot
	// be converted as expected
	//
	//	Notes
	//		The error message has the form:
	//      "the conversion '%s' could not be performed"
	//
	ErrConversionFailed = NewErrorTemplate(
		ErrConversionFailedID,
		"the conversion '%s' could not be performed",
		http.StatusBadRequest,
		false)

	// ErrDataValidationFailed is an error template used when something cannot
	// be converted as expected
	//
	//	Notes
	//		The error message has the form:
	//      "the validation of '%s' failed"
	//
	ErrDataValidationFailed = NewErrorTemplate(
		ErrDataValidationFailedID,
		"the validation of '%s' failed",
		http.StatusBadRequest,
		false)
)
