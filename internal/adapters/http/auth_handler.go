package http

import (
	domain "drones/internal/core/domain"
	ports "drones/internal/ports"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	service   ports.AuthService
	validator *validator.Validate
	logger    ports.Logger
}

func NewAuthHandler(service ports.AuthService, validator *validator.Validate, logger ports.Logger) *AuthHandler {
	return &AuthHandler{
		service:   service,
		validator: validator,
		logger:    logger,
	}
}

// RegisterRoutes registers all bank routes
func (h *AuthHandler) RegisterRoutes(r *mux.Router) {

	// Special routes
	r.HandleFunc("/token", h.HandleAuthorize).Methods("POST")
}

func (h *AuthHandler) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	var body domain.AuthRequest

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		ResponseWithError(w, domain.NewDomainError(domain.InvalidInputError, "Invalid request body", err))
		return
	}

	// Validate the request body
	if err := h.validator.Struct(body); err != nil {
		ResponseWithValidationError(w, http.StatusBadRequest, domain.GetValidationErrors(err.(validator.ValidationErrors)))
		return
	}

	response, err := h.service.Login(r.Context(), body.Name, body.Type)
	if err != nil {
		ResponseWithError(w, err)
		return
	}
	if response == nil {
		ResponseWithError(w, &domain.DomainError{
			Code:    domain.UserNotFoundError,
			Message: "User not found",
		})
		return
	}
	ResponseWithJSON(w, http.StatusOK, response.ToDTO())
}

// func (h *HTTPHandler) HandleAuthorizeVerification(w http.ResponseWriter, r *http.Request) {
// 	meta := domain.ExtractRequestInfo(r)

// 	// Extract token from URL
// 	token := mux.Vars(r)["token"]
// 	if token == "" {
// 		h.responseWithError(w, http.StatusBadRequest, &domain.DomainError{
// 			Code:    domain.MissingParameterError,
// 			Message: "Token is required",
// 		})
// 		return
// 	}

// 	var body ports.LoginWithPhoneNumberVerificationRequest
// 	// Decode the request body
// 	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
// 		// Log the error
// 		h.logError(err, "Verification failed - invalid request body", r)

// 		// Return a 400 Bad Request error
// 		h.writeError(w, http.StatusBadRequest, string(domain.InvalidBodyError))
// 		return
// 	}

// 	// Find opts with phoneNumber and token
// 	// Validate the request body
// 	if err := h.Validator.Struct(body); err != nil {
// 		// Log the validation error
// 		h.logError(err, "Verification failed - invalid request body", r)

// 		// Return a 400 Bad Request error with validation errors
// 		errors := domain.GetValidationErrors(err.(validator.ValidationErrors))

// 		// Log the validation errors
// 		h.writeValidationError(w, http.StatusBadRequest, errors)
// 		return
// 	}

// 	response, err := h.handleAuthenticateWithPhoneNumberVerification(r.Context(), token, &body, meta)
// 	if err != nil {
// 		h.responseWithError(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	if response == nil {
// 		h.responseWithError(w, http.StatusBadRequest, &domain.DomainError{
// 			Code:    domain.UserNotFoundError,
// 			Message: "User not found",
// 		})
// 		return
// 	}
// 	h.writeJSON(w, http.StatusOK, response)

// }

// func (h *HTTPHandler) handleAuthenticateWithPhoneNumber(ctx context.Context, userType domain.UserType, login *ports.LoginWithPhoneNumberRequest, meta domain.RequestInfo) (*ports.LoginWithPhoneNumberResponse, error) {
// 	// Implement the logic to handle authentication with phone number
// 	// Check if the user exists, create a new user if not, and return the response
// 	var user *domain.User
// 	h.logger.Info("Authenticating user with phone number", login.PhoneNumber, userType)
// 	user, err := h.usersService.GetByPhone(ctx, login.PhoneNumber, userType)

// 	if err != nil {
// 		// Log err
// 		h.logger.Error("Failed to get user by phone number", zap.Error(err))
// 	} else {
// 		h.logger.Info("Found user by phone number", zap.Any("user", user))
// 	}

