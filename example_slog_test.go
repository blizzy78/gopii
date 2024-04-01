package gopii_test

import (
	"log/slog"
	"os"

	"github.com/blizzy78/gopii"
)

func Example_slog() {
	// define string part types
	const (
		typePlain = gopii.PartType(iota + 1)
		typeLogin
		typePassword
	)

	// define output levels
	const (
		levelInternal = gopii.Level(iota + 1)
		levelDevLog
		levelProdLog

		levelWorld
	)

	// define how to obscure each part type at each output level -
	// higher levels are obscured more restrictively
	obscure := map[gopii.PartType]map[gopii.Level]gopii.ObscureFunc{
		// plain text is never obscured
		typePlain: {levelWorld: gopii.Plain()},

		typeLogin: {
			levelInternal: gopii.Plain(),
			levelDevLog:   gopii.KeepFirstLast(2, 2),
			levelProdLog:  gopii.KeepFirst(1),
		},

		typePassword: {
			levelInternal: gopii.Plain(),
			levelDevLog:   gopii.KeepFirst(3),
		},
	}

	// create a new string with parts of different types
	str := gopii.NewString(obscure,
		"JohnDoe", typeLogin,
		":", typePlain,
		"secret", typePassword,
	)

	// create a ReplaceAttr function to obscure strings
	replaceAttr := gopii.SlogReplaceAttr(levelDevLog)

	// create a handler that writes to stdout
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// use replaceAttr to obscure strings -
		// note: please ignore removeExtraAttrs(), this is just for the example
		ReplaceAttr: removeExtraAttrs(replaceAttr),
	})

	// create a new logger
	logger := slog.New(handler)

	// log the string
	logger.Info("", slog.Any("str", str))

	// Output:
	// str=Jo***oe:sec***
}

// removeExtraAttrs removes attributes that are not needed for the example.
func removeExtraAttrs(next func(groups []string, attr slog.Attr) slog.Attr) func(
	groups []string, attr slog.Attr) slog.Attr {
	return func(groups []string, attr slog.Attr) slog.Attr {
		if attr.Key == "time" || attr.Key == "level" || attr.Key == "msg" {
			return slog.Attr{}
		}

		return next(groups, attr)
	}
}
