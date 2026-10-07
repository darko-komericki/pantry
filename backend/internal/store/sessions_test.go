package store_test

import (
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/darko-komericki/pantry/backend/internal/store"
)

func createSession(t *testing.T, q *store.Queries, user store.User, hash string, expiresAt time.Time) {
	t.Helper()
	_, err := q.CreateSession(t.Context(), store.CreateSessionParams{
		UserID:    user.ID,
		TokenHash: []byte(hash),
		ExpiresAt: expiresAt,
	})
	if err != nil {
		t.Fatalf("create session %s: %v", hash, err)
	}
}

func TestGetActiveSessionByTokenHash(t *testing.T) {
	q := testQueries(t)
	user := createUser(t, q, "darko@example.com")
	createSession(t, q, user, "active", hoursFromNow(1))
	createSession(t, q, user, "expired", hoursFromNow(-1))

	tests := []struct {
		name    string
		hash    string
		wantErr error
	}{
		{"active session", "active", nil},
		{"expired session", "expired", pgx.ErrNoRows},
		{"unknown token", "unknown", pgx.ErrNoRows},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.GetActiveSessionByTokenHash(t.Context(), []byte(tt.hash))
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got.UserID != user.ID {
				t.Errorf("session user = %s, want %s", got.UserID, user.ID)
			}
		})
	}
}

func TestDeleteSessionByTokenHash(t *testing.T) {
	q := testQueries(t)
	user := createUser(t, q, "darko@example.com")
	createSession(t, q, user, "this-device", hoursFromNow(1))
	createSession(t, q, user, "other-device", hoursFromNow(1))

	if err := q.DeleteSessionByTokenHash(t.Context(), []byte("this-device")); err != nil {
		t.Fatalf("delete session: %v", err)
	}

	if _, err := q.GetActiveSessionByTokenHash(t.Context(), []byte("this-device")); !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("deleted session: err = %v, want pgx.ErrNoRows", err)
	}
	if _, err := q.GetActiveSessionByTokenHash(t.Context(), []byte("other-device")); err != nil {
		t.Errorf("other session should survive: %v", err)
	}
}

func TestDeleteUserSessions(t *testing.T) {
	q := testQueries(t)
	darko := createUser(t, q, "darko@example.com")
	ana := createUser(t, q, "ana@example.com")
	createSession(t, q, darko, "darko-phone", hoursFromNow(1))
	createSession(t, q, darko, "darko-laptop", hoursFromNow(1))
	createSession(t, q, ana, "ana-phone", hoursFromNow(1))

	if err := q.DeleteUserSessions(t.Context(), darko.ID); err != nil {
		t.Fatalf("delete user sessions: %v", err)
	}

	for _, hash := range []string{"darko-phone", "darko-laptop"} {
		if _, err := q.GetActiveSessionByTokenHash(t.Context(), []byte(hash)); !errors.Is(err, pgx.ErrNoRows) {
			t.Errorf("%s: err = %v, want pgx.ErrNoRows", hash, err)
		}
	}
	if _, err := q.GetActiveSessionByTokenHash(t.Context(), []byte("ana-phone")); err != nil {
		t.Errorf("other user's session should survive: %v", err)
	}
}

func TestDeleteExpiredSessions(t *testing.T) {
	q := testQueries(t)
	user := createUser(t, q, "darko@example.com")
	createSession(t, q, user, "expired-1", hoursFromNow(-1))
	createSession(t, q, user, "expired-2", hoursFromNow(-48))
	createSession(t, q, user, "active", hoursFromNow(1))

	deleted, err := q.DeleteExpiredSessions(t.Context())
	if err != nil {
		t.Fatalf("delete expired: %v", err)
	}
	if deleted != 2 {
		t.Errorf("deleted = %d, want 2", deleted)
	}
	if _, err := q.GetActiveSessionByTokenHash(t.Context(), []byte("active")); err != nil {
		t.Errorf("active session should survive: %v", err)
	}
}
