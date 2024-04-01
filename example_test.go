package gopii_test

import (
	"fmt"

	"github.com/blizzy78/gopii"
)

func Example_string() {
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

	// for each output level, print the string with the appropriate obscuring
	fmt.Println(str.String(levelInternal))
	fmt.Println(str.String(levelDevLog))
	fmt.Println(str.String(levelProdLog))

	// Output:
	// JohnDoe:secret
	// Jo***oe:sec***
	// J***:***
}
