package zaputil

import (
	"fmt"
	"go.uber.org/zap"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
)

func TestGetZapLevel(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		levelString string
		expected    zapcore.Level
		expectError bool
	}{
		{"debug", zapcore.DebugLevel, false},
		{"info", zapcore.InfoLevel, false},
		{"warn", zapcore.WarnLevel, false},
		{"error", zapcore.ErrorLevel, false},
		{"", zapcore.InfoLevel, false},
		{"foobar", zapcore.InfoLevel, true},
	}

	for _, tc := range testCases {
		actual, err := getZapLevel(tc.levelString)
		if tc.expectError && err == nil {
			t.Errorf("Expected error for level %q but got none", tc.levelString)
		} else if !tc.expectError && err != nil {
			t.Errorf("Unexpected error for level %q: %s", tc.levelString, err)
		}
		if actual != tc.expected {
			t.Errorf("For level %q expected %v but got %v", tc.levelString, tc.expected, actual)
		}
	}
}

func TestGetZapEncoder(t *testing.T) {
	t.Parallel()
	// Test valid formats
	testCases := []struct {
		format string
	}{
		{"text"},
		{"color"},
		{"json"},
		{"TEXT"},
		{"COLOR"},
		{"JSON"},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(
			fmt.Sprintf("valid format %s", tc.format), func(t *testing.T) {
				t.Parallel()
				encoder, err := getZapEncoder(Format(tc.format))
				assert.NoError(t, err)
				assert.NotNil(t, encoder)
			},
		)
	}

	// Test unknown format
	unknownFormat := "invalid"
	_, err := getZapEncoder(Format(unknownFormat))
	assert.EqualError(t, err, fmt.Sprintf("unknown log format [text,color,json]: %q", unknownFormat))
}

func TestNewLogger(t *testing.T) {
	t.Parallel()
	logger := NewLogger(WithFormat(FormatJSON))
	assert.NotNil(t, logger)

	logger.With(
		zap.String("url", "url"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	).Info("test")
}
