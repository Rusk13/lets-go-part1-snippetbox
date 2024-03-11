package main

import (
	"snippetbox.olegmonabaka.net/internal/assert"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 3, 6, 16, 45, 0, 0, time.UTC),
			want: "06 Mar 2024 at 16:45",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "JST",
			tm:   time.Date(2024, 3, 6, 16, 45, 0, 0, time.FixedZone("JST", -6*60*60)),
			want: "06 Mar 2024 at 22:45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}

}
