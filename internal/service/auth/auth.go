package auth

import (
	"context"
	"net/http"

	"vk-test-task/internal/core"
	"vk-test-task/internal/store/user"
	"vk-test-task/pkg/hash"
	"vk-test-task/pkg/jwt"
)

type (
	Service interface {
		SignUp(context.Context, SignUpModel) (user.Entity, error)
		GetPassHashAndRoleByUsername(context.Context, string) (string, string, error)
		CreateToken(context.Context, string, string) (jwt.Token, error)
		Verify(http.ResponseWriter, *http.Request) (*jwt.UserData, bool)
	}

	SignUpModel struct {
		Username string `json:"username" validate:"required,min=1,max=16"`
		Password string `json:"password" validate:"required,min=1,max=100"`
		Role     string `json:"role" validate:"required,oneof=admin user"`
	}

	LoginModel struct {
		Username string `json:"username" validate:"required,min=1,max=16"`
		Password string `json:"password" validate:"required,min=1,max=100"`
	}

	serviceImpl struct {
		jwtService *jwt.Service
		usersStore user.Store
	}
)

func New(users user.Store) (Service, error) {
	cfg, err := jwt.ParseConfig()
	if err != nil {
		return nil, err
	}

	return &serviceImpl{
		jwtService: jwt.New(cfg),
		usersStore: users,
	}, nil
}

func (s *serviceImpl) SignUp(ctx context.Context, signUpModel SignUpModel) (user.Entity, error) {
	exists, err := s.usersStore.CheckExistence(ctx, signUpModel.Username)
	if err != nil {
		return user.Entity{}, err
	}
	if exists {
		return user.Entity{}, core.ErrUsernameExists
	}

	entity := signUpModel.toCreateUserEntity()

	data, err := s.usersStore.Create(ctx, entity)
	if err != nil {
		return user.Entity{}, err
	}

	return data, nil
}

func (s *serviceImpl) GetPassHashAndRoleByUsername(ctx context.Context, username string) (string, string, error) {
	return s.usersStore.GetPassHashAndRoleByUsername(context.Background(), username)
}

func (s *serviceImpl) CreateToken(ctx context.Context, username, role string) (jwt.Token, error) {
	return s.jwtService.CreateToken(ctx, jwt.UserData{Username: username, Role: role})
}

func (s *serviceImpl) Verify(w http.ResponseWriter, r *http.Request) (*jwt.UserData, bool) {
	return s.jwtService.Verify(w, r)
}

func (m SignUpModel) toCreateUserEntity() user.CreateEntity {
	passHash := hash.CalculateHash(m.Password)

	return user.CreateEntity{
		Username: m.Username,
		PassHash: passHash,
		Role:     m.Role,
	}
}
