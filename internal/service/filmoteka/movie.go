package filmoteka

import (
	"context"
	"strings"
	"time"

	"vk-test-task/internal/core"
	"vk-test-task/internal/store/movie"
	"vk-test-task/pkg/web"
)

type (
	moviesService interface {
		GetMovies(context.Context, GetMoviesModel) ([]movie.Entity, int, error)
		GetMovieByID(context.Context, int) (movie.Entity, error)
		CreateMovie(context.Context, CreateMovieModel) (movie.Entity, error)
		UpdateMovie(context.Context, int, UpdateMovieModel) (movie.Entity, error)
		DeleteMovie(context.Context, int) error
	}

	GetMoviesModel struct {
		web.PaginationQuery
		SearchTerm string `query:"q" validate:"omitempty,max=150"`
		Sort       string `query:"sort" validate:"sort_params"`
	}

	CreateMovieModel struct {
		Title       string `json:"title" validate:"required,min=1,max=150"`
		Description string `json:"description" validate:"required,min=1,max=1000"`
		ReleaseDate string `json:"release_date" validate:"required,date"`
		Rating      int    `json:"rating" validate:"required,min=1,max=10"`
		StarsID     []int  `json:"stars_id" validate:"required"`
	}

	UpdateMovieModel struct {
		Title       *string `json:"title" validate:"required_without_all=Description ReleaseDate Rating,omitempty,min=1,max=150"`
		Description *string `json:"description" validate:"omitempty,min=1,max=1000"`
		ReleaseDate *string `json:"release_date" validate:"omitempty,date"`
		Rating      *int    `json:"rating" validate:"omitempty,min=0,max=10"`
		StarsID     []int   `json:"stars_id"`
	}
)

func (s *serviceImpl) GetMovies(ctx context.Context, model GetMoviesModel) ([]movie.Entity, int, error) {
	getAllParams := model.toGetAllParams()

	data, err := s.moviesStore.GetAll(ctx, getAllParams)
	if err != nil {
		return nil, 0, err
	}

	return data.Movies, data.TotalCount, nil
}

func (s *serviceImpl) GetMovieByID(ctx context.Context, id int) (movie.Entity, error) {
	data, err := s.moviesStore.GetByID(ctx, id)
	if err != nil {
		return movie.Entity{}, err
	}

	return data, nil
}

func (s *serviceImpl) CreateMovie(ctx context.Context, model CreateMovieModel) (movie.Entity, error) {
	entity, err := model.toCreateMovieEntity()
	if err != nil {
		return movie.Entity{}, err
	}

	data, err := s.moviesStore.Create(ctx, entity)
	if err != nil {
		return movie.Entity{}, err
	}

	return data, nil
}

func (s *serviceImpl) UpdateMovie(ctx context.Context, id int, model UpdateMovieModel) (movie.Entity, error) {
	for _, starID := range model.StarsID {
		exists, err := s.starsStore.CheckExistence(ctx, starID)
		if err != nil {
			return movie.Entity{}, err
		}
		if !exists {
			return movie.Entity{}, core.ErrStarIDNotExists
		}
	}
	entity, err := model.toUpdateMovieEntity()
	if err != nil {
		return movie.Entity{}, err
	}

	data, err := s.moviesStore.Update(ctx, id, entity)
	if err != nil {
		return movie.Entity{}, err
	}

	return data, nil
}

func (s *serviceImpl) DeleteMovie(ctx context.Context, id int) error {
	return s.moviesStore.Delete(ctx, id)
}

func (m GetMoviesModel) toGetAllParams() movie.GetAllParams {
	sortBy, sortOrder := "", ""
	if m.Sort != "" {
		splitted := strings.Split(m.Sort, ",")
		sortBy = splitted[0]
		sortOrder = splitted[1]

		if _, ok := core.AllowedSorts[splitted[0]]; !ok {
			sortBy = core.RatingOrder
		}
	}

	return movie.GetAllParams{
		SearchTerm: m.SearchTerm,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
		Limit:      m.GetLimit(),
		Offset:     m.GetOffset(),
	}
}

func (m CreateMovieModel) toCreateMovieEntity() (movie.CreateEntity, error) {
	date, err := time.Parse(time.RFC3339, m.ReleaseDate)
	if err != nil {
		return movie.CreateEntity{}, err
	}

	return movie.CreateEntity{
		Title:       m.Title,
		Description: m.Description,
		ReleaseDate: date,
		Rating:      m.Rating,
		StarsID:     m.StarsID,
	}, nil
}

func (m UpdateMovieModel) toUpdateMovieEntity() (movie.UpdateEntity, error) {
	var date *time.Time

	if m.ReleaseDate != nil {
		dateTime, err := time.Parse(time.RFC3339, *m.ReleaseDate)
		if err != nil {
			return movie.UpdateEntity{}, err
		}
		date = &dateTime
	}

	return movie.UpdateEntity{
		Title:       m.Title,
		Description: m.Description,
		ReleaseDate: date,
		Rating:      m.Rating,
		StarsID:     m.StarsID,
	}, nil
}
