package tests

import (
	"context"
	"fmt"
	"os"
	"time"

	"vk-test-task/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	PostgresFilmotekaTestHost     = "localhost"
	postgresFilmotekaTestPort     = "5432"
	PostgresFilmotekaTestUser     = "postgres"
	PostgresFilmotekaTestPassword = "postgres"
	PostgresFilmotekaTestDB       = "filmoteka"
)

func CreateFilmotekaTestPostgresContainer(ctx context.Context) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15.3-bullseye",
		ExposedPorts: []string{postgresFilmotekaTestPort + "/tcp"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		Env: map[string]string{
			"POSTGRES_USER":     PostgresFilmotekaTestUser,
			"POSTGRES_PASSWORD": PostgresFilmotekaTestPassword,
			"POSTGRES_DB":       PostgresFilmotekaTestDB,
		},
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic("starting postgres container: " + err.Error())
	}

	mappedPort, err := postgresContainer.MappedPort(ctx, postgresFilmotekaTestPort)
	if err != nil {
		panic("getting mapped port: " + err.Error())
	}

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		PostgresFilmotekaTestUser,
		PostgresFilmotekaTestPassword,
		PostgresFilmotekaTestHost+":"+mappedPort.Port(),
		PostgresFilmotekaTestDB,
	)

	// Wait for the container to be ready
	time.Sleep(500 * time.Millisecond)
	if err = runMigrations(ctx, databaseURL); err != nil {
		panic("run migrations: " + err.Error())
	}

	if err := logger.SetupLogger("local"); err != nil {
		panic("setting up logger: " + err.Error())
	}

	return postgresContainer, databaseURL
}

func TerminateFilmotekaTestContainer(ctx context.Context, container testcontainers.Container) {
	if err := container.Terminate(ctx); err != nil {
		panic(err)
	}
}

func runMigrations(ctx context.Context, connString string) error {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	migrator, err := migrate.NewMigrator(context.Background(), conn, "schema_version")
	if err != nil {
		return err
	}

	err = migrator.LoadMigrations(os.DirFS("../../../migrations/sql"))
	if err != nil {
		return err
	}

	return migrator.Migrate(context.Background())
}
