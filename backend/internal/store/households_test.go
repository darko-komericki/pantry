package store_test

import (
	"testing"

	"github.com/darko-komericki/pantry/backend/internal/store"
)

func addMember(t *testing.T, q *store.Queries, h store.Household, u store.User) {
	t.Helper()
	_, err := q.AddHouseholdMember(t.Context(), store.AddHouseholdMemberParams{
		HouseholdID: h.ID,
		UserID:      u.ID,
	})
	if err != nil {
		t.Fatalf("add %s to %s: %v", u.Email, h.Name, err)
	}
}

func TestListUserHouseholds(t *testing.T) {
	q := testQueries(t)
	darko := createUser(t, q, "darko@example.com")
	ana := createUser(t, q, "ana@example.com")
	home := createHousehold(t, q, "Home")
	parents := createHousehold(t, q, "Parents")
	other := createHousehold(t, q, "Ana's flat")
	addMember(t, q, home, darko)
	addMember(t, q, parents, darko)
	addMember(t, q, other, ana)

	got, err := q.ListUserHouseholds(t.Context(), darko.ID)
	if err != nil {
		t.Fatalf("list households: %v", err)
	}

	// Ordered by join time, and only Darko's households.
	want := []string{"Home", "Parents"}
	if len(got) != len(want) {
		t.Fatalf("got %d households, want %d", len(got), len(want))
	}
	for i, h := range got {
		if h.Name != want[i] {
			t.Errorf("household[%d] = %q, want %q", i, h.Name, want[i])
		}
	}
}

func TestIsHouseholdMember(t *testing.T) {
	q := testQueries(t)
	darko := createUser(t, q, "darko@example.com")
	ana := createUser(t, q, "ana@example.com")
	home := createHousehold(t, q, "Home")
	addMember(t, q, home, darko)

	tests := []struct {
		name string
		user store.User
		want bool
	}{
		{"member", darko, true},
		{"not a member", ana, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.IsHouseholdMember(t.Context(), store.IsHouseholdMemberParams{
				HouseholdID: home.ID,
				UserID:      tt.user.ID,
			})
			if err != nil {
				t.Fatalf("is member: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddHouseholdMemberRejectsDuplicate(t *testing.T) {
	q := testQueries(t)
	darko := createUser(t, q, "darko@example.com")
	home := createHousehold(t, q, "Home")
	addMember(t, q, home, darko)

	_, err := q.AddHouseholdMember(t.Context(), store.AddHouseholdMemberParams{
		HouseholdID: home.ID,
		UserID:      darko.ID,
	})
	if !isUniqueViolation(err) {
		t.Fatalf("err = %v, want unique violation", err)
	}
}
