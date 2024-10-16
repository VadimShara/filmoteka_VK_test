package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"vk-test-task/api/rest/presenters/movie"
	"vk-test-task/internal/core"
	"vk-test-task/internal/service/filmoteka"
	"vk-test-task/pkg/logger"
	"vk-test-task/pkg/web"
	"vk-test-task/pkg/webutil"

	"github.com/jackc/pgx/v5"
)

func (r *Resolver) handleMovies(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		if webutil.AllowedRoleChecker(w, req, core.AdminRole, core.UserRole) {
			r.getMovies(w, req)
		}
	case http.MethodPost:
		if webutil.AllowedRoleChecker(w, req, core.AdminRole) {
			r.createMovie(w, req)
		}
	default:
		webutil.SendJSONResponse(w, http.StatusMethodNotAllowed, web.ErrorResponse(core.UnsupportedMethodCode, nil, nil))
	}
}

func (r *Resolver) handleMovie(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		if webutil.AllowedRoleChecker(w, req, core.AdminRole, core.UserRole) {
			r.getMovieByID(w, req)
		}
	case http.MethodPatch:
		if webutil.AllowedRoleChecker(w, req, core.AdminRole) {
			r.updateMovie(w, req)
		}
	case http.MethodDelete:
		if webutil.AllowedRoleChecker(w, req, core.AdminRole) {
			r.deleteMovie(w, req)
		}
	default:
		webutil.SendJSONResponse(w, http.StatusMethodNotAllowed, web.ErrorResponse(core.UnsupportedMethodCode, nil, nil))
	}
}

// @Title Get Movies Paginated
// @Resource Movies
// @Description Get movies list paginated (with sorting and search term)
// @Param q query string false "Search term"
// @Param sort query string false "Sort result"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 array model.GetMoviesResponse "Successful get movies"
// @Failure 401 object model.UnauthorizedResponse "Unauthorized error"
// @Failure 403 object model.ForbiddenResponse "Forbidden error"
// @Failure 500 object model.InternalResponse "Internal server error"
// @Route /api/v1/filmoteka/movies [get]
func (r *Resolver) getMovies(w http.ResponseWriter, req *http.Request) {
	var model filmoteka.GetMoviesModel

	queries := webutil.QueryParser(req)
	if _, ok := queries["page"]; ok {
		num, err := strconv.Atoi(queries["page"])
		if err != nil {
			logger.Log.Error("error parse page", "error", err.Error())
			webutil.SendJSONResponse(w, http.StatusBadRequest, web.ErrorResponse(core.BadRequestErrorCode, nil, nil))
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
	if val, ok := queries["q"]; ok {
		model.SearchTerm = val
	}
	if val, ok := queries["sort"]; ok {
		model.Sort = val
	}
	if !webutil.QueryValidator(w, req, &model) {
		return
	}

	data, total, err := r.filmotekaService.GetMovies(context.Background(), model)
	if err != nil {
		webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
		return
	}

	pres := movie.PresentList(data, model.PaginationQuery, total)

	webutil.SendJSONResponse(w, http.StatusOK, pres.Response())
}

// @Title Get Movie By ID
// @Resource Movies
// @Param id path int true "Movie ID"
// @Success 200 object model.GetMovieByIDResponse "Successful get movie"
// @Failure 400 object model.BadRequestInvalidIDResponse "Bad request error"
// @Failure 401 object model.UnauthorizedResponse "Unauthorized error"
// @Failure 403 object model.ForbiddenResponse "Forbidden error"
// @Failure 404 object model.MovieNotFoundResponse "Not found error"
// @Failure 500 object model.InternalResponse "Internal server error"
// @Route /api/v1/filmoteka/movie/{id} [get]
func (r *Resolver) getMovieByID(w http.ResponseWriter, req *http.Request) {
	id := webutil.ParseID(w, req)
	if id == 0 {
		return
	}

	data, err := r.filmotekaService.GetMovieByID(context.Background(), id)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			webutil.SendJSONResponse(w, http.StatusNotFound, web.ErrorResponse(core.MovieNotFoundCode, nil, nil))
			return
		default:
			webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return
		}
	}

	pres := movie.PresentMovie(data)

	webutil.SendJSONResponse(w, http.StatusOK, pres.Response(core.MovieReceivedCode))
}

