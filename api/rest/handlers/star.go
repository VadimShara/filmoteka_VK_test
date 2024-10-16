package handlers

import (
	"context"
	"net/http"
	"strconv"

	"vk-test-task/api/rest/presenters/star"
	"vk-test-task/internal/core"
	"vk-test-task/internal/service/filmoteka"
	"vk-test-task/pkg/logger"
	"vk-test-task/pkg/web"
	"vk-test-task/pkg/webutil"

	"github.com/jackc/pgx/v5"
)

func (r *Resolver) handleStars(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		if webutil.AllowedRoleChecker(w, req, core.AdminRole, core.UserRole) {
			r.getStars(w, req)
		}
	case http.MethodPost:
		if webutil.AllowedRoleChecker(w, req, core.AdminRole) {
			r.createStar(w, req)
		}
	default:
		webutil.SendJSONResponse(w, http.StatusMethodNotAllowed, web.ErrorResponse(core.UnsupportedMethodCode, nil, nil))
	}
}

func (r *Resolver) handleStar(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		if webutil.AllowedRoleChecker(w, req, core.AdminRole, core.UserRole) {
			r.getStarByID(w, req)
		}
	case http.MethodPatch:
		if webutil.AllowedRoleChecker(w, req, core.AdminRole) {
			r.updateStar(w, req)
		}
	case http.MethodDelete:
		if webutil.AllowedRoleChecker(w, req, core.AdminRole) {
			r.deleteStar(w, req)
		}
	default:
		webutil.SendJSONResponse(w, http.StatusMethodNotAllowed, web.ErrorResponse(core.UnsupportedMethodCode, nil, nil))
	}
}

// @Title Get Stars Paginated
// @Resource Stars
// @Description Get stars list paginated
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 array model.GetStarsResponse "Successful get stars"
// @Failure 401 object model.UnauthorizedResponse "Unauthorized error"
// @Failure 403 object model.ForbiddenResponse "Forbidden error"
// @Failure 500 object model.InternalResponse "Internal server error"
// @Route /api/v1/filmoteka/stars [get]
func (r *Resolver) getStars(w http.ResponseWriter, req *http.Request) {
	var model filmoteka.GetStarsModel

	queries := webutil.QueryParser(req)
	if _, ok := queries["page"]; ok {
		num, err := strconv.Atoi(queries["page"])
		if err != nil {
			logger.Log.Error("error parse page", "error", err.Error())
			webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return
		}
		model.Page = num
	}

	if _, ok := queries["limit"]; ok {
		num, err := strconv.Atoi(queries["limit"])
		if err != nil {
			logger.Log.Error("error parse limit", "error", err.Error())
			webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return
		}
		model.Limit = num
	}
	webutil.QueryValidator(w, req, &model)

	data, total, err := r.filmotekaService.GetStars(context.Background(), model)
	if err != nil {
		webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
		return
	}

	pres := star.PresentList(data, model.PaginationQuery, total)

	webutil.SendJSONResponse(w, http.StatusOK, pres.Response())
}

// @Title Get Star By ID
// @Resource Stars
// @Param id path int true "Star ID"
// @Success 200 object model.GetStarByIDResponse "Successful get star"
// @Failure 400 object model.BadRequestInvalidIDResponse "Bad request error"
// @Failure 401 object model.UnauthorizedResponse "Unauthorized error"
// @Failure 403 object model.ForbiddenResponse "Forbidden error"
// @Failure 404 object model.StarNotFoundResponse "Not found error"
// @Failure 500 object model.InternalResponse "Internal server error"
// @Route /api/v1/filmoteka/star/{id} [get]
func (r *Resolver) getStarByID(w http.ResponseWriter, req *http.Request) {
	id := webutil.ParseID(w, req)
	if id == 0 {
		return
	}

	starData, moviesData, err := r.filmotekaService.GetStarByID(context.Background(), id)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			webutil.SendJSONResponse(w, http.StatusNotFound, web.ErrorResponse(core.StarNotFoundCode, nil, nil))
			return
		default:
			webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return
		}
	}

	pres := star.PresentStar(starData, moviesData)

	webutil.SendJSONResponse(w, http.StatusOK, pres.Response(core.StarReceivedCode))
}