// 	// If user exists, create a new user
// 	if user == nil {
// 		h.logger.Info("Creating new user with phone number", zap.Any("login", login))
// 		locale := login.Locale
// 		if login.Locale != nil {
// 			locale = utils.StringPtr("en_US")
// 		}
// 		password := utils.RandomPassword()
// 		salt := utils.RandomSalt()
// 		hashedPassword := utils.HashText(password, salt)
// 		userDto := &domain.CreateUser{
// 			Phone:             &login.PhoneNumber,
// 			Name:              "User",
// 			Username:          login.PhoneNumber,
// 			UserType:          userType,
// 			DeviceID:          login.DeviceID,
// 			HashedPassword:    hashedPassword,
// 			Salt:              salt,
// 			NotificationToken: login.NotificationToken,
// 			Locale:            locale,
// 			UserAgent:         meta.UserAgent,
// 			IPAddress:         meta.IpAddress,
// 		}
// 		var err error
// 		user, err = h.usersService.CreateUser(ctx, login.PhoneNumber, userDto)
// 		if err != nil {
// 			h.logger.Error("Failed to create user", zap.Error(err))
// 			// Return error
// 			return nil, err
// 		}
// 		if user == nil {
// 			// Return error
// 			return nil, &domain.DomainError{
// 				Code:    domain.UserErrorInternalServerError,
// 				Message: "Failed to create user",
// 			}
// 		}

// 		// Log activity
// 		h.eventPublisher.PublishUserCreated(ctx, events.UserCreatedEvent{
// 			UserID:            user.ID,
// 			UserType:          user.Type,
// 			Name:              user.Name,
// 			Phone:             user.Phone,
// 			Locale:            user.Locale,
// 			Country:           user.Country,
// 			Email:             user.Email,
// 			DeviceID:          login.DeviceID,
// 			NotificationToken: login.NotificationToken,
// 		})

// 	} else {
// 		// Check if user has one-time password (OTP) enabled
// 		if !user.Active {

// 			// Log activity
// 			h.eventPublisher.LogActivity(ctx, user.ID, "login_attempt_inactive_user", meta.ActivityLogMeta())

// 			// Return error
// 			return nil, domain.NewDomainError(domain.UserNotActiveError, "User is not active", nil)
// 		}
// 	}

// 	// check if user has more than 3 failed otp attempts in the last 15 minutes
// 	otps, err := h.otpsService.GetUserOtpLast15Min(ctx, user.ID, 3)
// 	if err != nil {
// 		// Log activity
// 		h.eventPublisher.LogActivity(ctx, user.ID, "login_attempt_failed_to_fetch_otps", meta.ActivityLogMeta())

// 		// Return error
// 		return nil, domain.NewDomainError(domain.UserErrorInternalServerError, "Internal server error", err)
// 	}

// 	if len(otps) >= 3 {
// 		// Suspend user if more than 3 failed attempts
// 		active := false
// 		updateDto := domain.UpdateUserInfo{
// 			Active: &active,
// 		}
// 		if _, err := h.usersService.UpdateUser(ctx, user.ID, &updateDto); err != nil {
// 			// Log activity
// 			h.logger.Error("Failed to suspend user", zap.Error(err))
// 		}
// 		// Log activity
// 		h.eventPublisher.LogActivity(ctx, user.ID, "login_attempt_failed_to_suspend_user", meta.ActivityLogMeta())

// 		// Return error
// 		return nil, domain.NewDomainError(domain.UserErrorTooManyRequests, "Too many failed attempts. User is suspended.", nil)
// 	}

// 	// Suspend user if more than 3 failed attempts
// 	otpCode, hashedOtp, salt, token := utils.GenerateVerfiyCred()

// 	h.logger.Info("Generated OTP", zap.String("otp", otpCode), zap.String("hashedOtp", hashedOtp), zap.String("salt", salt), zap.String("token", token))

// 	// Create otp
// 	otp := &domain.CreateOTP{
// 		HashCode:  hashedOtp,
// 		Salt:      salt,
// 		UserID:    user.ID,
// 		Token:     token,
// 		Username:  *user.Phone,
// 		ExpiresAt: time.Now().Add(2 * time.Minute),
// 	}

// 	_, err = h.otpsService.Create(ctx, otp)
// 	if err != nil {
// 		// Log activity
// 		h.eventPublisher.LogActivity(ctx, user.ID, "login_attempt_failed_to_create_otp", meta.ActivityLogMeta())

