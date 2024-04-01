package gopii

import (
	"log/slog"
	"strconv"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestSlogReplaceAttr_LogValue(t *testing.T) {
	obscureFuncs := map[PartType]map[Level]ObscureFunc{
		typePlain: {
			levelMax: Plain(),
		},

		typeLogin: {
			levelDevelopment: Plain(),             // user123456789 -> user123456789
			levelLog:         KeepFirstLast(4, 3), // user123456789 -> user***789
			levelAuditLog:    KeepFirst(3),        // user123456789 -> use***
			levelInternet:    KeepFirst(1),        // user123456789 -> u***
		},

		typePassword: {
			levelDevelopment: Plain(), // pass123456789 -> pass123456789
		},
	}

	str := NewString(obscureFuncs,
		"user123456789", typeLogin,
		":", typePlain,
		"pass123456789", typePassword,
	)

	tests := []struct {
		outLevel Level
		want     string
	}{
		{0, "user123456789:pass123456789"},
		{levelDevelopment, "user123456789:pass123456789"},
		{levelLog - 2, "user***789:***"},
		{levelLog, "user***789:***"},
		{levelAuditLog - 2, "use***:***"},
		{levelAuditLog, "use***:***"},
		{levelInternet - 2, "u***:***"},
		{levelInternet, "u***:***"},
		{100, "***:***"},
	}

	for _, test := range tests {
		t.Run("level "+strconv.Itoa(int(test.outLevel)), func(t *testing.T) {
			is := is.New(t)

			buf := strings.Builder{}

			handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
				Level:       slog.LevelInfo,
				ReplaceAttr: SlogReplaceAttr(test.outLevel),
			})

			logger := slog.New(handler)

			logger.Info("msg", slog.Any("s", str), slog.String("x", "y"))
			is.True(strings.Contains(buf.String(), "s="+test.want))
			is.True(strings.Contains(buf.String(), "x=y"))
		})
	}
}
