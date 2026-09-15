package main

import (
	"testing"
	"time"

	"github.com/lunalubra/snippetbox/internal/assert"
)

func TestCanExtend(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want bool
	}{
		{
			name: "Expired",
			tm:   time.Now().Add(-1 * time.Hour),
			want: false,
		},
		{
			name: "One hour away",
			tm:   time.Now().Add(1 * time.Hour),
			want: true,
		},
		{
			name: "Just inside the window",
			tm:   time.Now().Add(72*time.Hour - time.Minute),
			want: true,
		},
		{
			name: "Just outside the window",
			tm:   time.Now().Add(73 * time.Hour),
			want: false,
		},
		{
			name: "Far future",
			tm:   time.Now().Add(365 * 24 * time.Hour),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, canExtend(tt.tm), tt.want)
		})
	}
}

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2024 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2024 at 09:15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
		})
	}
}
