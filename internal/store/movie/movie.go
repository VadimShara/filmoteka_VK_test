package movie

import (
	"context"
	"fmt"
	"time"

	"vk-test-task/pkg/format"
	"vk-test-task/pkg/logger"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Store interface {
		Create(context.Context, CreateEntity) (Entity, error)
		GetByID(context.Context, int) (Entity, error)
		GetByStarID(context.Context, int) ([]Entity, error)
		GetAll(context.Context, GetAllParams) (EntityWithTotalCount, error)
		Update(context.Context, int, UpdateEntity) (Entity, error)
		Delete(context.Context, int) error
	}

	storeImpl struct {
		client      *pgxpool.Pool
		statBuilder sq.StatementBuilderType
	}

	GetAllParams struct {
		SearchTerm string
		SortBy     string
		SortOrder  string
		Limit      int
		Offset     int
	}

	CreateEntity struct {
		Title       string
		Description string
		ReleaseDate time.Time
		Rating      int
		StarsID     []int
	}

	UpdateEntity struct {
		Title       *string
		Description *string
		ReleaseDate *time.Time
		Rating      *int
		StarsID     []int
	}

	Entity struct {
		ID          int
		Title       string
		Description string
		ReleaseDate time.Time
		Rating      int
		CreatedAt   time.Time
		UpdatedAt   time.Time
		DeletedAt   *time.Time
	}

	EntityWithTotalCount struct {
		Movies     []Entity
		TotalCount int
	}
)