// 		// Return error
// 		return nil, err
// 	}

// 	// Log activity
// 	h.eventPublisher.LogActivity(ctx, user.ID, "login_attempt_created_otp", meta.ActivityLogMeta())

// 	// Send OTP code
// 	h.eventPublisher.SendOTP(ctx, user.ID, *user.Phone, otpCode)

// 	// Map the user to the response format
// 	response := &ports.LoginWithPhoneNumberResponse{
// 		UserID: user.ID,
// 		Token:  token,
// 	}

// 	return response, nil
// }

// func (h *HTTPHandler) handleAuthenticateWithPhoneNumberVerification(ctx context.Context, token string, verification *ports.LoginWithPhoneNumberVerificationRequest, info domain.RequestInfo) (*ports.LoginWithPhoneNumberVerificationResponse, error) {
// 	// Find otp with phoneNumber and token
// 	otp, err := h.otpsService.FindByTokenAndUsername(ctx, token, verification.PhoneNumber)
// 	if err != nil {
// 		// Activity log
// 		h.eventPublisher.LogActivity(ctx, "", "login_attempt_failed_to_find_otp", info.ActivityLogMeta())
// 		h.logger.Error("Failed to find OTP by token and username", zap.Error(err))
// 		return nil, err
// 	}
// 	if otp == nil {
// 		// Activity log
// 		h.eventPublisher.LogActivity(ctx, "", "login_attempt_invalid_token_or_phone", info.ActivityLogMeta())

// 		return nil, &domain.DomainError{
// 			Code:    domain.InvalidOtpError,
// 			Message: "Invalid token or phone number",
// 		}
// 	}

// 	if otp.Used {
// 		// Activity log
// 		h.eventPublisher.LogActivity(ctx, otp.UserID, "login_attempt_otp_already_used", info.ActivityLogMeta())

// 		return nil, &domain.DomainError{
// 			Code:    domain.InvalidOtpError,
// 			Message: "OTP already used",
// 		}
// 	}

// 	// Check if otp is expired
// 	if time.Now().After(otp.ExpiresAt) {
// 		// Activity log
// 		h.eventPublisher.LogActivity(ctx, otp.UserID, "login_attempt_expired_otp", info.ActivityLogMeta())

// 		return nil, &domain.DomainError{
// 			Code:    domain.InvalidOtpError,
// 			Message: "OTP has expired",
// 		}
// 	}

// 	// Verify otp code
// 	if !utils.VerifyOtpHash(verification.OtpCode, otp.Salt, otp.HashCode) {
// 		// Activity log
// 		h.eventPublisher.LogActivity(ctx, otp.UserID, "login_attempt_invalid_otp_code", info.ActivityLogMeta())

// 		return nil, &domain.DomainError{
// 			Code:    domain.InvalidOtpError,
// 			Message: "Invalid OTP code",
// 		}
// 	}

// 	// Get user
// 	user, err := h.usersService.GetUserByID(ctx, otp.UserID)
// 	if err != nil {
// 		// Activity log
// 		h.eventPublisher.LogActivity(ctx, otp.UserID, "login_attempt_failed_to_get_user", info.ActivityLogMeta())
// 		h.logger.Error("Failed to get user by ID", zap.Error(err))
// 		return nil, err
// 	}
// 	if user == nil {
// 		// Activity log
// 		h.eventPublisher.LogActivity(ctx, otp.UserID, "login_attempt_user_not_found", info.ActivityLogMeta())
// 		return nil, &domain.DomainError{
// 			Code:    domain.UserNotFoundError,
// 			Message: "User not found",
// 		}
// 	}

// 	// Update otp set as used
// 	if err := h.otpsService.SetUsed(ctx, otp.ID); err != nil {
// 		// Activity log
// 		h.eventPublisher.LogActivity(ctx, otp.UserID, "login_attempt_failed_to_set_otp_as_used", info.ActivityLogMeta())
// 		h.logger.Error("Failed to set OTP as used", zap.Error(err))
// 		return nil, err
// 	}

// 	// Generate JWT token
// 	accessToken, refreshToken, err := h.authService.Login(user.ID, string(user.Type), []string{})
// 	if err != nil {
// 		h.logger.Error("Failed to generate JWT token", zap.Error(err))
// 		return nil, err
// 	}

