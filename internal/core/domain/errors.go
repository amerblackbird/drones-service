package domain

import "fmt"

type UserError string

const (
	UserErrorInvalidEmail       UserError = "invalid_email"
	UserErrorInvalidPhoneNumber UserError = "invalid_phone_number"
	UserErrorUserNotFound       UserError = "user_not_found"
	UserErrorUserAlreadyExists  UserError = "user_already_exists"
	UserErrorInvalidUserType    UserError = "invalid_user_type"
	UserErrorDatabaseError      UserError = "database_error"
	UserErrorUnauthorized       UserError = "unauthorized"
	UserErrorInvalidInput       UserError = "invalid_input"
	EmailAlreadyInUseError      UserError = "email_already_in_use"
	LoginNotFoundError          UserError = "login_not_found"

	// Http Status Codes
	UserErrorBadRequest          UserError = "bad_request"
	UserErrorNotFound            UserError = "not_found"
	UserErrorInternalServerError UserError = "internal_server_error"
	UserErrorMethodNotAllowed    UserError = "method_not_allowed"
	UserErrorConflict            UserError = "conflict"
	UserErrorTooManyRequests     UserError = "too_many_requests"
	UserErrorServiceUnavailable  UserError = "service_unavailable"
	InvalidBodyError             UserError = "invalid_body_error"
	UserNotFoundError            UserError = "user_not_found_error"
	UserNotActiveError           UserError = "user_not_active_error"
	UnauthenticatedError         UserError = "unauthenticated_error"
	UnauthorizedError            UserError = "unauthorized_error"

	// Auth
	InvalidOtpError              UserError = "invalid_otp_error"
	InvalidTokenError            UserError = "invalid_token_error"
	TokenExpiredError            UserError = "token_expired_error"
	AccessDeniedError            UserError = "access_denied_error"
	InvalidCredentialsError      UserError = "invalid_credentials_error"
	PasswordTooWeakError         UserError = "password_too_weak_error"
	PasswordMismatchError        UserError = "password_mismatch_error"
	EmailAlreadyVerifiedError    UserError = "email_already_verified_error"
	PhoneAlreadyVerifiedError    UserError = "phone_already_verified_error"
	OtpSendFailedError           UserError = "otp_send_failed_error"
	OtpExpiredError              UserError = "otp_expired_error"
	OtpNotFoundError             UserError = "otp_not_found_error"
	RefreshTokenNotFoundError    UserError = "refresh_token_not_found_error"
	RefreshTokenExpiredError     UserError = "refresh_token_expired_error"
	InvalidRefreshTokenError     UserError = "invalid_refresh_token_error"
	AuthTokenInvalidError        UserError = "auth_token_invalid_error"
	AuthTokenExpiredError        UserError = "auth_token_expired_error"
	InsufficientPermissionsError UserError = "insufficient_permissions_error"
	InvalidAuthTokenFormatError  UserError = "invalid_auth_token_format_error"
	InvalidAuthTokenTypeError    UserError = "invalid_auth_token_type_error"

	// Resource
	ResourceNotFoundError    UserError = "resource_not_found_error"
	ResourceConflictError    UserError = "resource_conflict_error"
	InvalidResourceError     UserError = "invalid_resource_error"
	UnableToProcessError     UserError = "unable_to_process_error"
	UnableToUpdateError      UserError = "unable_to_update_error"
	UnableToDeleteError      UserError = "unable_to_delete_error"
	UnableToCreateError      UserError = "unable_to_create_error"
	ActivityLogNotFoundError UserError = "activity_log_not_found_error"
	AuditLogNotFoundError    UserError = "audit_log_not_found_error"

	// Form validation
	InvalidInputError UserError = "invalid_input_error"

	// Params
	MissingParameterError UserError = "missing_parameter_error"
	BodyIsRequiredError   UserError = "body_is_required_error"
	ConflictError         UserError = "conflict_error"

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
	Code    UserError `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
	Details *string   `json:"details,omitempty"`
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s → %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *UserError) Error() string {
	return string(*e)
}

func NewDomainError(code UserError, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
