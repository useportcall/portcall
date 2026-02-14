//go:build integration

package testutil

import (
	"context"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/useportcall/portcall/libs/go/dbx"
	gormPG "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupResult holds the DB handle and a cleanup function.
type SetupResult struct {
	DB      dbx.IORM
	Cleanup func()
}

// SetupPostgres spins up a throwaway Postgres container and returns
// an IORM with all tables migrated. Call Cleanup() when done.
func SetupPostgres() *SetupResult {
	ctx := context.Background()
	pg, err := postgres.Run(ctx, "postgres:16",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
		),
	)
	if err != nil {
		log.Fatalf("start postgres container: %v", err)
	}
	dsn, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("get connection string: %v", err)
	}
	gormDB, err := gorm.Open(gormPG.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("open gorm: %v", err)
	}
	db := dbx.NewFromDB(gormDB)
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("auto-migrate: %v", err)
	}
	fmt.Println("integration: postgres ready")
	return &SetupResult{
		DB:      db,
		Cleanup: func() { pg.Terminate(ctx) },
	}
}
