package gopii

import (
	"strconv"
	"testing"

	"github.com/matryer/is"
)

type obscureFuncTestCase struct {
	s    string
	want string
}

func TestPlain(t *testing.T) {
	testObscureFunc(t, Plain(), []obscureFuncTestCase{
		{"", ""},
		{"a", "a"},
		{"ab", "ab"},
		{"abc", "abc"},
		{"abcde", "abcde"},
		{"abcdefghijklm", "abcdefghijklm"},
	})
}

func TestObscured(t *testing.T) {
	testObscureFunc(t, Obscured(), []obscureFuncTestCase{
		{"", "***"},
		{"a", "***"},
		{"ab", "***"},
		{"abc", "***"},
		{"abcde", "***"},
		{"abcdefghijklm", "***"},
	})
}

func TestKeepFirst(t *testing.T) {
	testObscureFunc(t, KeepFirst(3), []obscureFuncTestCase{
		{"", ""},
		{"a", "a"},
		{"ab", "ab"},
		{"abc", "abc"},
		{"abcde", "abc***"},
		{"abcdefghijklm", "abc***"},
	})
}

func TestKeepLast(t *testing.T) {
	testObscureFunc(t, KeepLast(3), []obscureFuncTestCase{
		{"", ""},
		{"a", "a"},
		{"ab", "ab"},
		{"abc", "abc"},
		{"abcde", "***cde"},
		{"abcdefghijklm", "***klm"},
	})
}

func TestKeepFirstLast(t *testing.T) {
	testObscureFunc(t, KeepFirstLast(2, 3), []obscureFuncTestCase{
		{"", ""},
		{"a", "a"},
		{"ab", "ab"},
		{"abc", "abc"},
		{"abcde", "abcde"},
		{"abcdefg", "ab***efg"},
		{"abcdefghijklm", "ab***klm"},
	})
}

func testObscureFunc(t *testing.T, obscureFunc ObscureFunc, tests []obscureFuncTestCase) {
	t.Helper()

	for _, test := range tests {
		t.Run(test.s, func(t *testing.T) {
			is := is.New(t)

			is.Equal(obscureFunc(test.s), test.want)
		})
	}
}

func TestObscure(t *testing.T) {
	obscureFuncs := map[Level]ObscureFunc{
		levelDevelopment: Plain(),             // user123456789 -> user123456789
		levelLog:         KeepFirstLast(4, 3), // user123456789 -> user***789
		levelAuditLog:    KeepFirst(3),        // user123456789 -> use***
		levelInternet:    KeepFirst(1),        // user123456789 -> u***
	}

	tests := []struct {
		outLevel Level
		want     string
	}{
		{0, "user123456789"},
		{levelDevelopment, "user123456789"},
		{levelLog - 2, "user***789"},
		{levelLog, "user***789"},
		{levelAuditLog - 2, "use***"},
		{levelAuditLog, "use***"},
		{levelInternet - 2, "u***"},
		{levelInternet, "u***"},
		{100, "***"},
	}

	for _, test := range tests {
		t.Run("level "+strconv.Itoa(int(test.outLevel)), func(t *testing.T) {
			is := is.New(t)

			is.Equal(obscure("user123456789", test.outLevel, obscureFuncs), test.want)
		})
	}
}

func TestObscure_NoFunc(t *testing.T) {
	is := is.New(t)

	is.Equal(obscure("user123456789", levelLog, map[Level]ObscureFunc{}), "***")
}