// 	// Emit login event
// 	h.publishUserLoggedInEvent(ctx, user, info)

// 	// Store refresh token in the login table
// 	if err := h.usersService.StoreRefreshToken(ctx, user.ID, refreshToken, info.DeviceID, info.IpAddress); err != nil {
// 		h.logger.Error("Failed to store refresh token", zap.Error(err))
// 	}

// 	return &ports.LoginWithPhoneNumberVerificationResponse{
// 		UserID:       user.ID,
// 		Name:         user.Name,
// 		Type:         string(user.Type),
// 		Phone:        user.Phone,
// 		Email:        user.Email,
// 		AccessToken:  accessToken,
// 		RefreshToken: refreshToken,
// 	}, nil
// }

// func (h *HTTPHandler) publishUserLoggedInEvent(ctx context.Context, user *domain.User, info domain.RequestInfo) {
// 	location := info.GetLocation()
// 	ipAddress := info.GetIpAddress()
// 	deviceId := info.GetDeviceId()

// 	h.eventPublisher.PublishUserLoggedIn(ctx, events.UserLoggedInEvent{
// 		UserID:    user.ID,
// 		DeviceID:  &deviceId,
// 		IPAddress: &ipAddress,
// 		Location:  &location,
// 	})
// }

// func (h *HTTPHandler) HandleAuthorizeVerificationResend(w http.ResponseWriter, r *http.Request) {
// 	info := domain.ExtractRequestInfo(r)

// 	// Extract token from URL
// 	token := mux.Vars(r)["token"]

// 	if token == "" {
// 		h.responseWithError(w, http.StatusBadRequest, &domain.DomainError{
// 			Code:    domain.MissingParameterError,
// 			Message: "Token is required",
// 		})
// 		return
// 	}

// 	response, err := h.HandleResendVerificationCode(r.Context(), token, info)
// 	if err != nil {
// 		h.responseWithError(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	if response == nil {
// 		h.responseWithError(w, http.StatusBadRequest, &domain.DomainError{
// 			Code:    domain.UserNotFoundError,
// 			Message: "User not found",
// 		})
// 		return
// 	}
// 	h.writeJSON(w, http.StatusOK, response)

// }

// func (h *HTTPHandler) HandleResendVerificationCode(ctx context.Context, token string, info domain.RequestInfo) (*ports.LoginWithPhoneNumberResendResponse, error) {
// 	// Get otp by token
// 	otp, err := h.otpsService.GetByToken(ctx, token)
// 	if err != nil {
// 		// Activity log
// 		h.eventPublisher.LogActivity(ctx, "", "resend_otp_failed_to_find_otp", info.ActivityLogMeta())
// 		h.logger.Error("Failed to find OTP by token and username", zap.Error(err))
// 		return nil, err
// 	}
// 	if otp == nil {
// 		// Activity log
// 		h.eventPublisher.LogActivity(ctx, "", "resend_otp_invalid_token_or_phone", info.ActivityLogMeta())

// 		return nil, &domain.DomainError{
// 			Code:    domain.InvalidOtpError,
// 			Message: "Invalid token or phone number",
// 		}
// 	}

// 	if otp.Used {
// 		// Activity log
// 		h.eventPublisher.LogActivity(ctx, otp.UserID, "resend_otp_already_used", info.ActivityLogMeta())

// 		return nil, &domain.DomainError{
// 			Code:    domain.InvalidOtpError,
// 			Message: "OTP already used",
// 		}
// 	}

// 	// Check if otp is not expired
// 	if time.Now().Before(otp.ExpiresAt) {
// 		// Activity log
// 		h.eventPublisher.LogActivity(ctx, otp.UserID, "resend_otp_not_expired", info.ActivityLogMeta())

// 		return nil, &domain.DomainError{
// 			Code:    domain.InvalidOtpError,
// 			Message: "OTP is not expired yet",
// 		}
// 	}

// 	// Generate new otp code
// 	otpCode, hashedOtp, salt, newToken := utils.GenerateVerfiyCred()

// 	h.logger.Info("Generated new OTP for resend", zap.String("otp", otpCode), zap.String("hashedOtp", hashedOtp), zap.String("salt", salt), zap.String("token", newToken))

