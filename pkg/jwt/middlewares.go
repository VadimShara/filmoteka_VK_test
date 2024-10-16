package jwt

import (
	"net/http"

	"vk-test-task/internal/core"
	"vk-test-task/pkg/web"
	"vk-test-task/pkg/webutil"
)

func (s *Service) Verify(w http.ResponseWriter, r *http.Request) (*UserData, bool) {
	tokenString := webutil.AuthHeaderChecker(w, r)
	if tokenString == "" {
		return nil, false
	}

	user, err := s.ValidateToken(tokenString)
	if err != nil {
		switch err {
		case ErrInvalidToken:
			webutil.SendJSONResponse(w, http.StatusUnauthorized, web.ErrorResponse(core.UnauthorizedCode, nil, nil))
			return nil, false
		default:
			webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return nil, false
		}
	}

	return user, true
}
