package core

const (
	// auth resps
	LoginSuccessCode = "login_success"
	UserCreatedCode  = "user_created"
	JWTRecievedCode  = "jwt_recieved"

	UsernameIsTaken      = "username_is_taken"
	WrongCredentialsCode = "wrong_credentials"
	InvalidJWTCode       = "invalid_jwt"

	// parsing resps
	InvalidIDCode          = "invalid_id"
	InvalidBodyCode        = "invalid_request_body"
	InvalidHeaderCode      = "invalid_header"
	InvalidQueryParamsCode = "invalid_query_params"

	// validation resps
	IDRequiredCode         = "id_is_required"
	BodyRequiredCode       = "request_body_is_required"
	HeaderRequiredCode     = "header_is_required"
	AuthHeaderRequiredCode = "auth_header_is_required"

	// stars resps
	StarReceivedCode  = "star_received"
	StarsReceivedCode = "stars_received"
	StarCreatedCode   = "star_created"
	StarUpdatedCode   = "star_updated"
	StarDeletedCode   = "star_deleted"

	StarNotFoundCode = "star_not_found"

	// movie resps
	MovieReceivedCode  = "movie_received"
	MoviesReceivedCode = "movies_received"
	MovieCreatedCode   = "movie_created"
	MovieUpdatedCode   = "movie_updated"
	MovieDeletedCode   = "movie_deleted"

	MovieNotFoundCode = "movie_not_found"

	// general
	UnauthorizedCode      = "general_unauthorized"
	AccessDeniedCode      = "general_access_denied"
	InternalErrorCode     = "general_internal"
	BadRequestErrorCode   = "general_bad_request_error"
	UnsupportedMethodCode = "general_unsupported_method"
	ForbiddenErrorCode    = "general_forbidden"
)
