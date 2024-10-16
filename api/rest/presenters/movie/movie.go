package movie

import (
	"time"

	"vk-test-task/internal/core"
	"vk-test-task/internal/store/movie"
	"vk-test-task/pkg/web"
)

type Presenter struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ReleaseDate time.Time  `json:"release_date"`
	Rating      int        `json:"rating"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func PresentMovie(entity movie.Entity) Presenter {
	return Presenter{
		ID:          entity.ID,
		Title:       entity.Title,
		Description: entity.Description,
		ReleaseDate: entity.ReleaseDate,
		Rating:      entity.Rating,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		DeletedAt:   entity.DeletedAt,
	}
}

func (p *Presenter) Response(msg string) web.Response {
	return web.OKResponse(msg, *p, nil)
}

type ListPresenter struct {
	movies     []Presenter
	pagination web.PaginationBody
}

func PresentList(entities []movie.Entity, pq web.PaginationQuery, total int) ListPresenter {
	pres := ListPresenter{}

	for _, entity := range entities {
		moviePresenter := Presenter{
			ID:          entity.ID,
			Title:       entity.Title,
			Description: entity.Description,
			ReleaseDate: entity.ReleaseDate,
			Rating:      entity.Rating,
			CreatedAt:   entity.CreatedAt,
			UpdatedAt:   entity.UpdatedAt,
			DeletedAt:   entity.DeletedAt,
		}
		pres.movies = append(pres.movies, moviePresenter)
	}

	pres.pagination = web.PaginationBody{
		TotalCount:  total,
		PageCount:   pq.GetPageCount(total),
		CurrentPage: pq.GetPage(),
		PerPage:     pq.GetLimit(),
	}

	return pres
}

func (p *ListPresenter) Response() web.Response {
	return web.OKResponse(core.MoviesReceivedCode, p.movies, p.pagination)
}