// @Title Create Movie
// @Resource Movies
// @Param movie body model.CreateMovieRequest true "Movie to create"
// @Success 201 object model.CreateMovieResponse "Successful create movie"
// @Failure 400 object model.BadRequestInvalidBodyResponse "Bad request error"
// @Failure 401 object model.UnauthorizedResponse "Unauthorized error"
// @Failure 403 object model.ForbiddenResponse "Forbidden error"
// @Failure 422 object model.ValidationResponse "Validation error"
// @Failure 500 object model.InternalResponse "Internal server error"
// @Route /api/v1/filmoteka/movies [post]
func (r *Resolver) createMovie(w http.ResponseWriter, req *http.Request) {
	var model filmoteka.CreateMovieModel

	if !webutil.BodyCheck(w, req, &model) {
		return
	}

	data, err := r.filmotekaService.CreateMovie(context.Background(), model)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrStarIDNotExists):
			webutil.SendJSONResponse(w, http.StatusNotFound, web.ErrorResponse(core.StarNotFoundCode, nil, nil))
			return
		case err == pgx.ErrNoRows:
			webutil.SendJSONResponse(w, http.StatusNotFound, web.ErrorResponse(core.MovieNotFoundCode, nil, nil))
			return
		default:
			webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return
		}
	}

	pres := movie.PresentMovie(data)

	webutil.SendJSONResponse(w, http.StatusCreated, pres.Response(core.MovieCreatedCode))
}

// @Title Update Movie
// @Resource Movies
// @Param id path int true "Movie ID"
// @Param movie body model.UpdateMovieRequest true "Movie to update"
// @Success 200 object model.UpdateMovieResponse "Successful update movie"
// @Failure 400 object model.BadRequestInvalidBodyResponse "Bad request error"
// @Failure 401 object model.UnauthorizedResponse "Unauthorized error"
// @Failure 403 object model.ForbiddenResponse "Forbidden error"
// @Failure 404 object model.MovieNotFoundResponse "Not found error"
// @Failure 422 object model.ValidationResponse "Validation error"
// @Failure 500 object model.InternalResponse "Internal server error"
// @Route /api/v1/filmoteka/movie/{id} [patch]
func (r *Resolver) updateMovie(w http.ResponseWriter, req *http.Request) {
	id := webutil.ParseID(w, req)
	if id == 0 {
		return
	}

	var model filmoteka.UpdateMovieModel

	if !webutil.BodyCheck(w, req, &model) {
		return
	}

	data, err := r.filmotekaService.UpdateMovie(context.Background(), id, model)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrStarIDNotExists):
			webutil.SendJSONResponse(w, http.StatusNotFound, web.ErrorResponse(core.StarNotFoundCode, nil, nil))
			return
		case err == pgx.ErrNoRows:
			webutil.SendJSONResponse(w, http.StatusNotFound, web.ErrorResponse(core.MovieNotFoundCode, nil, nil))
			return
		default:
			webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return
		}
	}

	pres := movie.PresentMovie(data)

	webutil.SendJSONResponse(w, http.StatusOK, pres.Response(core.MovieUpdatedCode))
}

// @Title Delete Movie
// @Resource Movies
// @Param id path int true "Movie ID"
// @Success 200 object model.DeleteMovieResponse "Successful delete movie"
// @Failure 400 object model.BadRequestInvalidIDResponse "Bad request error"
// @Failure 401 object model.UnauthorizedResponse "Unauthorized error"
// @Failure 403 object model.ForbiddenResponse "Forbidden error"
// @Failure 404 object model.MovieNotFoundResponse "Not found error"
// @Failure 500 object model.InternalResponse "Internal server error"
// @Route /api/v1/filmoteka/movie/{id} [delete]
func (r *Resolver) deleteMovie(w http.ResponseWriter, req *http.Request) {
	id := webutil.ParseID(w, req)
	if id == 0 {
		return
	}

	err := r.filmotekaService.DeleteMovie(context.Background(), id)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			webutil.SendJSONResponse(w, http.StatusNotFound, web.ErrorResponse(core.MovieNotFoundCode, nil, nil))
			return
		default:
			webutil.SendJSONResponse(w, http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return
		}
	}

	webutil.SendJSONResponse(w, http.StatusOK, web.OKResponse(core.MovieDeletedCode, nil, nil))
}
