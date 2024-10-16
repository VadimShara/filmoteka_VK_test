package movie

import (
	"context"
	"math/rand"
	"strings"
	"testing"
	"time"

	"vk-test-task/internal/tests"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type star struct {
	id   int
	name string
}

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

	starsID, err := addExistingStars(ctx, postgresClient)
	if err != nil {
		t.Errorf("error with adding existing data: %s", err.Error())
	}

	store := New(postgresClient)

	for _, test := range []struct {
		Name    string
		Data    CreateEntity
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful creating a movie",
			Data: CreateEntity{
				Title: "Drive",
				Description: `I'm giving you a night call to tell you how I feel (We'll go all, all, all night long)
				I want to drive you through the night, down the hills (We'll go all, all, all night long)
				I'm gonna tell you something you don't want to hear (We'll go all, all, all night long)
				I'm gonna show you where it's dark, but have no fear (We'll go all, all, all night long)`,
				ReleaseDate: time.Date(2011, time.November, 3, 0, 0, 0, 0, time.UTC),
				Rating:      10,
				StarsID:     starsID,
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
			if entity.Title != test.Data.Title {
				t.Errorf("wrong title. Expected %q but got %q", test.Data.Title, entity.Title)
			}
			if entity.Description != test.Data.Description {
				t.Errorf("wrong description. Expected %q but got %q", test.Data.Description, entity.Description)
			}
			if entity.ReleaseDate != test.Data.ReleaseDate {
				t.Errorf("wrong release date. Expected %q but got %q", test.Data.ReleaseDate, entity.ReleaseDate)
			}
			if entity.Rating != test.Data.Rating {
				t.Errorf("wrong rating. Expected %d but got %d", test.Data.Rating, entity.Rating)
			}
			if entity.DeletedAt != nil {
				t.Error("deleted at should be nil")
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

	starsID, err := addExistingStars(ctx, postgresClient)
	if err != nil {
		t.Errorf("error with adding existing data: %s", err.Error())
	}

	store := New(postgresClient)

	data := CreateEntity{
		Title:       "Drive",
		Description: "I'm giving you a night call to tell you how I feel (We'll go all, all, all night long)",
		ReleaseDate: time.Date(2011, time.November, 3, 0, 0, 0, 0, time.UTC),
		Rating:      10,
		StarsID:     starsID,
	}

	movie, err := store.Create(ctx, data)
	if err != nil {
		t.Errorf("error with creating test movie: %s", err.Error())
	}

	for _, test := range []struct {
		Name    string
		ID      int
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful get a movie",
			ID:   movie.ID,
		},
		{
			Name:    "Get non-existent movie",
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
			if entity.Title != movie.Title {
				t.Errorf("wrong title. Expected %q but got %q", movie.Title, entity.Title)
			}
			if entity.Description != movie.Description {
				t.Errorf("wrong description. Expected %q but got %q", movie.Description, entity.Description)
			}
			if entity.ReleaseDate != movie.ReleaseDate {
				t.Errorf("wrong release date. Expected %q but got %q", movie.ReleaseDate, entity.ReleaseDate)
			}
			if entity.Rating != movie.Rating {
				t.Errorf("wrong rating. Expected %d but got %d", movie.Rating, entity.Rating)
			}
			if entity.CreatedAt != movie.CreatedAt {
				t.Errorf("wrong created at. Expected %q but got %q", movie.CreatedAt, entity.CreatedAt)
			}
			if entity.UpdatedAt != movie.UpdatedAt {
				t.Errorf("wrong updated at. Expected %q but got %q", movie.UpdatedAt, entity.UpdatedAt)
			}
			if entity.DeletedAt != nil {
				t.Error("wrong deleted at. Expected nil")
			}
		})
	}
}

func TestGetByStarID(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresContainer, databaseURL := tests.CreateFilmotekaTestPostgresContainer(ctx)
	defer tests.TerminateFilmotekaTestContainer(ctx, postgresContainer)

	postgresClient, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Errorf("error with starting postgres client: %s", err.Error())
	}
	defer postgresClient.Close()

	_, err = addExistingStars(ctx, postgresClient)
	if err != nil {
		t.Errorf("error with adding existing data: %s", err.Error())
	}

	store := New(postgresClient)

	cases := []CreateEntity{
		{
			Title:       "Drive",
			Description: "I'm giving you a night call to tell you how I feel (We'll go all, all, all night long)",
			ReleaseDate: time.Date(2011, time.November, 3, 0, 0, 0, 0, time.UTC),
			Rating:      10,
			StarsID:     []int{1},
		},
		{
			Title:       "Drive 2",
			Description: "I want to drive you through the night, down the hills (We'll go all, all, all night long)",
			ReleaseDate: time.Date(2012, time.November, 3, 0, 0, 0, 0, time.UTC),
			Rating:      10,
			StarsID:     []int{1},
		},
	}

	var movies []Entity
	for _, c := range cases {
		movie, err := store.Create(ctx, c)
		if err != nil {
			t.Errorf("error with creating test movie: %s", err.Error())
		}

		movies = append(movies, movie)
	}

	for _, test := range []struct {
		Name    string
		ID      int
		Count   int
		WantErr bool
		Err     string
	}{
		{
			Name:  "Successful get a movies with Ryan Gosling",
			ID:    1,
			Count: 2,
		},
		{
			Name:  "Get movies for Leonardo DiCaprio",
			ID:    1,
			Count: 0,
		},
		{
			Name:  "Get movies for non-existent star",
			ID:    100,
			Count: 0,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			entity, err := store.GetByStarID(ctx, test.ID)
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
			for i, e := range entity {
				if e.Title != movies[i].Title {
					t.Errorf("wrong title. Expected %q but got %q", movies[i].Title, e.Title)
				}
				if e.Description != movies[i].Description {
					t.Errorf("wrong description. Expected %q but got %q", movies[i].Description, e.Description)
				}
				if e.ReleaseDate != movies[i].ReleaseDate {
					t.Errorf("wrong release date. Expected %q but got %q", movies[i].ReleaseDate, e.ReleaseDate)
				}
				if e.Rating != movies[i].Rating {
					t.Errorf("wrong rating. Expected %d but got %d", movies[i].Rating, e.Rating)
				}
				if e.CreatedAt != movies[i].CreatedAt {
					t.Errorf("wrong created at. Expected %q but got %q", movies[i].CreatedAt, e.CreatedAt)
				}
				if e.UpdatedAt != movies[i].UpdatedAt {
					t.Errorf("wrong updated at. Expected %q but got %q", movies[i].UpdatedAt, e.UpdatedAt)
				}
				if e.DeletedAt != nil {
					t.Error("wrong deleted at. Expected nil")
				}
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

	starsID, err := addExistingStars(ctx, postgresClient)
	if err != nil {
		t.Errorf("error with adding existing data: %s", err.Error())
	}

	store := New(postgresClient)

	var movies []Entity
	for _, data := range []CreateEntity{
		{
			Title:       "Drive",
			Description: "I'm giving you a night call to tell you how I feel (We'll go all, all, all night long)",
			ReleaseDate: time.Date(2011, time.November, 3, 0, 0, 0, 0, time.UTC),
			Rating:      10,
			StarsID:     []int{starsID[0]},
		},
		{
			Title:       "Oppenheimer",
			Description: "Boom",
			ReleaseDate: time.Date(2023, time.July, 19, 0, 0, 0, 0, time.UTC),
			Rating:      9,
			StarsID:     []int{starsID[rand.Intn(len(starsID))]},
		},
		{
			Title:       "Barbie",
			Description: "Life in plastic, it's fantastic!",
			ReleaseDate: time.Date(2011, time.November, 3, 0, 0, 0, 0, time.UTC),
			Rating:      9,
			StarsID:     []int{starsID[rand.Intn(len(starsID))]},
		},
		{
			Title:       "Interstellar",
			Description: "1 hour here is 7 years on earth",
			ReleaseDate: time.Date(2014, time.November, 6, 0, 0, 0, 0, time.UTC),
			Rating:      9,
			StarsID:     []int{starsID[rand.Intn(len(starsID))]},
		},
		{
			Title:       "Barbieneimer",
			Description: "Another Boom",
			ReleaseDate: time.Date(2023, time.June, 18, 0, 0, 0, 0, time.UTC),
			Rating:      10,
			StarsID:     []int{starsID[rand.Intn(len(starsID))]},
		},
	} {
		entity, err := store.Create(ctx, data)
		if err != nil {
			t.Errorf("error with creating movie: %s", err.Error())
		}

		movies = append(movies, entity)
	}

	for _, test := range []struct {
		Name          string
		Data          GetAllParams
		WantErr       bool
		IsSorted      bool
		SearchByTitle bool
		SearchByStar  bool
		Err           string
	}{
		{
			Name: "Successful get all movies paginated",
			Data: GetAllParams{
				Limit: len(movies),
			},
		},
		{
			Name: "Successful search film by title",
			Data: GetAllParams{
				SearchTerm: "Barbie",
			},
			SearchByTitle: true,
		},
		{
			Name: "Successful search film by star name",
			Data: GetAllParams{
				SearchTerm: "Gosling",
			},
			SearchByStar: true,
		},
		{
			Name: "Successful sort movies by title asc",
			Data: GetAllParams{
				SortBy:    "m.title",
				SortOrder: "ASC",
			},
			IsSorted: true,
		},
		{
			Name: "Successful sort movies by title desc",
			Data: GetAllParams{
				SortBy:    "m.title",
				SortOrder: "DESC",
			},
			IsSorted: true,
		},
		{
			Name: "Successful sort movies by rating asc",
			Data: GetAllParams{
				SortBy:    "m.rating",
				SortOrder: "ASC",
			},
			IsSorted: true,
		},
		{
			Name: "Successful sort movies by rating desc",
			Data: GetAllParams{
				SortBy:    "m.rating",
				SortOrder: "DESC",
			},
			IsSorted: true,
		},
		{
			Name: "Successful sort movies by release date asc",
			Data: GetAllParams{
				SortBy:    "m.release_date",
				SortOrder: "ASC",
			},
			IsSorted: true,
		},
		{
			Name: "Successful sort movies by release_date desc",
			Data: GetAllParams{
				SortBy:    "m.release_date",
				SortOrder: "DESC",
			},
			IsSorted: true,
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
			if test.SearchByTitle {
				for _, movie := range entity.Movies {
					if !strings.Contains(movie.Title, test.Data.SearchTerm) {
						t.Errorf("wrong title. Expected %q to be substring of %q", "Barbie", movie.Title)
					}
				}
			}
			if test.SearchByStar {
				for _, movie := range entity.Movies {
					stars, err := getStarsNameByMovieID(ctx, movie.ID, postgresClient)
					if err != nil {
						t.Errorf("error with getting stars name: %s", err.Error())
					}

					var nameFound bool
					for _, star := range stars {
						if strings.Contains(star.name, test.Data.SearchTerm) {
							nameFound = true
							break
						}
					}
					if !nameFound {
						t.Error("star name not found")
					}
				}
			}
			if test.IsSorted {
				switch test.Data.SortBy {
				case "m.title":
					if test.Data.SortOrder == "ASC" {
						for i := 0; i < len(entity.Movies)-1; i++ {
							if entity.Movies[i].Title > entity.Movies[i+1].Title {
								t.Errorf("wrong sort order. Expected %q but got %q", entity.Movies[i].Title, entity.Movies[i+1].Title)
							}
						}
					} else {
						for i := 0; i < len(entity.Movies)-1; i++ {
							if entity.Movies[i].Title < entity.Movies[i+1].Title {
								t.Errorf("wrong sort order. Expected %q but got %q", entity.Movies[i].Title, entity.Movies[i+1].Title)
							}
						}
					}
				case "m.rating":
					if test.Data.SortOrder == "ASC" {
						for i := 0; i < len(entity.Movies)-1; i++ {
							if entity.Movies[i].Rating > entity.Movies[i+1].Rating {
								t.Errorf("wrong sort order. Expected %d but got %d", entity.Movies[i].Rating, entity.Movies[i+1].Rating)
							}
						}
					} else {
						for i := 0; i < len(entity.Movies)-1; i++ {
							if entity.Movies[i].Rating < entity.Movies[i+1].Rating {
								t.Errorf("wrong sort order. Expected %d but got %d", entity.Movies[i].Rating, entity.Movies[i+1].Rating)
							}
						}
					}
				case "m.release_date":
					if test.Data.SortOrder == "ASC" {
						for i := 0; i < len(entity.Movies)-1; i++ {
							if entity.Movies[i].ReleaseDate.After(entity.Movies[i+1].ReleaseDate) {
								t.Errorf("wrong sort order. Expected %s but got %s", entity.Movies[i].ReleaseDate, entity.Movies[i+1].ReleaseDate)
							}
						}
					} else {
						for i := 0; i < len(entity.Movies)-1; i++ {
							if entity.Movies[i].ReleaseDate.Before(entity.Movies[i+1].ReleaseDate) {
								t.Errorf("wrong sort order. Expected %s but got %s", entity.Movies[i].ReleaseDate, entity.Movies[i+1].ReleaseDate)
							}
						}
					}
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

	starsID, err := addExistingStars(ctx, postgresClient)
	if err != nil {
		t.Errorf("error with adding existing data: %s", err.Error())
	}

	store := New(postgresClient)

	var movies []Entity
	for _, data := range []CreateEntity{
		{
			Title:       "Drive",
			Description: "I'm giving you a night call to tell you how I feel (We'll go all, all, all night long)",
			ReleaseDate: time.Date(2011, time.November, 3, 0, 0, 0, 0, time.UTC),
			Rating:      10,
			StarsID:     []int{starsID[0], starsID[2]},
		},
		{
			Title:       "Oppenheimer",
			Description: "Boom",
			ReleaseDate: time.Date(2023, time.July, 19, 0, 0, 0, 0, time.UTC),
			Rating:      9,
			StarsID:     []int{starsID[rand.Intn(len(starsID))]},
		},
		{
			Title:       "Barbie",
			Description: "Life in plastic, it's fantastic!",
			ReleaseDate: time.Date(2011, time.November, 3, 0, 0, 0, 0, time.UTC),
			Rating:      9,
			StarsID:     []int{starsID[rand.Intn(len(starsID))]},
		},
		{
			Title:       "Interstellar",
			Description: "1 hour here is 7 years on earth.",
			ReleaseDate: time.Date(2014, time.November, 6, 0, 0, 0, 0, time.UTC),
			Rating:      9,
			StarsID:     []int{starsID[rand.Intn(len(starsID))]},
		},
		{
			Title:       "Barbieneimer",
			Description: "Another Boom",
			ReleaseDate: time.Date(2023, time.June, 18, 0, 0, 0, 0, time.UTC),
			Rating:      10,
			StarsID:     []int{starsID[rand.Intn(len(starsID))]},
		},
		{
			Title:       "Star Wars: Episode II - Attack of the Clones",
			Description: "I don't like sand. It's coarse and rough and irritating and it gets everywhere.",
			ReleaseDate: time.Date(2002, time.May, 16, 0, 0, 0, 0, time.UTC),
			Rating:      1,
			StarsID:     []int{starsID[3]},
		},
	} {
		entity, err := store.Create(ctx, data)
		if err != nil {
			t.Errorf("error with creating movie: %s", err.Error())
		}

		movies = append(movies, entity)
	}

	newTitle := "Drive 2"
	newDescription := "Wow"
	newReleaseDate := time.Date(2024, time.December, 12, 0, 0, 0, 0, time.UTC)
	newRating := 5
	newStarsID := []int{starsID[0], starsID[1]}

	var id int
	for _, test := range []struct {
		Name    string
		Data    UpdateEntity
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful full update a movie",
			Data: UpdateEntity{
				Title:       &newTitle,
				Description: &newDescription,
				ReleaseDate: &newReleaseDate,
				Rating:      &newRating,
				StarsID:     newStarsID,
			},
		},
		{
			Name: "Successful update title only",
			Data: UpdateEntity{
				Title: &newTitle,
			},
		},
		{
			Name: "Successful update description only",
			Data: UpdateEntity{
				Description: &newDescription,
			},
		},
		{
			Name: "Successful update release date only",
			Data: UpdateEntity{
				ReleaseDate: &newReleaseDate,
			},
		},
		{
			Name: "Successful update rating only",
			Data: UpdateEntity{
				Rating: &newRating,
			},
		},
		{
			Name: "Successful update stars only",
			Data: UpdateEntity{
				StarsID: newStarsID,
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
			if test.Data.Title != nil {
				if entity.Title != *test.Data.Title {
					t.Errorf("wrong title. Expected %q but got %q", *test.Data.Title, entity.Title)
				}
			} else {
				if entity.Title != movies[id-1].Title {
					t.Errorf("wrong title. Expected %q but got %q", movies[id-1].Title, entity.Title)
				}
			}
			if test.Data.Description != nil {
				if entity.Description != *test.Data.Description {
					t.Errorf("wrong description. Expected %q but got %q", *test.Data.Description, entity.Description)
				}
			} else {
				if entity.Description != movies[id-1].Description {
					t.Errorf("wrong description. Expected %q but got %q", movies[id-1].Description, entity.Description)
				}
			}
			if test.Data.ReleaseDate != nil {
				if entity.ReleaseDate != *test.Data.ReleaseDate {
					t.Errorf("wrong release date. Expected %q but got %q", *test.Data.ReleaseDate, entity.ReleaseDate)
				}
			} else {
				if entity.ReleaseDate != movies[id-1].ReleaseDate {
					t.Errorf("wrong release date. Expected %q but got %q", movies[id-1].ReleaseDate, entity.ReleaseDate)
				}
			}
			if test.Data.Rating != nil {
				if entity.Rating != *test.Data.Rating {
					t.Errorf("wrong rating. Expected %q but got %q", *test.Data.Rating, entity.Rating)
				}
			} else {
				if entity.Rating != movies[id-1].Rating {
					t.Errorf("wrong rating. Expected %q but got %q", movies[id-1].Rating, entity.Rating)
				}
			}
			if test.Data.StarsID != nil {
				stars, err := getStarsNameByMovieID(ctx, id, postgresClient)
				if err != nil {
					t.Errorf("getting stars list for movie: %s", err.Error())
				}

				for _, starID := range test.Data.StarsID {
					found := false
					for _, star := range stars {
						if starID == star.id {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("star %d not found", starID)
					}
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

	starsID, err := addExistingStars(ctx, postgresClient)
	if err != nil {
		t.Errorf("error with adding existing data: %s", err.Error())
	}

	store := New(postgresClient)

	data := CreateEntity{
		Title: "Drive",
		Description: `I'm giving you a night call to tell you how I feel (We'll go all, all, all night long)
				I want to drive you through the night, down the hills (We'll go all, all, all night long)
				I'm gonna tell you something you don't want to hear (We'll go all, all, all night long)
				I'm gonna show you where it's dark, but have no fear (We'll go all, all, all night long)`,
		ReleaseDate: time.Date(2011, time.November, 3, 0, 0, 0, 0, time.UTC),
		Rating:      10,
		StarsID:     starsID,
	}

	movie, err := store.Create(ctx, data)
	if err != nil {
		t.Errorf("error with creating test movie: %s", err.Error())
	}

	for _, test := range []struct {
		Name    string
		ID      int
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful delete a movie",
			ID:   movie.ID,
		},
		{
			Name:    "Delete already deleted movie",
			ID:      movie.ID,
			WantErr: true,
			Err:     pgx.ErrNoRows.Error(),
		},
		{
			Name:    "Delete non-existent movie",
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

func addExistingStars(ctx context.Context, postgresClient *pgxpool.Pool) ([]int, error) {
	var starsID []int

	for _, data := range []struct {
		Name      string
		Sex       string
		BirthDate time.Time
	}{
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
		var id int
		if err := postgresClient.QueryRow(ctx,
			`
				INSERT INTO stars (name, sex, birth_date, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id
			`,
			data.Name,
			data.Sex,
			data.BirthDate,
			time.Now(),
			time.Now(),
		).Scan(&id); err != nil {
			return []int{}, err
		}

		starsID = append(starsID, id)
	}

	return starsID, nil
}

func getStarsNameByMovieID(ctx context.Context, movieID int, postgresClient *pgxpool.Pool) ([]star, error) {
	var stars []star

	rows, err := postgresClient.Query(
		ctx,
		`
			SELECT s.id, s.name
			FROM movies m 
			JOIN movie_stars ms ON m.id = ms.movie_id
			JOIN stars s ON ms.star_id = s.id
			WHERE m.id = $1
		`,
		movieID,
	)
	if err != nil {
		return []star{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var starID int
		var starName string

		if err := rows.Scan(&starID, &starName); err != nil {
			return []star{}, err
		}

		stars = append(stars, star{id: starID, name: starName})
	}

	return stars, nil
}
