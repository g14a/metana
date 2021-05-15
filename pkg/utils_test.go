package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetComponents(t *testing.T) {
	tests := []struct {
		input       string
		resultTS    int
		resultMName string
		wantErr     bool
	}{
		{
			input:       "1621081055-Random.go",
			resultTS:    1621081055,
			resultMName: "Random",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		ts, mn, err := GetComponents(tt.input)
		assert.Equal(t, tt.resultTS, ts)
		assert.Equal(t, tt.resultMName, mn)
		assert.Equal(t, tt.wantErr, err != nil)
	}
}
