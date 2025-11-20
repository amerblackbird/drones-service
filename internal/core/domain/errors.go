package domain

import "fmt"

type DomainErrorCode string

const (
	DomainErrorCodeInvalidEmail       DomainErrorCode = "invalid_email"
	DomainErrorCodeInvalidPhoneNumber DomainErrorCode = "invalid_phone_number"
	DomainErrorCodeUserNotFound       DomainErrorCode = "user_not_found"
	DomainErrorCodeUserAlreadyExists  DomainErrorCode = "user_already_exists"
	DomainErrorCodeInvalidUserType    DomainErrorCode = "invalid_user_type"
	DomainErrorCodeDatabaseError      DomainErrorCode = "database_error"
	DomainErrorCodeUnauthorized       DomainErrorCode = "unauthorized"
	DomainErrorCodeInvalidInput       DomainErrorCode = "invalid_input"
	EmailAlreadyInUseError            DomainErrorCode = "email_already_in_use"
	LoginNotFoundError                DomainErrorCode = "login_not_found"

	// Http Status Codes
	DomainErrorCodeBadRequest          DomainErrorCode = "bad_request"
	DomainErrorCodeNotFound            DomainErrorCode = "not_found"
	DomainErrorCodeInternalServerError DomainErrorCode = "internal_server_error"
	DomainErrorCodeMethodNotAllowed    DomainErrorCode = "method_not_allowed"
	DomainErrorCodeConflict            DomainErrorCode = "conflict"
	DomainErrorCodeTooManyRequests     DomainErrorCode = "too_many_requests"
	DomainErrorCodeServiceUnavailable  DomainErrorCode = "service_unavailable"
	InvalidBodyError                   DomainErrorCode = "invalid_body_error"
	UserNotFoundError                  DomainErrorCode = "user_not_found_error"
	UserNotActiveError                 DomainErrorCode = "user_not_active_error"
	UnauthenticatedError               DomainErrorCode = "unauthenticated_error"
	UnauthorizedError                  DomainErrorCode = "unauthorized_error"

	// Auth
	InvalidOtpError              DomainErrorCode = "invalid_otp_error"
	InvalidTokenError            DomainErrorCode = "invalid_token_error"
	TokenExpiredError            DomainErrorCode = "token_expired_error"
	AccessDeniedError            DomainErrorCode = "access_denied_error"
	InvalidCredentialsError      DomainErrorCode = "invalid_credentials_error"
	PasswordTooWeakError         DomainErrorCode = "password_too_weak_error"
	PasswordMismatchError        DomainErrorCode = "password_mismatch_error"
	EmailAlreadyVerifiedError    DomainErrorCode = "email_already_verified_error"
	PhoneAlreadyVerifiedError    DomainErrorCode = "phone_already_verified_error"
	OtpSendFailedError           DomainErrorCode = "otp_send_failed_error"
	OtpExpiredError              DomainErrorCode = "otp_expired_error"
	OtpNotFoundError             DomainErrorCode = "otp_not_found_error"
	RefreshTokenNotFoundError    DomainErrorCode = "refresh_token_not_found_error"
	RefreshTokenExpiredError     DomainErrorCode = "refresh_token_expired_error"
	InvalidRefreshTokenError     DomainErrorCode = "invalid_refresh_token_error"
	AuthTokenInvalidError        DomainErrorCode = "auth_token_invalid_error"
	AuthTokenExpiredError        DomainErrorCode = "auth_token_expired_error"
	InsufficientPermissionsError DomainErrorCode = "insufficient_permissions_error"
	InvalidAuthTokenFormatError  DomainErrorCode = "invalid_auth_token_format_error"
	InvalidAuthTokenTypeError    DomainErrorCode = "invalid_auth_token_type_error"

	// Resource
	ResourceNotFoundError    DomainErrorCode = "resource_not_found_error"
	ResourceConflictError    DomainErrorCode = "resource_conflict_error"
	InvalidResourceError     DomainErrorCode = "invalid_resource_error"
	UnableToProcessError     DomainErrorCode = "unable_to_process_error"
	UnableToUpdateError      DomainErrorCode = "unable_to_update_error"
	UnableToDeleteError      DomainErrorCode = "unable_to_delete_error"
	UnableToCreateError      DomainErrorCode = "unable_to_create_error"
	ActivityLogNotFoundError DomainErrorCode = "activity_log_not_found_error"
	AuditLogNotFoundError    DomainErrorCode = "audit_log_not_found_error"

	// Form validation
	InvalidInputError DomainErrorCode = "invalid_input_error"

	// Params
	MissingParameterError DomainErrorCode = "missing_parameter_error"
	BodyIsRequiredError   DomainErrorCode = "body_is_required_error"
	ConflictError         DomainErrorCode = "conflict_error"

	// Orders
)

