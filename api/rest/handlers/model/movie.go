package model

import (
	"time"

	"vk-test-task/pkg/web"
)

type Movie struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ReleaseDate time.Time  `json:"release_date"`
	Rating      int        `json:"rating"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type CreateMovieRequest struct {
	Title       string    `json:"title" example:"Title"`
	Description string    `json:"description" example:"Description"`
	ReleaseDate time.Time `json:"release_date" example:"2024-06-06T00:00:00Z"`
	Rating      int       `json:"rating" example:"9"`
}

type UpdateMovieRequest struct {
	Title       string    `json:"title" example:"New Title"`
	Description string    `json:"description" example:"New Description"`
	ReleaseDate time.Time `json:"release_date" example:"2024-10-10T00:00:00Z"`
	Rating      int       `json:"rating" example:"8"`
}

type GetMoviesResponse struct {
	Status  string             `json:"status" example:"OK"`
	MsgCode string             `json:"msg_code" example:"movies_received"`
	Data    []Movie            `json:"data" example:"[{\"id\":1,\"title\":\"Drive\",\"description\":\"Night Call\",\"release_date\":\"2012-01-26T00:00:00Z\",\"rating\":8,\"created_at\":\"2024-03-15T21:16:36Z\",\"updated_at\":\"2024-03-15T22:16:03Z\",\"deleted_at\":null},{\"id\":2,\"title\":\"Oppenheimer\",\"description\":\"Boom\",\"release_date\":\"2023-07-20T00:00:00Z\",\"rating\":9,\"created_at\":\"2024-03-15T20:04:29Z\",\"updated_at\":\"2024-03-15T20:04:29Z\",\"deleted_at\":null}]"`
	Meta    web.PaginationBody `json:"_meta" example:"{\"total_count\": 2, \"page_count\": 1, \"current_page\": 1, \"per_page\": 20}"`
}

type GetMovieByIDResponse struct {
	Status  string `json:"status" example:"OK"`
	MsgCode string `json:"msg_code" example:"movie_received"`
	Data    Movie  `json:"data" example:"{\"id\":1,\"title\":\"Drive\",\"description\":\"Night Call\",\"release_date\":\"2012-26-01T00:00:00Z\",\"rating\":8,\"created_at\":\"2024-03-15T21:16:36Z\",\"updated_at\":\"2024-03-15T22:16:03Z\",\"deleted_at\":null}"`
}

type CreateMovieResponse struct {
	Status  string `json:"status" example:"OK"`
	MsgCode string `json:"msg_code" example:"movie_created"`
	Data    Movie  `json:"data" example:"{\"id\":1,\"title\":\"Drive\",\"description\":\"Night Call\",\"release_date\":\"2012-26-01T00:00:00Z\",\"rating\":8,\"created_at\":\"2024-03-15T21:16:36Z\",\"updated_at\":\"2024-03-15T22:16:03Z\",\"deleted_at\":null}"`
}

type UpdateMovieResponse struct {
	Status  string `json:"status" example:"OK"`
	MsgCode string `json:"msg_code" example:"movie_updated"`
	Data    Movie  `json:"data" example:"{\"id\":1,\"title\":\"Drive 2\",\"description\":\"Night Call 2\",\"release_date\":\"2022-02-01T00:00:00Z\",\"rating\":10,\"created_at\":\"2024-03-15T21:16:36Z\",\"updated_at\":\"2024-03-16T14:16:03Z\",\"deleted_at\":null}"`
}

type DeleteMovieResponse struct {
	Status  string `json:"status" example:"OK"`
	MsgCode string `json:"msg_code" example:"movie_deleted"`
}
