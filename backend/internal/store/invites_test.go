package store_test

import (
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/darko-komericki/pantry/backend/internal/store"
)

func createInvite(t *testing.T, q *store.Queries, h store.Household, by store.User, code string, expiresAt time.Time) {
	t.Helper()
	_, err := q.CreateInvite(t.Context(), store.CreateInviteParams{
		HouseholdID: h.ID,
		Code:        code,
		CreatedBy:   by.ID,
		ExpiresAt:   expiresAt,
	})
	if err != nil {
		t.Fatalf("create invite %s: %v", code, err)
	}
}

func TestConsumeInvite(t *testing.T) {
	q := testQueries(t)
	darko := createUser(t, q, "darko@example.com")
	home := createHousehold(t, q, "Home")
	createInvite(t, q, home, darko, "VALID234", hoursFromNow(24))
	createInvite(t, q, home, darko, "EXPIRED2", hoursFromNow(-1))

	tests := []struct {
		name    string
		code    string
		wantErr error
	}{
		{"valid code", "VALID234", nil},
		{"same code again is gone", "VALID234", pgx.ErrNoRows},
		{"expired code", "EXPIRED2", pgx.ErrNoRows},
		{"unknown code", "NOPE2345", pgx.ErrNoRows},
	}

	// Order matters: the second case relies on the first having used the code.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.ConsumeInvite(t.Context(), tt.code)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got.HouseholdID != home.ID {
				t.Errorf("household = %s, want %s", got.HouseholdID, home.ID)
			}
		})
	}
}

func TestDeleteExpiredInvites(t *testing.T) {
	q := testQueries(t)
	darko := createUser(t, q, "darko@example.com")
	home := createHousehold(t, q, "Home")
	createInvite(t, q, home, darko, "EXPIRED1", hoursFromNow(-1))
	createInvite(t, q, home, darko, "VALID234", hoursFromNow(24))

	deleted, err := q.DeleteExpiredInvites(t.Context())
	if err != nil {
		t.Fatalf("delete expired: %v", err)
	}
	if deleted != 1 {
		t.Errorf("deleted = %d, want 1", deleted)
	}
	if _, err := q.ConsumeInvite(t.Context(), "VALID234"); err != nil {
		t.Errorf("valid invite should survive: %v", err)
	}
}
