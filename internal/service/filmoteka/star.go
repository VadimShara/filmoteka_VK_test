package filmoteka

import (
	"context"
	"time"

	"vk-test-task/internal/store/movie"
	"vk-test-task/internal/store/star"
	"vk-test-task/pkg/web"
)

type (
	starsService interface {
		GetStars(context.Context, GetStarsModel) ([]star.Entity, int, error)
		GetStarByID(context.Context, int) (star.Entity, []movie.Entity, error)
		CreateStar(context.Context, CreateStarModel) (star.Entity, error)
		UpdateStar(context.Context, int, UpdateStarModel) (star.Entity, []movie.Entity, error)
		DeleteStar(context.Context, int) error
	}

	GetStarsModel struct {
		web.PaginationQuery
	}

	CreateStarModel struct {
		Name      string `json:"name" validate:"required,min=1,max=100"`
		Sex       string `json:"sex" validate:"required,oneof=male female"`
		BirthDate string `json:"birth_date" validate:"date"`
	}

	UpdateStarModel struct {
		Name      *string `json:"name" validate:"required_without_all=Sex BirthDate,omitempty,min=1,max=100"`
		Sex       *string `json:"sex" validate:"omitempty,oneof=male female"`
		BirthDate *string `json:"birth_date" validate:"omitempty,date"`
	}
)

func (s *serviceImpl) GetStars(ctx context.Context, model GetStarsModel) ([]star.Entity, int, error) {
	getAllParams := model.toGetAllParams()

	data, err := s.starsStore.GetAll(ctx, getAllParams)
	if err != nil {
		return nil, 0, err
	}

	return data.Stars, data.TotalCount, nil
}

func (s *serviceImpl) GetStarByID(ctx context.Context, id int) (star.Entity, []movie.Entity, error) {
	starData, err := s.starsStore.GetByID(ctx, id)
	if err != nil {
		return star.Entity{}, []movie.Entity{}, err
	}

	moviesData, err := s.moviesStore.GetByStarID(ctx, id)
	if err != nil {
		return star.Entity{}, []movie.Entity{}, err
	}

	return starData, moviesData, nil
}

func (s *serviceImpl) CreateStar(ctx context.Context, model CreateStarModel) (star.Entity, error) {
	entity, err := model.toCreateStarEntity()
	if err != nil {
		return star.Entity{}, err
	}

	data, err := s.starsStore.Create(ctx, entity)
	if err != nil {
		return star.Entity{}, err
	}

	return data, nil
}

func (s *serviceImpl) UpdateStar(ctx context.Context, id int, model UpdateStarModel) (star.Entity, []movie.Entity, error) {
	entity, err := model.toUpdateStarEntity()
	if err != nil {
		return star.Entity{}, []movie.Entity{}, err
	}

	data, err := s.starsStore.Update(ctx, id, entity)
	if err != nil {
		return star.Entity{}, []movie.Entity{}, err
	}

	moviesData, err := s.moviesStore.GetByStarID(ctx, id)
	if err != nil {
		return star.Entity{}, []movie.Entity{}, err
	}

	return data, moviesData, nil
}

func (s *serviceImpl) DeleteStar(ctx context.Context, id int) error {
	return s.starsStore.Delete(ctx, id)
}

func (m GetStarsModel) toGetAllParams() star.GetAllParams {
	return star.GetAllParams{
		Limit:  m.PaginationQuery.GetLimit(),
		Offset: m.PaginationQuery.GetOffset(),
	}
}

func (m CreateStarModel) toCreateStarEntity() (star.CreateEntity, error) {
	date, err := time.Parse(time.RFC3339, m.BirthDate)
	if err != nil {
		return star.CreateEntity{}, err
	}

	return star.CreateEntity{
		Name:      m.Name,
		Sex:       m.Sex,
		BirthDate: date,
	}, nil
}

func (m UpdateStarModel) toUpdateStarEntity() (star.UpdateEntity, error) {
	var date *time.Time

	if m.BirthDate != nil {
		dateTime, err := time.Parse(time.RFC3339, *m.BirthDate)
		if err != nil {
			return star.UpdateEntity{}, err
		}
		date = &dateTime
	}

	return star.UpdateEntity{
		Name:      m.Name,
		Sex:       m.Sex,
		BirthDate: date,
	}, nil
}
