package webutil

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"vk-test-task/internal/core"
	"vk-test-task/pkg/web"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validate *validator.Validate
}

const (
	TitleAsc        = "title,asc"
	TitleDesc       = "title,desc"
	RatingAsc       = "rating,asc"
	RatingDesc      = "rating,desc"
	ReleaseDateAsc  = "release_date,asc"
	ReleaseDateDesc = "release_date,desc"
	NoParam         = ""
)

func (cv *CustomValidator) ValidateSortParams(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	switch value {
	case TitleAsc, TitleDesc, RatingAsc, RatingDesc, ReleaseDateAsc, ReleaseDateDesc, NoParam:
		return true
	default:
		return false
	}
}

func (cv *CustomValidator) ValidateAuthHeader(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	splittedAuthHeader := strings.Split(value, " ")
	if len(splittedAuthHeader) != 2 || splittedAuthHeader[0] != "Bearer" {
		return false
	}

	return true
}

func (cv *CustomValidator) ValidateDate(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	_, err := time.Parse(time.RFC3339, value)
	return err == nil
}

func Validate(w http.ResponseWriter, r *http.Request, entity interface{}) bool {
	validate = validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(entity); err != nil {
		var verrors validator.ValidationErrors
		ok := errors.As(err, &verrors)
		if ok {
			SendJSONResponse(w, http.StatusUnprocessableEntity, web.ValidationErrorResponse(web.ValidationErrors(verrors), nil))
		} else {
			SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
		}
		return false
	}

	return true
}
