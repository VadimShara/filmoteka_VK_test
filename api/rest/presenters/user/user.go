package user

import (
	"time"

	"vk-test-task/internal/store/user"
	"vk-test-task/pkg/web"
)

type Presenter struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func PresentUser(entity user.Entity) Presenter {
	return Presenter{
		ID:        entity.ID,
		Username:  entity.Username,
		Role:      entity.Role,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func (p *Presenter) Response(msg string) web.Response {
	return web.OKResponse(msg, *p, nil)
}