// Common error variables
var (
	ErrDroneNotFound = &DomainError{
		Code:    ResourceNotFoundError,
		Message: "Drone not found",
	}
	ErrOrderNotFound = &DomainError{
		Code:    ResourceNotFoundError,
		Message: "Order not found",
	}

	ErrWithdrawNotAllowed = &DomainError{
		Code:    UnableToProcessError,
		Message: "Only pending orders can be withdrawn",
	}
	ErrAlreadyWithdrawed = &DomainError{
		Code:    UnableToProcessError,
		Message: "Order has already been withdrawn",
	}
	ErrReserveNotAllowed = &DomainError{
		Code:    UnableToProcessError,
		Message: "Only pending orders can be reserved",
	}
	ErrAlreadyReserved = &DomainError{
		Code:    UnableToProcessError,
		Message: "Order has already been reserved",
	}
	ErrConfirmNotAllowed = &DomainError{
		Code:    UnableToProcessError,
		Message: "Only reserved orders can be confirmed",
	}
	ErrTransitNotAllowed = &DomainError{
		Code:    UnableToProcessError,
		Message: "Only picked up orders can be marked as in transit",
	}
	ErrHandoffNotAllowed = &DomainError{
		Code:    UnableToProcessError,
		Message: "Only failed orders can be handed off",
	}
	ErrArriveNotAllowed = &DomainError{
		Code:    UnableToProcessError,
		Message: "Only in transit orders can be marked as arrived",
	}
	ErrDeliverNotAllowed = &DomainError{
		Code:    UnableToProcessError,
		Message: "Only arrived orders can be marked as delivered",
	}
	ErrDeliverFailedNotAllowed = &DomainError{
		Code:    UnableToProcessError,
		Message: "Only in picked or transit or arrived orders can be marked as delivery failed",
	}
	ErrDroneFailedNotAllowed = &DomainError{
		Code:    UnableToProcessError,
		Message: "Only in transit or pickup orders can be marked as drone failed",
	}

	ErrReassignNotAllowed = &DomainError{
		Code:    UnableToProcessError,
		Message: "Only handoff orders can be reassigned",
	}
	ErrDroneMustBeIdle = &DomainError{
		Code:    UnableToProcessError,
		Message: "Drone must be in idle status to reserve an order",
	}
	ErrDroneIsLoading = &DomainError{
		Code:    UnableToProcessError,
		Message: "Drone is currently loading a package",
	}
	ErrDroneIsDelivering = &DomainError{
		Code:    UnableToProcessError,
		Message: "Drone is currently delivering a package",
	}
	ErrDroneIsReturning = &DomainError{
		Code:    UnableToProcessError,
		Message: "Drone is currently returning to base",
	}
	ErrDroneIsCharging = &DomainError{
		Code:    UnableToProcessError,
		Message: "Drone is currently charging",
	}
	ErrDroneInMaintenance = &DomainError{
		Code:    UnableToProcessError,
		Message: "Drone is under maintenance",
	}
	ErrDroneIsBroken = &DomainError{
		Code:    UnableToProcessError,
		Message: "Drone is broken",
	}
	ErrDroneUnderRepair = &DomainError{
		Code:    UnableToProcessError,
		Message: "Drone is under repair",
	}
	ErrIdleTransition = &DomainError{
		Code:    UnableToProcessError,
		Message: "idle can only transition to loading, charging, maintenanced, broken",
	}
	ErrLoadingTransition = &DomainError{
		Code:    UnableToProcessError,
		Message: "loading can only transition to delivering, broken",
	}
	ErrDeliveringTransition = &DomainError{
		Code:    UnableToProcessError,
		Message: "delivering can only transition to returning, broken",
	}
	ErrReturningTransition = &DomainError{
		Code:    UnableToProcessError,
		Message: "returning can only transition to idle, charging, broken",
	}
	ErrChargingTransition = &DomainError{
		Code:    UnableToProcessError,
		Message: "charging can only transition to idle, broken, returning",
	}
	ErrBrokenTransition = &DomainError{
		Code:    UnableToProcessError,
		Message: "broken can only transition to under_repair",
	}
	ErrUnderRepairTransition = &DomainError{
		Code:    UnableToProcessError,
		Message: "under_repair can only transition to maintenanced",
	}
	ErrMaintenancedTransition = &DomainError{
		Code:    UnableToProcessError,
		Message: "maintenanced can only transition to idle, returning",
	}
)

type DomainError struct {
	Code    DomainErrorCode `json:"code"`
	Message string          `json:"message"`
	Err     error           `json:"-"`
	Details *string         `json:"details,omitempty"`
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s → %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DomainErrorCode) Error() string {
	return string(*e)
}

func NewDomainError(code DomainErrorCode, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