func New(client *pgxpool.Pool) Store {
	return &storeImpl{
		client:      client,
		statBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *storeImpl) Create(ctx context.Context, entity CreateEntity) (Entity, error) {
	var newMovie Entity

	tx, err := s.client.Begin(ctx)
	if err != nil {
		return Entity{}, err
	}

	err = tx.QueryRow(
		ctx,
		`
			INSERT INTO movies (title, description, release_date, rating, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, title, description, release_date, rating, created_at, updated_at
		`,
		entity.Title,
		entity.Description,
		entity.ReleaseDate,
		entity.Rating,
		format.TimeNow(),
		format.TimeNow(),
	).Scan(&newMovie.ID,
		&newMovie.Title,
		&newMovie.Description,
		&newMovie.ReleaseDate,
		&newMovie.Rating,
		&newMovie.CreatedAt,
		&newMovie.UpdatedAt)
	if err != nil {
		logger.Log.Error("create new movie",
			"error", err.Error())
		if txErr := tx.Rollback(ctx); txErr != nil {
			logger.Log.Error("rollback",
				"error", txErr.Error())
			return Entity{}, txErr
		}
		return Entity{}, err
	}

	if err = s.addStarsToMovieInTx(ctx, newMovie.ID, entity.StarsID, tx); err != nil {
		logger.Log.Error("adding stars to movie",
			"error", err.Error())
		if txErr := tx.Rollback(ctx); txErr != nil {
			logger.Log.Error("rollback",
				"error", txErr.Error())
			return Entity{}, txErr
		}
		return Entity{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.Log.Error("committing create transaction",
			"error", err.Error())
	}

	return newMovie, err
}

func (s *storeImpl) GetByID(ctx context.Context, id int) (Entity, error) {
	var movie Entity

	err := s.client.QueryRow(
		ctx,
		`
			SELECT id, title, description, release_date, rating, created_at, updated_at, deleted_at
			FROM movies
			WHERE id = $1
		`,
		id,
	).Scan(&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.ReleaseDate,
		&movie.Rating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
		&movie.DeletedAt)
	if err != nil {
		logger.Log.Error("get movie by id",
			"error", err.Error())
	}

	return movie, err
}

func (s *storeImpl) GetByStarID(ctx context.Context, starID int) ([]Entity, error) {
	var movies []Entity

	rows, err := s.client.Query(
		ctx,
		`
			SELECT m.id, m.title, m.description, m.release_date, m.rating, m.created_at, m.updated_at, m.deleted_at
			FROM movies m
			LEFT JOIN movie_stars ms ON m.id = ms.movie_id
			LEFT JOIN stars s ON ms.star_id = s.id
			WHERE s.id = $1
		`,
		starID,
	)
	if err != nil {
		logger.Log.Error("get movies by star id",
			"error", err.Error())
		return []Entity{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie Entity

		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.ReleaseDate,
			&movie.Rating,
			&movie.CreatedAt,
			&movie.UpdatedAt,
			&movie.DeletedAt,
		); err != nil {
			logger.Log.Error("scan star",
				"error", err.Error())
			return []Entity{}, err
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

func (s *storeImpl) GetAll(ctx context.Context, params GetAllParams) (EntityWithTotalCount, error) {
	var movies EntityWithTotalCount

	selectQuery := s.statBuilder.
		Select("m.id", "m.title", "m.description", "m.release_date", "m.rating", "m.created_at", "m.updated_at", "m.deleted_at", "COUNT(*) OVER() AS total").
		From("movies m").
		Where("m.deleted_at IS NULL").
		Limit(uint64(params.Limit)).
		Offset(uint64(params.Offset))

	if params.SearchTerm != "" {
		selectQuery = selectQuery.
			Join("movie_stars ms ON m.id = ms.movie_id").
			Join("stars s ON ms.star_id = s.id").
			Where(sq.Or{
				sq.ILike{"m.title": "%" + params.SearchTerm + "%"},
				sq.ILike{"s.name": "%" + params.SearchTerm + "%"},
			})
	}

	if params.SortBy != "" {
		selectQuery = selectQuery.
			OrderBy(fmt.Sprintf("%s %s", params.SortBy, params.SortOrder))
	} else {
		selectQuery = selectQuery.
			OrderBy("m.rating DESC")
	}

	sqlQuery, args, err := selectQuery.ToSql()
	if err != nil {
		logger.Log.Error("generate select query",
			"error", err.Error())
		return EntityWithTotalCount{}, err
	}

	rows, err := s.client.Query(ctx, sqlQuery, args...)
	if err != nil {
		logger.Log.Error("get movies",
			"error", err.Error())
		return EntityWithTotalCount{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie Entity
		var total int

		if err := rows.Scan(&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.ReleaseDate,
			&movie.Rating,
			&movie.CreatedAt,
			&movie.UpdatedAt,
			&movie.DeletedAt,
			&total,
		); err != nil {
			logger.Log.Error("scan star",
				"error", err.Error())
			return EntityWithTotalCount{}, err
		}

		movies.Movies = append(movies.Movies, movie)
		movies.TotalCount = total
	}

	return movies, nil
}

func (s *storeImpl) Update(ctx context.Context, id int, entity UpdateEntity) (Entity, error) {
	item, err := s.GetByID(ctx, id)
	if err != nil {
		logger.Log.Error("get star by id",
			"error", err.Error())
		return Entity{}, err
	}
	if entity.Title != nil {
		item.Title = *entity.Title
	}
	if entity.Description != nil {
		item.Description = *entity.Description
	}
	if entity.ReleaseDate != nil {
		item.ReleaseDate = *entity.ReleaseDate
	}
	if entity.Rating != nil {
		item.Rating = *entity.Rating
	}
	oldStarsID, err := s.getStarsIDForMovie(ctx, id)
	if err != nil {
		return Entity{}, err
	}

	tx, err := s.client.Begin(ctx)
	if err != nil {
		return Entity{}, err
	}

	var updatedStar Entity
	if err = tx.QueryRow(ctx,
		`
			UPDATE movies
			SET title = $1, description = $2, release_date = $3, rating = $4, updated_at = $5
			WHERE id = $6 AND deleted_at IS NULL
			RETURNING id, title, description, release_date, rating, created_at, updated_at
		`,
		item.Title,
		item.Description,
		item.ReleaseDate,
		item.Rating,
		format.TimeNow(),
		id,
	).Scan(&updatedStar.ID,
		&updatedStar.Title,
		&updatedStar.Description,
		&updatedStar.ReleaseDate,
		&updatedStar.Rating,
		&updatedStar.CreatedAt,
		&updatedStar.UpdatedAt,
	); err != nil {
		logger.Log.Error("update star",
			"error", err.Error())
		if txErr := tx.Rollback(ctx); txErr != nil {
			logger.Log.Error("rollback",
				"error", txErr.Error())
			return Entity{}, txErr
		}
		return Entity{}, err
	}

	// We believe that a movie cannot exist without actors
	if len(entity.StarsID) != 0 {
		if err := s.updateStarsForMovieInTx(ctx, id, oldStarsID, entity.StarsID, tx); err != nil {
			logger.Log.Error("update stars for movie",
				"error", err.Error())
			if txErr := tx.Rollback(ctx); txErr != nil {
				logger.Log.Error("rollback",
					"error", txErr.Error())
				return Entity{}, txErr
			}
			return Entity{}, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.Log.Error("committing update transaction",
			"error", err.Error())
	}

	return updatedStar, err
}

func (s *storeImpl) Delete(ctx context.Context, id int) error {
	n, err := s.client.Exec(
		ctx,
		`
			UPDATE movies
			SET updated_at = $1, deleted_at = $2
			WHERE id = $3 AND deleted_at IS NULL
		`,
		format.TimeNow(),
		format.TimeNow(),
		id,
	)
	if err != nil {
		logger.Log.Error("delete movie",
			"error", err.Error())
		return err
	}

	if n.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (s *storeImpl) addStarsToMovieInTx(ctx context.Context, movieID int, starsID []int, tx pgx.Tx) error {
	for _, starID := range starsID {
		_, err := tx.Exec(
			ctx,
			`
                INSERT INTO movie_stars (movie_id, star_id)
                VALUES ($1, $2)
            `,
			movieID,
			starID,
		)
		if err != nil {
			logger.Log.Error("add star to movie",
				"error", err.Error())
			return err
		}
	}

	return nil
}

func (s *storeImpl) deleteStarsFromMovieInTx(ctx context.Context, movieID int, starsID []int, tx pgx.Tx) error {
	for _, starID := range starsID {
		_, err := tx.Exec(
			ctx,
			`
                DELETE FROM movie_stars
                WHERE movie_id = $1 AND star_id = $2
            `,
			movieID,
			starID,
		)
		if err != nil {
			logger.Log.Error("delete star from movie",
				"error", err.Error())
			return err
		}
	}

	return nil
}

func (s *storeImpl) updateStarsForMovieInTx(ctx context.Context, movieID int, oldStarsID, newStarsID []int, tx pgx.Tx) error {
	if err := s.deleteStarsFromMovieInTx(ctx, movieID, oldStarsID, tx); err != nil {
		return err
	}

	return s.addStarsToMovieInTx(ctx, movieID, newStarsID, tx)
}

func (s *storeImpl) getStarsIDForMovie(ctx context.Context, movieID int) ([]int, error) {
	var starsID []int

	rows, err := s.client.Query(
		ctx,
		`
            SELECT star_id
            FROM movie_stars
            WHERE movie_id = $1
        `,
		movieID,
	)
	if err != nil {
		logger.Log.Error("get stars for movie",
			"error", err.Error())
		return []int{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var starID int

		if err := rows.Scan(&starID); err != nil {
			logger.Log.Error("scan star",
				"error", err.Error())
			return []int{}, err
		}

		starsID = append(starsID, starID)
	}

	return starsID, nil
}
