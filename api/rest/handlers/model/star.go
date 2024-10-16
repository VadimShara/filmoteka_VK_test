package model

import (
	"time"

	"vk-test-task/pkg/web"
)

type Star struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Sex       string     `json:"sex"`
	BirthDate time.Time  `json:"birth_date"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type StarWithMovies struct {
	Star   Star    `json:"star"`
	Movies []Movie `json:"movies"`
}

type CreateStarRequest struct {
	Name      string    `json:"name" example:"Ryan Gosling"`
	Sex       string    `json:"sex" example:"male"`
	BirthDate time.Time `json:"birth_date" example:"1980-11-12T00:00:00Z"`
}

type UpdateStarRequest struct {
	Name      string    `json:"name" example:"Raisa Goslingova"`
	Sex       string    `json:"sex" example:"female"`
	BirthDate time.Time `json:"birth_date" example:"1981-02-21T00:00:00Z"`
}

type GetStarsResponse struct {
	Status  string             `json:"status" example:"OK"`
	MsgCode string             `json:"msg_code" example:"stars_received"`
	Data    []Star             `json:"data" example:"[{\"id\":1,\"name\":\"Ryan Gosling\",\"sex\":\"male\",\"birth_date\":\"1980-11-12T00:00:00Z\",\"created_at\":\"2024-03-15T21:16:36Z\",\"updated_at\":\"2024-03-15T22:16:03Z\",\"deleted_at\":null},{\"id\":2,\"name\":\"Zendaya\",\"sex\":\"female\",\"birth_date\":\"1996-09-01T00:00:00Z\",\"created_at\":\"2024-03-15T20:04:29Z\",\"updated_at\":\"2024-03-15T20:04:29Z\",\"deleted_at\":null}]"`
	Meta    web.PaginationBody `json:"_meta" example:"{\"total_count\": 2, \"page_count\": 1, \"current_page\": 1, \"per_page\": 20}"`
}

type GetStarByIDResponse struct {
	Status  string         `json:"status" example:"OK"`
	MsgCode string         `json:"msg_code" example:"star_received"`
	Data    StarWithMovies `json:"data" example:"{\"star\":{\"id\":1,\"name\":\"Ryan Gosling\",\"sex\":\"male\",\"birth_date\":\"1980-11-12T00:00:00Z\",\"created_at\":\"2024-03-16T10:40:05Z\",\"updated_at\":\"2024-03-16T10:40:05Z\",\"deleted_at\":null},\"movies\":[{\"id\":1,\"title\":\"Drive\",\"description\":\"I Drive\",\"release_date\":\"2012-01-26T00:00:00Z\",\"rating\":9,\"created_at\":\"2024-03-16T10:41:18Z\",\"updated_at\":\"2024-03-16T10:41:18Z\",\"deleted_at\":null},{\"id\":2,\"title\":\"La La Land\",\"description\":\"I Dance\",\"release_date\":\"2017-01-12T00:00:00Z\",\"rating\":8,\"created_at\":\"2024-03-16T10:42:20Z\",\"updated_at\":\"2024-03-16T10:42:20Z\",\"deleted_at\":null}]}"`
}

type CreateStarResponse struct {
	Status  string         `json:"status" example:"OK"`
	MsgCode string         `json:"msg_code" example:"star_created"`
	Data    StarWithMovies `json:"data" example:"{\"star\":{\"id\":1,\"name\":\"Ryan Gosling\",\"sex\":\"male\",\"birth_date\":\"1980-11-12T00:00:00Z\",\"created_at\":\"2024-03-16T10:40:05Z\",\"updated_at\":\"2024-03-16T10:40:05Z\",\"deleted_at\":null},\"movies\":[]}"`
}

type UpdateStarResponse struct {
	Status  string         `json:"status" example:"OK"`
	MsgCode string         `json:"msg_code" example:"star_updated"`
	Data    StarWithMovies `json:"data" example:"{\"star\":{\"id\":1,\"name\":\"Raisa Goslingova\",\"sex\":\"female\",\"birth_date\":\"1981-02-21T00:00:00Z\",\"created_at\":\"2024-03-16T10:40:05Z\",\"updated_at\":\"2024-03-16T10:40:05Z\",\"deleted_at\":null},\"movies\":[{\"id\":1,\"title\":\"Drive\",\"description\":\"I Drive\",\"release_date\":\"2012-01-26T00:00:00Z\",\"rating\":9,\"created_at\":\"2024-03-16T10:41:18Z\",\"updated_at\":\"2024-03-16T10:41:18Z\",\"deleted_at\":null},{\"id\":2,\"title\":\"La La Land\",\"description\":\"I Dance\",\"release_date\":\"2017-01-12T00:00:00Z\",\"rating\":8,\"created_at\":\"2024-03-16T10:42:20Z\",\"updated_at\":\"2024-03-16T10:42:20Z\",\"deleted_at\":null}]}"`
}

type DeleteStarResponse struct {
	Status  string `json:"status" example:"OK"`
	MsgCode string `json:"msg_code" example:"star_deleted"`
}
