package store

import (
	"github.com/g14a/metana/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProcessLogs(t *testing.T) {
	tests := []struct {
		input       string
		wantedTrack types.Track
		wantedNum   int
	}{
		{
			input: "1621081055-InitSchema.go\n1621084125-AddIndexes.go\n1621084135-AddFKeys.go",
			wantedTrack: types.Track{
				LastRun:   "1621084135-AddFKeys.go",
				LastRunTS: 1621084135,
				Migrations: []types.Migration{
					{
						Title:     "1621081055-InitSchema.go",
						Timestamp: 1621081055,
					}, {
						Title:     "1621084125-AddIndexes.go",
						Timestamp: 1621084125,
					}, {
						Title:     "1621084135-AddFKeys.go",
						Timestamp: 1621084135,
					},
				},
			},
			wantedNum: 3,
		}, {}, {
			input:       "",
			wantedTrack: types.Track{},
			wantedNum:   0,
		},
	}

	for _, tt := range tests {
		wantedTrack, num := ProcessLogs(tt.input)
		assert.Equal(t, tt.wantedNum, num)
		assert.Equal(t, tt.wantedTrack, wantedTrack)
	}
}

func TestTrackToSetDown(t *testing.T) {
	tests := []struct {
		inputTrack  types.Track
		inputNum    int
		wantedTrack types.Track
	}{
		{
			inputTrack: types.Track{
				LastRun:   "1621095067-Abc.go",
				LastRunTS: int(time.Now().Unix()),
				Migrations: []types.Migration{
					{
						Title:     "1621095067-Abc.go",
						Timestamp: 1621095067,
					},
				},
			},
			inputNum:    1,
			wantedTrack: types.Track{},
		},
		{
			inputTrack: types.Track{
				LastRun:   "1621097000-InitSchema.go",
				LastRunTS: int(time.Now().Unix()),
				Migrations: []types.Migration{
					{
						Title:     "1621095067-Abc.go",
						Timestamp: 1621095067,
					}, {
						Title:     "1621096992-Random.go",
						Timestamp: 1621096992,
					}, {
						Title:     "1621096995-AddIndexes.go",
						Timestamp: 1621096995,
					}, {
						Title:     "1621097000-InitSchema.go",
						Timestamp: 1621097000,
					},
				},
			},
			inputNum: 3,
			wantedTrack: types.Track{
				LastRun: "1621095067-Abc.go",
				LastRunTS: 1621095067,
				Migrations: []types.Migration{
					{
						Title:     "1621095067-Abc.go",
						Timestamp: 1621095067,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		resultTrack := TrackToSetDown(tt.inputTrack, tt.inputNum)
		assert.Equal(t, tt.wantedTrack, resultTrack)
	}
}
