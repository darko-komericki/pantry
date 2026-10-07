package store_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/darko-komericki/pantry/backend/internal/store"
)

func TestGetUserByEmail(t *testing.T) {
	q := testQueries(t)
	created := createUser(t, q, "Darko@Example.com")

	tests := []struct {
		name    string
		email   string
		wantErr error
	}{
		{"exact case", "Darko@Example.com", nil},
		{"different case", "darko@example.COM", nil},
		{"unknown email", "nobody@example.com", pgx.ErrNoRows},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.GetUserByEmail(t.Context(), tt.email)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got.ID != created.ID {
				t.Errorf("got user %s, want %s", got.ID, created.ID)
			}
		})
	}
}

func TestCreateUserRejectsDuplicateEmailIgnoringCase(t *testing.T) {
	q := testQueries(t)
	createUser(t, q, "darko@example.com")

	// The failed insert aborts the transaction, so this must be the last query.
	_, err := q.CreateUser(t.Context(), store.CreateUserParams{
		Email:        "DARKO@example.com",
		PasswordHash: "not-a-real-hash",
		DisplayName:  "Other",
	})
	if !isUniqueViolation(err) {
		t.Fatalf("err = %v, want unique violation", err)
	}
}

func TestGetUserByID(t *testing.T) {
	q := testQueries(t)
	created := createUser(t, q, "darko@example.com")

	got, err := q.GetUserByID(t.Context(), created.ID)
	if err != nil {
		t.Fatalf("get existing user: %v", err)
	}
	if got.Email != created.Email {
		t.Errorf("email = %q, want %q", got.Email, created.Email)
	}

	_, err = q.GetUserByID(t.Context(), uuid.New())
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("unknown id: err = %v, want pgx.ErrNoRows", err)
	}
}
