package star

import (
	"context"
	"testing"
	"time"

	"vk-test-task/internal/tests"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
			Name: "Successful creating a star",
			Data: CreateEntity{
				Name:      "Ryan Gosling",
				Sex:       "male",
				BirthDate: time.Date(1980, time.November, 12, 0, 0, 0, 0, time.UTC),
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
			if entity.Name != test.Data.Name {
				t.Errorf("wrong name. Expected %q but got %q", test.Data.Name, entity.Name)
			}
			if entity.Sex != test.Data.Sex {
				t.Errorf("wrong sex. Expected %q but got %q", test.Data.Sex, entity.Sex)
			}
			if entity.BirthDate != test.Data.BirthDate {
				t.Errorf("wrong birth date. Expected %q but got %q", test.Data.BirthDate, entity.BirthDate)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
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

	star, err := store.Create(ctx, CreateEntity{
		Name:      "Ryan Gosling",
		Sex:       "male",
		BirthDate: time.Date(1980, time.November, 12, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Errorf("error with creating test data: %s", err.Error())
	}

	for _, test := range []struct {
		Name    string
		ID      int
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful get a star",
			ID:   star.ID,
		},
		{
			Name:    "Get non-existent star",
			WantErr: true,
			Err:     pgx.ErrNoRows.Error(),
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			entity, err := store.GetByID(ctx, test.ID)
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
			if entity.Name != star.Name {
				t.Errorf("wrong name. Expected %q but got %q", star.Name, entity.Name)
			}
			if entity.Sex != star.Sex {
				t.Errorf("wrong sex. Expected %q but got %q", star.Sex, entity.Sex)
			}
			if entity.BirthDate != star.BirthDate {
				t.Errorf("wrong birth date. Expected %q but got %q", star.BirthDate, entity.BirthDate)
			}
			if entity.CreatedAt != star.CreatedAt {
				t.Errorf("wrong created at. Expected %q but got %q", star.CreatedAt, entity.CreatedAt)
			}
			if entity.UpdatedAt != star.UpdatedAt {
				t.Errorf("wrong updated at. Expected %q but got %q", star.UpdatedAt, entity.UpdatedAt)
			}
			if entity.DeletedAt != nil {
				t.Error("wrong deleted at. Expected nil")
			}
		})
	}
}

func TestGetAll(t *testing.T) {
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

	var stars []Entity
	for _, data := range []CreateEntity{
		{
			Name:      "Ryan Gosling",
			Sex:       "male",
			BirthDate: time.Date(1980, time.November, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:      "Leonardo DiCaprio",
			Sex:       "male",
			BirthDate: time.Date(1974, time.November, 11, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:      "Emma Stone",
			Sex:       "female",
			BirthDate: time.Date(1988, time.November, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:      "Zendaya",
			Sex:       "female",
			BirthDate: time.Date(1996, time.September, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:      "Timothee Chalamet",
			Sex:       "male",
			BirthDate: time.Date(1995, time.December, 27, 0, 0, 0, 0, time.UTC),
		},
	} {
		entity, err := store.Create(ctx, data)
		if err != nil {
			t.Errorf("error with creating stars: %s", err.Error())
		}

		stars = append(stars, entity)
	}

	for _, test := range []struct {
		Name    string
		Data    GetAllParams
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful get all stars paginated",
			Data: GetAllParams{
				Limit: len(stars),
			},
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			entity, err := store.GetAll(ctx, test.Data)
			if err != nil {
				if !test.WantErr {
					t.Errorf("unexpected error: %s", err.Error())
				}
				if err.Error() != test.Err {
					t.Errorf("unexpected error. Expected %s but got %s", test.Err, err.Error())
				}
				return
			}
			if test.WantErr {
				t.Errorf("expected error but nothing got")
			}
			if len(entity.Stars) > test.Data.Limit {
				t.Errorf("wrong number of stars per page. Expected %d but got %d", test.Data.Limit, len(entity.Stars))
			}
			if entity.TotalCount != len(stars) {
				t.Errorf("wrong number of stars. Expected %d but got %d", len(stars), entity.TotalCount)
			}
			for i, star := range stars {
				j := len(entity.Stars) - i - 1
				if entity.Stars[j].Name != star.Name {
					t.Errorf("wrong name. Expected %q but got %q", star.Name, entity.Stars[j].Name)
				}
				if entity.Stars[j].Sex != star.Sex {
					t.Errorf("wrong sex. Expected %q but got %q", star.Sex, entity.Stars[j].Sex)
				}
				if entity.Stars[j].BirthDate != star.BirthDate {
					t.Errorf("wrong birth date. Expected %q but got %q", star.BirthDate, entity.Stars[j].BirthDate)
				}
				if entity.Stars[j].CreatedAt != star.CreatedAt {
					t.Errorf("wrong created at. Expected %q but got %q", star.CreatedAt, entity.Stars[j].CreatedAt)
				}
				if entity.Stars[j].UpdatedAt != star.UpdatedAt {
					t.Errorf("wrong updated at. Expected %q but got %q", star.UpdatedAt, entity.Stars[j].UpdatedAt)
				}
				if entity.Stars[j].DeletedAt != nil {
					t.Error("wrong deleted at. Expected nil")
				}
			}
			for i := 0; i < len(entity.Stars)-1; i++ {
				if entity.Stars[i].ID < entity.Stars[i+1].ID {
					t.Errorf("wrong order. Expected %d to be before %d", i+1, i)
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
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

	var stars []Entity
	for _, data := range []CreateEntity{
		{
			Name:      "Ryan Gosling",
			Sex:       "male",
			BirthDate: time.Date(1980, time.November, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:      "Emma Stone",
			Sex:       "female",
			BirthDate: time.Date(1988, time.November, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:      "Leonardo DiCaprio",
			Sex:       "male",
			BirthDate: time.Date(1974, time.November, 11, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:      "Zendaya",
			Sex:       "female",
			BirthDate: time.Date(1996, time.September, 1, 0, 0, 0, 0, time.UTC),
		},
	} {
		entity, err := store.Create(ctx, data)
		if err != nil {
			t.Errorf("error with creating stars: %s", err.Error())
		}

		stars = append(stars, entity)
	}

	newName := "Ryan Gosling 2"
	newSex := "female"
	newBirthDate := time.Date(1980, time.December, 12, 0, 0, 0, 0, time.UTC)

	var id int
	for _, test := range []struct {
		Name    string
		Data    UpdateEntity
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful full update a star",
			Data: UpdateEntity{
				Name:      &newName,
				Sex:       &newSex,
				BirthDate: &newBirthDate,
			},
		},
		{
			Name: "Successful update name only",
			Data: UpdateEntity{
				Name: &newName,
			},
		},
		{
			Name: "Successful update sex only",
			Data: UpdateEntity{
				Sex: &newSex,
			},
		},
		{
			Name: "Successful update birth date only",
			Data: UpdateEntity{
				BirthDate: &newBirthDate,
			},
		},
		{
			Name:    "Update non-existent star",
			WantErr: true,
			Err:     pgx.ErrNoRows.Error(),
		},
	} {
		id++

		t.Run(test.Name, func(t *testing.T) {
			entity, err := store.Update(ctx, id, test.Data)
			if err != nil {
				if !test.WantErr {
					t.Errorf("unexpected error: %s", err.Error())
				}
				if err.Error() != test.Err {
					t.Errorf("unexpected error. Expected %s but got %s", test.Err, err.Error())
				}
				return
			}
			if test.WantErr {
				t.Errorf("expected error but nothing got")
			}
			if test.Data.Name != nil {
				if entity.Name != *test.Data.Name {
					t.Errorf("wrong name. Expected %q but got %q", *test.Data.Name, entity.Name)
				}
			} else {
				if entity.Name != stars[id-1].Name {
					t.Errorf("wrong name. Expected %q but got %q", stars[id-1].Name, entity.Name)
				}
			}
			if test.Data.Sex != nil {
				if entity.Sex != *test.Data.Sex {
					t.Errorf("wrong sex. Expected %q but got %q", *test.Data.Sex, entity.Sex)
				}
			} else {
				if entity.Sex != stars[id-1].Sex {
					t.Errorf("wrong sex. Expected %q but got %q", stars[id-1].Sex, entity.Sex)
				}
			}
			if test.Data.BirthDate != nil {
				if entity.BirthDate != *test.Data.BirthDate {
					t.Errorf("wrong birth date. Expected %q but got %q", *test.Data.BirthDate, entity.BirthDate)
				}
			} else {
				if entity.BirthDate != stars[id-1].BirthDate {
					t.Errorf("wrong birth date. Expected %q but got %q", stars[id-1].BirthDate, entity.BirthDate)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
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

	star, err := store.Create(ctx, CreateEntity{
		Name:      "Ryan Gosling",
		Sex:       "male",
		BirthDate: time.Date(1980, time.November, 12, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Errorf("error with creating test data: %s", err.Error())
	}

	for _, test := range []struct {
		Name    string
		ID      int
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful delete a star",
			ID:   star.ID,
		},
		{
			Name:    "Delete already deleted star",
			ID:      star.ID,
			WantErr: true,
			Err:     pgx.ErrNoRows.Error(),
		},
		{
			Name:    "Delete non-existent star",
			ID:      100,
			WantErr: true,
			Err:     pgx.ErrNoRows.Error(),
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			err := store.Delete(ctx, test.ID)
			if err != nil {
				if !test.WantErr {
					t.Errorf("unexpected error: %s", err.Error())
				}
				if err.Error() != test.Err {
					t.Errorf("unexpected error. Expected %s but got %s", test.Err, err.Error())
				}
				return
			}
			if test.WantErr {
				t.Errorf("expected error but nothing got")
			}

			entity, err := store.GetByID(ctx, test.ID)
			if err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}
			if entity.DeletedAt == nil {
				t.Errorf("deleted at should not be nil")
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

	star, err := store.Create(ctx, CreateEntity{
		Name:      "Ryan Gosling",
		Sex:       "male",
		BirthDate: time.Date(1980, time.November, 12, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Errorf("error with creating test data: %s", err.Error())
	}

	for _, test := range []struct {
		Name    string
		ID      int
		Exists  bool
		WantErr bool
		Err     string
	}{
		{
			Name:   "Check existent star",
			ID:     star.ID,
			Exists: true,
		},
		{
			Name: "Check non-existent star",
			ID:   100,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			exists, err := store.CheckExistence(ctx, test.ID)
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
