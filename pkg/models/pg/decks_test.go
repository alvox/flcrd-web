package pg

import (
	"alexanderpopov.me/flcrd/pkg/models"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var testDecks = [3]*models.Deck{
	{
		ID:          "9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e",
		Name:        "Test Name 1",
		Description: "Test Description 1",
		Created:     time.Date(2019, 1, 1, 10, 0, 0, 0, time.UTC),
		Public:      false,
		CardsCount:  3,
		CreatedBy: models.DeckCreator{
			ID:   "40afbc9a-27e3-4b38-97f9-2930b8790a9f",
			Name: "Testuser1",
		},
	}, {
		ID:          "2601da50-56a6-41a1-a92e-5624598a7d19",
		Name:        "Test Name 2",
		Description: "Test Description 2",
		Created:     time.Date(2019, 2, 2, 12, 22, 0, 0, time.UTC),
		Public:      true,
		CardsCount:  2,
		CreatedBy: models.DeckCreator{
			ID:   "dd4a5e3a-4d95-44c1-8aa3-e29fa9a29570",
			Name: "Testuser2",
		},
	}, {
		ID:          "4735bdb2-45a5-42d9-a6d4-6db29787b5f1",
		Name:        "Test Name 3",
		Description: "Test Description 3",
		Created:     time.Date(2019, 3, 3, 12, 22, 0, 0, time.UTC),
		Public:      true,
		CardsCount:  0,
		CreatedBy: models.DeckCreator{
			ID:   "40afbc9a-27e3-4b38-97f9-2930b8790a9f",
			Name: "Testuser1",
		},
	},
}

func TestDeckModel_Create_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := DeckModel{db}

	id, err := model.Create("Test Deck", "Deck, created from test", "40AFBC9A-27E3-4B38-97F9-2930B8790A9F", true)
	require.Nil(t, err)
	require.NotNil(t, id)

	deck, err := model.Get(*id)
	require.Nil(t, err)
	require.Equal(t, "Test Deck", deck.Name, "invalid deck name")
	require.Equal(t, "Deck, created from test", deck.Description, "invalid deck description")
}

func TestDeckModel_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	tests := []struct {
		name      string
		deckId    string
		wantDeck  *models.Deck
		wantError error
	}{
		{
			name:      "Deck 1",
			deckId:    "9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e",
			wantDeck:  testDecks[0],
			wantError: nil,
		},
		{
			name:      "Deck 2",
			deckId:    "2601da50-56a6-41a1-a92e-5624598a7d19",
			wantDeck:  testDecks[1],
			wantError: nil,
		},
		{
			name:      "Non-existent ID",
			deckId:    "7af65126-d46c-4797-a329-09d283acc664",
			wantDeck:  nil,
			wantError: models.ErrNoRecord,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			model := DeckModel{db}
			deck, err := model.Get(tt.deckId)
			require.Equal(t, tt.wantError, err)
			require.Equal(t, tt.wantDeck, deck)
		})
	}
}

func TestDeckModel_GetForUser(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	tests := []struct {
		name      string
		userId    string
		wantDecks []*models.Deck
	}{
		{
			name:      "User 1",
			userId:    "40afbc9a-27e3-4b38-97f9-2930b8790a9f",
			wantDecks: []*models.Deck{testDecks[0], testDecks[2]},
		},
		{
			name:      "User 2",
			userId:    "dd4a5e3a-4d95-44c1-8aa3-e29fa9a29570",
			wantDecks: []*models.Deck{testDecks[1]},
		},
		{
			name:      "Non-existent user",
			userId:    "7af65126-d46c-4797-a329-09d283acc664",
			wantDecks: []*models.Deck{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			model := DeckModel{db}
			decks, err := model.GetForUser(tt.userId)
			require.Nil(t, err)
			for i, deck := range tt.wantDecks {
				require.Equal(t, decks[i], deck)
			}
		})
	}
}

func TestDeckModel_GetPublic_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := DeckModel{db}
	decks, total, err := model.GetPublic(0, 5)
	require.Nil(t, err)
	require.Equal(t, 2, len(decks))
	require.Equal(t, decks[0], testDecks[1])
	require.Equal(t, 2, total)
}

func TestDeckModel_Update_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := DeckModel{db}
	deck := testDecks[0]
	deck.Name = "Updated name"
	deck.Description = "Updated description"
	err := model.Update(deck)
	require.Nil(t, err)
	updated, err := model.Get("9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e")
	require.Nil(t, err)
	require.Equal(t, deck, updated)
}

func TestDeckModel_Update_NonExistent(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := DeckModel{db}
	deck := &models.Deck{
		ID:          "7af65126-d46c-4797-a329-09d283acc664",
		Name:        "Updated 1",
		Description: "Updated Description 1",
		Created:     time.Date(2019, 1, 1, 10, 0, 0, 0, time.UTC),
	}
	err := model.Update(deck)
	require.Equal(t, models.ErrNoRecord, err)
}

func TestDeckModel_Delete_Positive(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := DeckModel{db}
	err := model.Delete("9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e")
	require.Nil(t, err)
	_, err = model.Get("9f2556fb-0b84-4b8d-ab0a-b5acb0c89f6e")
	require.Equal(t, models.ErrNoRecord, err)
}

func TestDeckModel_Delete_NonExistent(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	db, teardown := newTestDB(t)
	defer teardown()
	model := DeckModel{db}
	err := model.Delete("7af65126-d46c-4797-a329-09d283acc664")
	require.Equal(t, models.ErrNoRecord, err)
}

func TestDeckModel_Search(t *testing.T) {
	if testing.Short() {
		t.Skip("pg: skipping database test")
	}
	tests := []struct {
		name    string
		terms   []string
		wantLen int
	}{
		{
			name:    "All",
			terms:   []string{"Test"},
			wantLen: 2,
		},
		{
			name:    "Third deck",
			terms:   []string{"Description", "3"},
			wantLen: 1,
		},
		{
			name:    "All, two terms",
			terms:   []string{"Test", "Name"},
			wantLen: 2,
		},
		{
			name:    "None",
			terms:   []string{"Not", "Here"},
			wantLen: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			model := DeckModel{db}
			decks, err := model.Search(tt.terms)
			require.Nil(t, err)
			require.Equal(t, tt.wantLen, len(decks))
		})
	}
}
