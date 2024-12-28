package errors

import (
	"fmt"
)

// A custom error for api responses, contains error code and message
type ApiError struct {
	Code          int
	Message       string
	RequiresLogin bool
}

// Generic ApiErrors
var (
	// Token Errors - Requires Re-Login
	ApiErrorToken        = ApiError{Code: 9811, Message: "DFE: Invalid Token", RequiresLogin: true}
	ApiErrorTokenExpired = ApiError{Code: 9812, Message: "DFE: Token Expired", RequiresLogin: true}

	// Common Errors
	ApiErrorSystem       = ApiError{Code: 9001, Message: "DFE: System error"}
	ApiErrorParamInvalid = ApiError{Code: 9002, Message: "DFE: Input param invalid"}
	ApiErrorParamMissing = ApiError{Code: 9003, Message: "DFE: Input param missing"}
	ApiError404          = ApiError{Code: 9004, Message: "404 Not Found"}

	// Authorization and Authentication Errors
	ApiErrorCredentials      = ApiError{Code: 9801, Message: "DFE: Invalid Credentials"}
	ApiErrorPermissionDenied = ApiError{Code: 9802, Message: "DFE: Permission Denied"}
	ApiErrorRequiresAdmin    = ApiError{Code: 9803, Message: "DFE: Permission Denied, Admin Credentials Required"}

	// Request/Response Errors
	ApiErrorDuplicateMobile = ApiError{Code: 9101, Message: "DFE: Mobile number not available"}

	// Communicator Errors
	ApiErrorCommunicator = ApiError{Code: 9905, Message: "DFE: Communicator general error"}
)

// Implement Error() string function so that ApiError implements error interface
func (a ApiError) Error() string {
	return fmt.Sprintf("ApiError: (%d, %s)", a.Code, a.Message)
}

// Convert any err to an ApiError
func ErrorToApiError(err error) ApiError {
	return ApiError{Code: 9999, Message: err.Error()}
}
