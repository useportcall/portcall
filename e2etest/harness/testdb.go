package harness

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/useportcall/portcall/libs/go/dbx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func swapDBName(dsn, name string) string {
	idx := strings.LastIndex(dsn, "/")
	qIdx := strings.Index(dsn[idx:], "?")
	if qIdx > 0 {
		return dsn[:idx] + "/" + name + dsn[idx+qIdx:]
	}
	return dsn[:idx] + "/" + name
}

const defaultDSN = "postgresql://admin:adminpassword@localhost:5432/postgres?sslmode=disable"

var fallbackDSNs = []string{
	defaultDSN,
	"postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable",
}

var resolvedDSN string

func baseDSN() string {
	if resolvedDSN != "" {
		return resolvedDSN
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		resolvedDSN = v
		return v
	}
	for _, dsn := range fallbackDSNs {
		if tryConnect(dsn) {
			resolvedDSN = dsn
			return dsn
		}
	}
	return defaultDSN
}

func ensurePostgresRunning(t *testing.T) {
	t.Helper()
	dsn := swapDBName(baseDSN(), "postgres")
	if tryConnect(dsn) {
		return
	}
	composeDir := findComposeDir()
	if composeDir == "" {
		t.Fatal("Cannot connect to Postgres and docker-compose/ directory not found. " +
			"Start Postgres with: docker compose -f docker-compose/docker-compose.db.yml up -d")
	}
	exec.Command("docker", "rm", "-f", "postgres_instance").Run()
	t.Logf("Starting postgres_instance via docker compose...")
	cmd := exec.Command("docker", "compose", "-f",
		composeDir+"/docker-compose.db.yml", "up", "-d", "postgres")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to start Postgres container: %v", err)
	}
	for i := 0; i < 30; i++ {
		if tryConnect(dsn) {
			return
		}
		exec.Command("sleep", "1").Run()
	}
	t.Fatal("Postgres container started but cannot connect after 30 retries")
}

func tryConnect(dsn string) bool {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return false
	}
	defer db.Close()
	return db.Ping() == nil
}

func findComposeDir() string {
	dir, _ := os.Getwd()
	for {
		candidate := dir + "/docker-compose"
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
		parent := dir[:strings.LastIndex(dir, "/")]
		if parent == dir || parent == "" {
			return ""
		}
		dir = parent
	}
}

// NewTestDB creates a temporary Postgres database, runs AutoMigrate,
// and registers a cleanup that closes connections and drops the database.
func NewTestDB(t *testing.T) dbx.IORM {
	t.Helper()
	ensurePostgresRunning(t)
	dbName := fmt.Sprintf("e2e_%s_%d", sanitize(t.Name()), os.Getpid())
	adminDSN := swapDBName(baseDSN(), "postgres")

	admin, err := sql.Open("pgx", adminDSN)
	if err != nil {
		t.Fatalf("connect admin db: %v", err)
	}
	defer admin.Close()

	admin.Exec("DROP DATABASE IF EXISTS " + dbName)
	if _, err := admin.Exec("CREATE DATABASE " + dbName); err != nil {
		t.Fatalf("create temp db: %v", err)
	}

	gormDB, err := gorm.Open(postgres.Open(swapDBName(baseDSN(), dbName)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open temp db: %v", err)
	}

	t.Cleanup(func() {
		if sqlDB, err := gormDB.DB(); err == nil {
			sqlDB.Close()
		}
		a, _ := sql.Open("pgx", adminDSN)
		defer a.Close()
		a.Exec(fmt.Sprintf(
			"SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname='%s'", dbName))
		a.Exec("DROP DATABASE IF EXISTS " + dbName)
	})

	orm := dbx.NewFromDB(gormDB)
	if err := orm.AutoMigrate(); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return orm
}

func sanitize(s string) string {
	return strings.ToLower(strings.NewReplacer("/", "_", " ", "_", "-", "_").Replace(s))
}
