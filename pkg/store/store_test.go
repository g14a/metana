package store

import (
	"testing"

	"github.com/g14a/metana/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestProcessLogs(t *testing.T) {
	tests := []struct {
		input      string
		wantNum    int
		wantTitles []string
	}{
		{
			input:   "1621081055-InitSchema.go\n1621084125-AddIndexes.go\n1621084135-AddFKeys.go",
			wantNum: 3,
			wantTitles: []string{
				"1621081055-InitSchema.go",
				"1621084125-AddIndexes.go",
				"1621084135-AddFKeys.go",
			},
		},
		{
			input:      "",
			wantNum:    0,
			wantTitles: nil,
		},
	}

	for _, tt := range tests {
		gotTrack, gotNum := ProcessLogs(tt.input)

		assert.Equal(t, tt.wantNum, gotNum)
		if tt.wantNum == 0 {
			assert.Empty(t, gotTrack.Migrations)
			assert.Empty(t, gotTrack.LastRun)
		} else {
			assert.Equal(t, tt.wantTitles[len(tt.wantTitles)-1], gotTrack.LastRun)
			var gotTitles []string
			for _, m := range gotTrack.Migrations {
				gotTitles = append(gotTitles, m.Title)
				assert.NotEmpty(t, m.ExecutedAt, "ExecutedAt should be set")
			}
			assert.Equal(t, tt.wantTitles, gotTitles)
		}
	}
}

func TestTrackToSetDown(t *testing.T) {
	tests := []struct {
		inputTrack  types.Track
		inputNum    int
		wantTitles  []string
		wantLastRun string
	}{
		{
			inputTrack: types.Track{
				LastRun: "1621095067-Abc.go",
				Migrations: []types.Migration{
					{Title: "1621095067-Abc.go", ExecutedAt: "now"},
				},
			},
			inputNum:    1,
			wantTitles:  nil,
			wantLastRun: "",
		},
		{
			inputTrack: types.Track{
				LastRun: "1621097000-InitSchema.go",
				Migrations: []types.Migration{
					{Title: "1621095067-Abc.go", ExecutedAt: "now"},
					{Title: "1621096992-Random.go", ExecutedAt: "now"},
					{Title: "1621096995-AddIndexes.go", ExecutedAt: "now"},
					{Title: "1621097000-InitSchema.go", ExecutedAt: "now"},
				},
			},
			inputNum: 3,
			wantTitles: []string{
				"1621095067-Abc.go",
			},
			wantLastRun: "1621095067-Abc.go",
		},
	}

	for _, tt := range tests {
		got := TrackToSetDown(tt.inputTrack, tt.inputNum)
		var gotTitles []string
		for _, m := range got.Migrations {
			gotTitles = append(gotTitles, m.Title)
		}
		assert.Equal(t, tt.wantTitles, gotTitles)
		assert.Equal(t, tt.wantLastRun, got.LastRun)
	}
}