// 	// Update otp
// 	updateOtp := &domain.UpdateOTP{
// 		Token:     newToken,
// 		HashCode:  hashedOtp,
// 		Salt:      salt,
// 		ExpiresAt: time.Now().Add(2 * time.Minute).Format(time.RFC3339),
// 		Expiry:    120, // 2 minutes in seconds
// 	}

// 	if _, err := h.otpsService.Update(ctx, otp.ID, updateOtp); err != nil {
// 		// Activity log
// 		h.eventPublisher.LogActivity(ctx, otp.UserID, "resend_otp_failed_to_update_otp", info.ActivityLogMeta())
// 		h.logger.Error("Failed to update OTP", zap.Error(err))
// 		return nil, err
// 	}

// 	// Log activity
// 	h.eventPublisher.LogActivity(ctx, otp.UserID, "resend_otp_updated_otp", info.ActivityLogMeta())

// 	// Send OTP code
// 	h.eventPublisher.SendOTP(ctx, otp.UserID, otp.Username, otpCode)

// 	// Log activity
// 	h.eventPublisher.LogActivity(ctx, otp.UserID, "resend_otp_sent_otp", info.ActivityLogMeta())

// 	// Return response
// 	return &ports.LoginWithPhoneNumberResendResponse{
// 		Token:  newToken,
// 		UserID: otp.UserID,
// 	}, nil
// }

// func (h *HTTPHandler) HandleGetUserHandler(w http.ResponseWriter, r *http.Request) {
// 	// Get user ID from URL
// 	id := mux.Vars(r)["id"]
// 	if id == "" {
// 		h.writeError(w, http.StatusBadRequest, "User ID is required")
// 		return
// 	}
// 	// Get user
// 	user, err := h.usersService.GetUserByID(r.Context(), id)
// 	if err != nil {
// 		h.responseWithError(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	if user == nil {
// 		h.responseWithError(w, http.StatusBadRequest, &domain.DomainError{
// 			Code:    domain.UserNotFoundError,
// 			Message: "User not found",
// 		})
// 	}
// 	h.writeJSON(w, http.StatusOK, user)
// }

// func (h *HTTPHandler) HandleGetUserPreference(w http.ResponseWriter, r *http.Request) {
// 	user, ok := UserFromContext(r.Context())
// 	if !ok || user == nil {
// 		h.responseWithError(w, http.StatusUnauthorized, &domain.DomainError{
// 			Code:    domain.UserNotFoundError,
// 			Message: "User not found in context",
// 		})
// 		return
// 	}
// 	h.logger.Info("Fetching user preference", zap.String("userID", user.ID), zap.String("userName", user.Name),
// 		zap.String("userType", string(user.Type)))

// 	log.Println("User fetched from context:", user.AvatarUrl)

// 	result := ports.PreferenceResponse{
// 		UserID:    user.ID,
// 		Name:      user.Name,
// 		Type:      string(user.Type),
// 		Phone:     user.Phone,
// 		Email:     user.Email,
// 		Locale:    user.Locale,
// 		Country:   user.Country,
// 		AvatarUrl: user.AvatarUrl,
// 		Active:    user.Active,
// 		Bio:       user.Bio,
// 	}

// 	h.writeJSON(w, http.StatusOK, result)

// }

// func (h *HTTPHandler) HandleGetUser(ctx context.Context, userID string) (*domain.User, error) {
// 	// Get user
// 	user, err := h.usersService.GetUserByID(ctx, userID)
// 	if err != nil {
// 		h.logger.Error("Failed to get user by ID", zap.Error(err))
// 		return nil, err
// 	}
// 	if user == nil {
// 		return nil, &domain.DomainError{
// 			Code:    domain.UserNotFoundError,
// 			Message: "User not found",
// 		}
// 	}
// 	return user, nil
// }

// func (h *HTTPHandler) HandleUpdateUserPreference(w http.ResponseWriter, r *http.Request) {
// 	info := domain.ExtractRequestInfo(r)

// 	user, ok := UserFromContext(r.Context())
// 	if !ok || user == nil {
// 		h.responseWithError(w, http.StatusUnauthorized, &domain.DomainError{
// 			Code:    domain.UserNotFoundError,
// 			Message: "User not found in context",
// 		})
// 		return
// 	}

// 	var body domain.UpdatePreferenceDto
// 	// Decode the request body
// 	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
// 		// Log the error
// 		h.logError(err, "Update preference failed - invalid request body", r)

