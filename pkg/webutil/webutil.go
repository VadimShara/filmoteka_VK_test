package webutil

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"vk-test-task/internal/core"
	"vk-test-task/pkg/logger"
	"vk-test-task/pkg/web"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func SendJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response) //nolint
}

func BodyCheck(w http.ResponseWriter, r *http.Request, entity interface{}) bool {
	if r.Body == nil {
		SendJSONResponse(w, http.StatusBadRequest, web.ErrorResponse(core.BodyRequiredCode, nil, nil))
		return false
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&entity); err != nil {
		SendJSONResponse(w, http.StatusBadRequest, web.ErrorResponse(core.InvalidBodyCode, nil, nil))
		return false
	}

	err := validate.RegisterValidation("date", (&CustomValidator{validate}).ValidateDate)
	if err != nil {
		logger.Log.Error("init validator", "error", err.Error())
		SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
		return false
	}

	if err := validate.Struct(entity); err != nil {
		var verrors validator.ValidationErrors
		ok := errors.As(err, &verrors)
		if !ok {
			logger.Log.Debug("error validation", "error", err.Error())
			SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return false
		}
		SendJSONResponse(w, http.StatusUnprocessableEntity, web.ValidationErrorResponse(web.ValidationErrors(verrors), nil))
		return false
	}

	return true
}

func ParseID(w http.ResponseWriter, r *http.Request) int {
	parts := strings.Split(r.URL.Path, "/")
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, web.ErrorResponse(core.InvalidIDCode, nil, nil))
		return 0
	}

	return id
}

func QueryParser(r *http.Request) map[string]string {
	queries := r.URL.Query()

	m := make(map[string]string, len(queries))
	for key, value := range queries {
		m[key] = value[0]
	}

	return m
}

func QueryValidator(w http.ResponseWriter, r *http.Request, entity interface{}) bool {
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("sort_params", (&CustomValidator{validate}).ValidateSortParams)
	if err != nil {
		logger.Log.Error("init validator", "error", err.Error())
		SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
		return false
	}

	if err := validate.Struct(entity); err != nil {
		var verrors validator.ValidationErrors
		ok := errors.As(err, &verrors)
		if !ok {
			logger.Log.Debug("error validation", "error", err.Error())
			SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return false
		}
		SendJSONResponse(w, http.StatusUnprocessableEntity, web.ValidationErrorResponse(web.ValidationErrors(verrors), nil))
		return false
	}

	return true
}

func AuthHeaderChecker(w http.ResponseWriter, r *http.Request) string {
	headerValue := r.Header.Get("Authorization")
	if headerValue == "" {
		SendJSONResponse(w, http.StatusUnauthorized, web.ErrorResponse(core.AuthHeaderRequiredCode, nil, nil))
		return ""
	}

	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("jwt_auth_header", (&CustomValidator{validate}).ValidateAuthHeader)
	if err != nil {
		logger.Log.Error("init validator", "error", err.Error())
		SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
		return ""
	}

	if err := validate.Var(headerValue, "jwt_auth_header"); err != nil {
		var verrors validator.ValidationErrors
		ok := errors.As(err, &verrors)
		if !ok {
			logger.Log.Debug("error validation", "error", err.Error())
			SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return ""
		}
		SendJSONResponse(w, http.StatusUnprocessableEntity, web.ValidationErrorResponse(web.ValidationErrors(verrors), nil))
		return ""
	}

	return strings.Split(headerValue, " ")[1]
}

func roleCheckerFromCtx(w http.ResponseWriter, r *http.Request) (string, bool) {
	userRole := r.Context().Value("user_role")

	var role string
	switch r := userRole.(type) {
	case string:
		role = r
	default:
		SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
		return "", false
	}

	return role, true
}

func AllowedRoleChecker(w http.ResponseWriter, req *http.Request, allowedRoles ...string) bool {
	role, ok := roleCheckerFromCtx(w, req)
	if !ok {
		return false
	}

	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return true
		}
	}

	SendJSONResponse(w, http.StatusForbidden, web.ErrorResponse(core.ForbiddenErrorCode, nil, nil))
	return false
}
