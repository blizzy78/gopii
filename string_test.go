package gopii

import (
	"strconv"
	"testing"

	"github.com/matryer/is"
)

const (
	typePlain = PartType(iota)
	typeLogin
	typePassword
)

const (
	levelDevelopment = Level(4)
	levelLog         = Level(8)
	levelAuditLog    = Level(12)
	levelInternet    = Level(16)

	levelMax = Level(100)
)

func TestString_String(t *testing.T) {
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

			is.Equal(str.String(test.outLevel), test.want)
		})
	}
}

func TestString_String_NoFunc(t *testing.T) {
	is := is.New(t)

	obscureFuncs := map[PartType]map[Level]ObscureFunc{}

	str := NewString(obscureFuncs,
		"user123456789", typeLogin,
		":", typePlain,
		"pass123456789", typePassword,
	)

	is.Equal(str.String(levelLog), "*********")
}
