package star

import (
	"context"
	"time"

	"vk-test-task/pkg/format"
	"vk-test-task/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Store interface {
		Create(context.Context, CreateEntity) (Entity, error)
		GetByID(context.Context, int) (Entity, error)
		GetAll(context.Context, GetAllParams) (EntityWithTotalCount, error)
		Update(context.Context, int, UpdateEntity) (Entity, error)
		Delete(context.Context, int) error
		CheckExistence(context.Context, int) (bool, error)
	}

	storeImpl struct {
		client *pgxpool.Pool
	}

	GetAllParams struct {
		Limit  int
		Offset int
	}

	CreateEntity struct {
		Name      string
		Sex       string
		BirthDate time.Time
	}

	UpdateEntity struct {
		Name      *string
		Sex       *string
		BirthDate *time.Time
	}

	Entity struct {
		ID        int
		Name      string
		Sex       string
		BirthDate time.Time
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}

	EntityWithTotalCount struct {
		Stars      []Entity
		TotalCount int
	}
)

func New(client *pgxpool.Pool) Store {
	return &storeImpl{
		client: client,
	}
}

func (s *storeImpl) Create(ctx context.Context, entity CreateEntity) (Entity, error) {
	var newStar Entity

	err := s.client.QueryRow(
		ctx,
		`
			INSERT INTO stars (name, sex, birth_date, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, name, sex, birth_date, created_at, updated_at
		`,
		entity.Name,
		entity.Sex,
		entity.BirthDate,
		format.TimeNow(),
		format.TimeNow(),
	).Scan(&newStar.ID,
		&newStar.Name,
		&newStar.Sex,
		&newStar.BirthDate,
		&newStar.CreatedAt,
		&newStar.UpdatedAt)
	if err != nil {
		logger.Log.Error("create new star",
			"error", err.Error())
	}

	return newStar, err
}

func (s *storeImpl) GetByID(ctx context.Context, id int) (Entity, error) {
	var star Entity

	err := s.client.QueryRow(
		ctx,
		`
			SELECT id, name, sex, birth_date, created_at, updated_at, deleted_at
			FROM stars
			WHERE id = $1
		`,
		id,
	).Scan(&star.ID,
		&star.Name,
		&star.Sex,
		&star.BirthDate,
		&star.CreatedAt,
		&star.UpdatedAt,
		&star.DeletedAt)
	if err != nil {
		logger.Log.Error("get star by id",
			"error", err.Error())
	}

	return star, err
}

func (s *storeImpl) GetAll(ctx context.Context, params GetAllParams) (EntityWithTotalCount, error) {
	var entities EntityWithTotalCount

	rows, err := s.client.Query(
		ctx,
		`
			SELECT id, name, sex, birth_date, created_at, updated_at, deleted_at, COUNT(*) OVER() AS total
			FROM stars
			WHERE deleted_at IS NULL
			ORDER BY id DESC
			LIMIT $1
			OFFSET $2
		`,
		params.Limit,
		params.Offset,
	)
	if err != nil {
		logger.Log.Error("get stars",
			"error", err.Error())
		return EntityWithTotalCount{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var entity Entity
		var total int

		if err := rows.Scan(
			&entity.ID,
			&entity.Name,
			&entity.Sex,
			&entity.BirthDate,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&total,
		); err != nil {
			logger.Log.Error("scan star",
				"error", err.Error())
			return EntityWithTotalCount{}, err
		}

		entities.Stars = append(entities.Stars, entity)
		entities.TotalCount = total
	}

	return entities, nil
}

func (s *storeImpl) Update(ctx context.Context, id int, entity UpdateEntity) (Entity, error) {
	item, err := s.GetByID(ctx, id)
	if err != nil {
		logger.Log.Error("get star by id",
			"error", err.Error())
		return Entity{}, err
	}
	if entity.Name != nil {
		item.Name = *entity.Name
	}
	if entity.Sex != nil {
		item.Sex = *entity.Sex
	}
	if entity.BirthDate != nil {
		item.BirthDate = *entity.BirthDate
	}

	var updatedStar Entity
	err = s.client.QueryRow(ctx,
		`
			UPDATE stars
			SET name = $1, sex = $2, birth_date = $3, updated_at = $4
			WHERE id = $5 AND deleted_at IS NULL
			RETURNING id, name, sex, birth_date, created_at, updated_at, deleted_at
		`,
		item.Name,
		item.Sex,
		item.BirthDate,
		format.TimeNow(),
		id,
	).Scan(&updatedStar.ID,
		&updatedStar.Name,
		&updatedStar.Sex,
		&updatedStar.BirthDate,
		&updatedStar.CreatedAt,
		&updatedStar.UpdatedAt,
		&updatedStar.DeletedAt)
	if err != nil {
		logger.Log.Error("update star",
			"error", err.Error())
	}

	return updatedStar, err
}

func (s *storeImpl) Delete(ctx context.Context, id int) error {
	n, err := s.client.Exec(
		ctx,
		`
			UPDATE stars
			SET updated_at = $1, deleted_at = $2
			WHERE id = $3 AND deleted_at IS NULL
		`,
		format.TimeNow(),
		format.TimeNow(),
		id,
	)
	if err != nil {
		logger.Log.Error("delete star",
			"error", err.Error())
		return err
	}

	if n.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (s *storeImpl) CheckExistence(ctx context.Context, id int) (bool, error) {
	var exists bool

	err := s.client.QueryRow(
		ctx,
		`
            SELECT EXISTS (
                SELECT 1
                FROM stars
                WHERE id = $1
            )
        `,
		id,
	).Scan(&exists)
	if err != nil {
		logger.Log.Error("check star existence",
			"error", err.Error())
	}

	return exists, err
}
