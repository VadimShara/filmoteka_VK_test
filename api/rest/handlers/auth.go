package handlers

import (
	"context"
	"errors"
	"net/http"

	_ "vk-test-task/api/rest/handlers/model" // to register swagger models
	"vk-test-task/api/rest/presenters/user"
	"vk-test-task/internal/core"
	"vk-test-task/internal/service/auth"
	"vk-test-task/pkg/hash"
	"vk-test-task/pkg/web"
	"vk-test-task/pkg/webutil"

	"github.com/jackc/pgx/v5"
)

// @Title Login
// @Resource Auth
// @Description Login user to get JWT token
// @Param user body model.LoginRequest true "Login Request"
// @Success 200 object model.LoginResponse "Successful login"
// @Failure 400 object model.BadRequestInvalidBodyResponse "Bad request error"
// @Failure 401 object model.WrongCredentialsResponse "Unauthorized error"
// @Failure 422 object model.ValidationResponse "Validation error"
// @Failure 500 object model.InternalResponse "Internal server error"
// @Route /api/v1/auth/login [post]
func (r *Resolver) login(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		webutil.SendJSONResponse(w, http.StatusMethodNotAllowed, web.ErrorResponse(core.UnsupportedMethodCode, nil, nil))
	}

	var model auth.LoginModel

	if !webutil.BodyCheck(w, req, &model) {
		return
	}

	reqPassHash := hash.CalculateHash(model.Password)

	passHash, role, err := r.authService.GetPassHashAndRoleByUsername(context.Background(), model.Username)
	if err != nil || passHash != reqPassHash {
		switch {
		case err == pgx.ErrNoRows || passHash != reqPassHash:
			webutil.SendJSONResponse(w, http.StatusUnauthorized, web.ErrorResponse(core.WrongCredentialsCode, nil, nil))
			return
		default:
			webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return
		}
	}

	data, err := r.authService.CreateToken(context.Background(), model.Username, role)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrUsernameExists):
			webutil.SendJSONResponse(w, http.StatusConflict, web.ErrorResponse(core.UsernameIsTaken, nil, nil))
			return
		default:
			webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return
		}
	}

	webutil.SendJSONResponse(w, http.StatusOK, web.OKResponse(core.JWTRecievedCode, data, nil))
}

// @Title Sign-up
// @Resource Auth
// @Description Sign-up user
// @Param user body model.SignUpRequest true "Sign-up Request"
// @Success 200 object model.SignUpResponse "Successful sign-up"
// @Failure 400 object model.BadRequestInvalidBodyResponse "Bad request error"
// @Failure 409 object model.ConflictUsernameResponse "Username is taken error"
// @Failure 422 object model.ValidationResponse "Validation error"
// @Failure 500 object model.InternalResponse "Internal server error"
// @Route /api/v1/auth/signup [post]
func (r *Resolver) signup(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		webutil.SendJSONResponse(w, http.StatusMethodNotAllowed, web.ErrorResponse(core.UnsupportedMethodCode, nil, nil))
	}

	var model auth.SignUpModel

	if !webutil.BodyCheck(w, req, &model) {
		return
	}

	data, err := r.authService.SignUp(context.Background(), model)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrUsernameExists):
			webutil.SendJSONResponse(w, http.StatusConflict, web.ErrorResponse(core.UsernameIsTaken, nil, nil))
			return
		default:
			webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return
		}
	}

	pres := user.PresentUser(data)

	webutil.SendJSONResponse(w, http.StatusCreated, pres.Response(core.UserCreatedCode))
}
