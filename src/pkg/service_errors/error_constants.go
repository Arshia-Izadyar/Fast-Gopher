package service_errors

const (
	OtpExists     = "otp exits"
	OtpUsed       = "otp used"
	OtpInvalid    = "otp invalid"
	ClaimNotFound = "claim not found"

	// user
	EmailExists        = "email already exits"
	UsernameExists     = "Username already exits"
	WrongPassword      = "Wrong Password!"
	PasswordsDontMatch = "Password don't match"
	BadRequest         = "Bad Request"

	TokenNotPresent    = "no token provided"
	TokenExpired       = "token is expired !"
	TokenInvalid       = "provided token is invalid"
	TokenInvalidFormat = "provided token has invalid format"
	NotRefreshToken    = "provided token is not a refresh token"
	InternalError      = "some thing happened"

	PermissionDenied = "Permission Denied"
)
