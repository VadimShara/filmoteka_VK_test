package star

import (
	"time"

	"vk-test-task/api/rest/presenters/movie"
	"vk-test-task/internal/core"
	moviestore "vk-test-task/internal/store/movie"
	"vk-test-task/internal/store/star"
	"vk-test-task/pkg/web"
)

type Presenter struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Sex       string     `json:"sex"`
	BirthDate time.Time  `json:"birth_date"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type PresenterWithMovies struct {
	Star   Presenter         `json:"star"`
	Movies []movie.Presenter `json:"movies"`
}

func PresentStar(entity star.Entity, movies []moviestore.Entity) PresenterWithMovies {
	moviesList := make([]movie.Presenter, len(movies))
	for i, m := range movies {
		moviesList[i] = movie.Presenter{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			ReleaseDate: m.ReleaseDate,
			Rating:      m.Rating,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
			DeletedAt:   m.DeletedAt,
		}
	}

	return PresenterWithMovies{
		Star: Presenter{
			ID:        entity.ID,
			Name:      entity.Name,
			Sex:       entity.Sex,
			BirthDate: entity.BirthDate,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
			DeletedAt: entity.DeletedAt,
		},
		Movies: moviesList,
	}
}

func (p *PresenterWithMovies) Response(msg string) web.Response {
	return web.OKResponse(msg, *p, nil)
}

type ListPresenter struct {
	stars      []Presenter
	pagination web.PaginationBody
}

func PresentList(entities []star.Entity, pq web.PaginationQuery, total int) ListPresenter {
	pres := ListPresenter{}

	for _, entity := range entities {
		starPresenter := Presenter{
			ID:        entity.ID,
			Name:      entity.Name,
			Sex:       entity.Sex,
			BirthDate: entity.BirthDate,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
			DeletedAt: entity.DeletedAt,
		}
		pres.stars = append(pres.stars, starPresenter)
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
	return web.OKResponse(core.StarsReceivedCode, p.stars, p.pagination)
}
