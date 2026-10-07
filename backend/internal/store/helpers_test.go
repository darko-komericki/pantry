package store_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/darko-komericki/pantry/backend/internal/store"
)

// testQueries returns Queries bound to a transaction on the test database.
// The transaction is rolled back when the test ends, so every test starts
// from empty tables and leaves nothing behind.
func testQueries(t *testing.T) *store.Queries {
	t.Helper()

	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Fatal("TEST_DATABASE_URL is not set; run tests with `make test`")
	}

	conn, err := pgx.Connect(t.Context(), url)
	if err != nil {
		t.Fatalf("connect to test db: %v", err)
	}
	tx, err := conn.Begin(t.Context())
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}

	// t.Context() is already cancelled when cleanup runs, so use a fresh one.
	t.Cleanup(func() {
		ctx := context.Background()
		_ = tx.Rollback(ctx) // nothing to do if it fails: closing the conn discards the tx anyway
		_ = conn.Close(ctx)
	})

	return store.New(tx)
}

func createUser(t *testing.T, q *store.Queries, email string) store.User {
	t.Helper()
	u, err := q.CreateUser(t.Context(), store.CreateUserParams{
		Email:        email,
		PasswordHash: "not-a-real-hash",
		DisplayName:  "Test User",
	})
	if err != nil {
		t.Fatalf("create user %s: %v", email, err)
	}
	return u
}

func createHousehold(t *testing.T, q *store.Queries, name string) store.Household {
	t.Helper()
	h, err := q.CreateHousehold(t.Context(), name)
	if err != nil {
		t.Fatalf("create household %s: %v", name, err)
	}
	return h
}

// isUniqueViolation reports whether err is Postgres error 23505 (unique_violation).
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func hoursFromNow(h int) time.Time {
	return time.Now().Add(time.Duration(h) * time.Hour)
}
