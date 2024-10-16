package user

import (
	"context"
	"testing"

	"vk-test-task/internal/tests"
	"vk-test-task/pkg/hash"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testUsername = "testuser"
	testPassHash = hash.CalculateHash("testpassword")
)

func TestCreate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresContainer, databaseURL := tests.CreateFilmotekaTestPostgresContainer(ctx)
	defer tests.TerminateFilmotekaTestContainer(ctx, postgresContainer)

	postgresClient, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Errorf("error with starting postgres client: %s", err.Error())
	}
	defer postgresClient.Close()

	store := New(postgresClient)

	for _, test := range []struct {
		Name    string
		Data    CreateEntity
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful creating an admin",
			Data: CreateEntity{
				Username: "RyanGosling2011",
				PassHash: hash.CalculateHash("idrive"),
				Role:     "admin",
			},
		},

		{
			Name: "Successful creating a user",
			Data: CreateEntity{
				Username: "PeterGriffin",
				PassHash: hash.CalculateHash("HeyLois"),
				Role:     "user",
			},
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			entity, err := store.Create(ctx, test.Data)
			if err != nil {
				if !test.WantErr {
					t.Errorf("unexpected error: %s", err.Error())
				}
				if err.Error() != test.Err {
					t.Errorf("unexpected error. Expected %q but got %q", test.Err, err.Error())
				}
				return
			}
			if test.WantErr {
				t.Errorf("expected error but nothing got")
			}
			if entity.Username != test.Data.Username {
				t.Errorf("wrong username. Expected %q but got %q", test.Data.Username, entity.Username)
			}
			if entity.Role != test.Data.Role {
				t.Errorf("wrong role. Expected %q but got %q", test.Data.Role, entity.Role)
			}
		})
	}
}

func TestGetPassHashAndRoleByUsername(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresContainer, databaseURL := tests.CreateFilmotekaTestPostgresContainer(ctx)
	defer tests.TerminateFilmotekaTestContainer(ctx, postgresContainer)

	postgresClient, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Errorf("error with starting postgres client: %s", err.Error())
	}
	defer postgresClient.Close()

	store := New(postgresClient)

	_, err = store.Create(ctx, CreateEntity{Username: testUsername, PassHash: testPassHash, Role: "user"})
	if err != nil {
		t.Errorf("error with adding existing user: %s", err.Error())
	}

	for _, test := range []struct {
		Name     string
		Username string
		Role     string
		WantErr  bool
		Err      string
	}{
		{
			Name:     "Successful get pass hash and role by username",
			Username: testUsername,
			Role:     "user",
		},
		{
			Name:     "Get pass hash by non-existent username",
			Username: "coolGOproger1337",
			WantErr:  true,
			Err:      pgx.ErrNoRows.Error(),
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			passHash, role, err := store.GetPassHashAndRoleByUsername(ctx, test.Username)
			if err != nil {
				if !test.WantErr {
					t.Errorf("unexpected error: %s", err.Error())
				}
				if err.Error() != test.Err {
					t.Errorf("unexpected error. Expected %q but got %q", test.Err, err.Error())
				}
				return
			}
			if test.WantErr {
				t.Errorf("expected error but nothing got")
			}
			if passHash != testPassHash {
				t.Errorf("wrong pass hash. Expected %q but got %q", testPassHash, passHash)
			}
			if role != test.Role {
				t.Errorf("wrong role. Expected %q but got %q", test.Role, role)
			}
		})
	}
}

func TestCheckExistence(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresContainer, databaseURL := tests.CreateFilmotekaTestPostgresContainer(ctx)
	defer tests.TerminateFilmotekaTestContainer(ctx, postgresContainer)

	postgresClient, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Errorf("error with starting postgres client: %s", err.Error())
	}
	defer postgresClient.Close()

	store := New(postgresClient)

	_, err = store.Create(ctx, CreateEntity{
		Username: testUsername,
		PassHash: testPassHash,
		Role:     "user",
	})
	if err != nil {
		t.Errorf("error with creating test data: %s", err.Error())
	}

	for _, test := range []struct {
		Name     string
		Username string
		Exists   bool
		WantErr  bool
		Err      string
	}{
		{
			Name:     "Check existent user",
			Username: testUsername,
			Exists:   true,
		},
		{
			Name:     "Check non-existent user",
			Username: "superMegaUser9000",
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			exists, err := store.CheckExistence(ctx, test.Username)
			if err != nil {
				if !test.WantErr {
					t.Errorf("unexpected error: %s", err.Error())
				}
				if err.Error() != test.Err {
					t.Errorf("unexpected error. Expected %q but got %q", test.Err, err.Error())
				}
				return
			}
			if test.WantErr {
				t.Errorf("expected error but nothing got")
			}
			if test.Exists != exists {
				t.Errorf("wrong existence. Expected %t but got %t", test.Exists, exists)
			}
		})
	}
}