// @Title Create Star
// @Resource Stars
// @Param star body model.CreateStarRequest true "Star to create"
// @Success 201 object model.CreateStarResponse "Successful create star"
// @Failure 400 object model.BadRequestInvalidBodyResponse "Bad request error"
// @Failure 401 object model.UnauthorizedResponse "Unauthorized error"
// @Failure 403 object model.ForbiddenResponse "Forbidden error"
// @Failure 422 object model.ValidationResponse "Validation error"
// @Failure 500 object model.InternalResponse "Internal server error"
// @Route /api/v1/filmoteka/stars [post]
func (r *Resolver) createStar(w http.ResponseWriter, req *http.Request) {
	var model filmoteka.CreateStarModel

	if !webutil.BodyCheck(w, req, &model) {
		return
	}

	data, err := r.filmotekaService.CreateStar(context.Background(), model)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			webutil.SendJSONResponse(w, http.StatusNotFound, web.ErrorResponse(core.StarNotFoundCode, nil, nil))
			return
		default:
			webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return
		}
	}

	pres := star.PresentStar(data, nil)

	webutil.SendJSONResponse(w, http.StatusCreated, pres.Response(core.StarCreatedCode))
}

// @Title Update Star
// @Resource Stars
// @Param id path int true "Star ID"
// @Param star body model.UpdateStarRequest true "Star to update"
// @Success 200 object model.UpdateStarResponse "Successful update star"
// @Failure 400 object model.BadRequestInvalidBodyResponse "Bad request error"
// @Failure 401 object model.UnauthorizedResponse "Unauthorized error"
// @Failure 403 object model.ForbiddenResponse "Forbidden error"
// @Failure 404 object model.StarNotFoundResponse "Not found error"
// @Failure 422 object model.ValidationResponse "Validation error"
// @Failure 500 object model.InternalResponse "Internal server error"
// @Route /api/v1/filmoteka/star/{id} [patch]
func (r *Resolver) updateStar(w http.ResponseWriter, req *http.Request) {
	id := webutil.ParseID(w, req)
	if id == 0 {
		return
	}

	var model filmoteka.UpdateStarModel

	if !webutil.BodyCheck(w, req, &model) {
		return
	}

	starData, moviesData, err := r.filmotekaService.UpdateStar(context.Background(), id, model)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			webutil.SendJSONResponse(w, http.StatusNotFound, web.ErrorResponse(core.StarNotFoundCode, nil, nil))
			return
		default:
			webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return
		}
	}

	pres := star.PresentStar(starData, moviesData)

	webutil.SendJSONResponse(w, http.StatusOK, pres.Response(core.StarUpdatedCode))
}

// @Title Delete Star
// @Resource Stars
// @Param id path int true "Star ID"
// @Success 200 object model.DeleteStarResponse "Successful delete star"
// @Failure 400 object model.BadRequestInvalidIDResponse "Bad request error"
// @Failure 401 object model.UnauthorizedResponse "Unauthorized error"
// @Failure 403 object model.ForbiddenResponse "Forbidden error"
// @Failure 404 object model.StarNotFoundResponse "Not found error"
// @Failure 500 object model.InternalResponse "Internal server error"
// @Route /api/v1/filmoteka/star/{id} [delete]
func (r *Resolver) deleteStar(w http.ResponseWriter, req *http.Request) {
	id := webutil.ParseID(w, req)
	if id == 0 {
		return
	}

	err := r.filmotekaService.DeleteStar(context.Background(), id)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			webutil.SendJSONResponse(w, http.StatusNotFound, web.ErrorResponse(core.StarNotFoundCode, nil, nil))
			return
		default:
			webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return
		}
	}

	webutil.SendJSONResponse(w, http.StatusOK, web.OKResponse(core.StarDeletedCode, nil, nil))
}