// 		// Activity log
// 		h.eventPublisher.LogActivity(r.Context(), user.ID, "preference_update_failed_invalid_body", info.ActivityLogMeta())

// 		// Return a 400 Bad Request error
// 		h.writeError(w, http.StatusBadRequest, string(domain.InvalidBodyError))
// 		return
// 	}

// 	// Validate the request body
// 	if err := h.Validator.Struct(body); err != nil {
// 		// Log the validation error
// 		h.logError(err, "Update preference failed - invalid request body", r)

// 		// Activity log
// 		h.eventPublisher.LogActivity(r.Context(), user.ID, "preference_update_failed_invalid_body", info.ActivityLogMeta())

// 		// Return a 400 Bad Request error with validation errors
// 		errors := domain.GetValidationErrors(err.(validator.ValidationErrors))

// 		// Log the validation errors
// 		h.writeValidationError(w, http.StatusBadRequest, errors)
// 		return
// 	}

// 	dto := domain.UpdateUserInfo{
// 		Name:    body.Name,
// 		Locale:  body.Locale,
// 		Country: body.Country,
// 		Email:   body.Email,
// 	}

// 	log.Println("Update preference request body:", body)

// 	updatedUser, err := h.usersService.UpdateUser(r.Context(), user.ID, &dto)
// 	if err != nil {
// 		// Log activity
// 		h.eventPublisher.LogActivity(r.Context(), user.ID, "preference_updated_failed", info.ActivityLogMeta())

// 		// Respond with error
// 		h.responseWithError(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	if updatedUser == nil {
// 		// Log activity
// 		h.eventPublisher.LogActivity(r.Context(), user.ID, "preference_updated_failed", info.ActivityLogMeta())

// 		// Respond with error
// 		h.responseWithError(w, http.StatusBadRequest, &domain.DomainError{
// 			Code:    domain.UserNotFoundError,
// 			Message: "User not found",
// 		})
// 		return
// 	}

// 	result := ports.PreferenceResponse{
// 		UserID:       updatedUser.ID,
// 		Name:         updatedUser.Name,
// 		Type:         string(updatedUser.Type),
// 		Phone:        updatedUser.Phone,
// 		Email:        updatedUser.Email,
// 		Locale:       updatedUser.Locale,
// 		Country:      updatedUser.Country,
// 		DeviceID:     updatedUser.DeviceID,
// 		Notification: updatedUser.NotificationToken,
// 	}

// 	/// Emit user updated event
// 	h.eventPublisher.PublishUserUpdated(r.Context(), events.UserUpdatedEvent{
// 		UserID:  updatedUser.ID,
// 		Name:    &updatedUser.Name,
// 		Email:   updatedUser.Email,
// 		Locale:  updatedUser.Locale,
// 		Country: updatedUser.Country,
// 	})

// 	// Activity log
// 	h.eventPublisher.LogActivity(r.Context(), updatedUser.ID, "preference_updated", info.ActivityLogMeta())

// 	h.writeJSON(w, http.StatusOK, result)
// }

// func (h *HTTPHandler) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
// 	user, ok := UserFromContext(r.Context())
// 	if !ok || user == nil {
// 		h.responseWithError(w, http.StatusUnauthorized, &domain.DomainError{
// 			Code:    domain.UserNotFoundError,
// 			Message: "User not found in context",
// 		})
// 		return
// 	}

// 	// Implement the logic to handle token refresh
// 	// Validate the refresh token, generate a new access token, and return the response
// 	accessToken, _, err := h.authService.Login(user.ID, string(user.Type), []string{})
// 	if err != nil {
// 		h.responseWithError(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	if accessToken == "" {
// 		h.responseWithError(w, http.StatusBadRequest, &domain.DomainError{
// 			Code:    domain.UserNotFoundError,
// 			Message: "User not found",
// 		})
// 		return
// 	}

// 	result := &ports.RefreshTokenResponse{
// 		AccessToken: accessToken,
// 	}

// 	// Log activity
// 	info := domain.ExtractRequestInfo(r)
// 	h.eventPublisher.LogActivity(r.Context(), user.ID, "token_refreshed", info.ActivityLogMeta())

// 	h.writeJSON(w, http.StatusOK, result)
// }
