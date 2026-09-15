package models

import (
	"errors"
	"testing"
	"time"

	"github.com/lunalubra/snippetbox/internal/assert"
)

func TestSnippetModelExtendExpiry(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name      string
		snippetID int
		wantErr   error
	}{
		{
			name:      "Expiring soon",
			snippetID: 1,
			wantErr:   nil,
		},
		{
			name:      "Long lived",
			snippetID: 2,
			wantErr:   ErrNotExtendable,
		},
		{
			name:      "Already expired",
			snippetID: 3,
			wantErr:   ErrNotExtendable,
		},
		{
			name:      "Non-existent ID",
			snippetID: 100,
			wantErr:   ErrNotExtendable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := SnippetModel{db}

			var before time.Time
			err := db.QueryRow("SELECT expires FROM snippets WHERE id = ?", tt.snippetID).Scan(&before)
			if err != nil && tt.snippetID != 100 {
				t.Fatal(err)
			}

			expires, err := m.ExtendExpiry(tt.snippetID)

			if tt.wantErr != nil {
				assert.True(t, errors.Is(err, tt.wantErr))

				var after time.Time
				if tt.snippetID != 100 {
					err = db.QueryRow("SELECT expires FROM snippets WHERE id = ?", tt.snippetID).Scan(&after)
					if err != nil {
						t.Fatal(err)
					}

					assert.Equal(t, after.Equal(before), true)
				}

				return
			}

			assert.Nil(t, err)
			assert.Equal(t, expires.Equal(before.Add(ExtensionPeriod)), true)
		})
	}
}
