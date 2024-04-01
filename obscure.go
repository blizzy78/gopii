package gopii

// Level is a level of obscuring.
type Level int

// ObscureFunc is a function that partially or completely obscures s.
type ObscureFunc func(s string) string

// Plain returns an ObscureFunc that does not obscure the string.
func Plain() ObscureFunc {
	return func(s string) string {
		return s
	}
}

// Obscured returns an ObscureFunc that completely obscures the string.
func Obscured() ObscureFunc {
	return func(_ string) string {
		return "***"
	}
}

// KeepFirst returns an ObscureFunc that obscures all but the first num runes of the string.
func KeepFirst(num int) ObscureFunc {
	return func(s string) string {
		runes := []rune(s)

		if len(runes) <= num {
			return s
		}

		return string(runes[:num]) + "***"
	}
}

// KeepLast returns an ObscureFunc that obscures all but the last num runes of the string.
func KeepLast(num int) ObscureFunc {
	return func(s string) string {
		runes := []rune(s)

		if len(runes) <= num {
			return s
		}

		return "***" + string(runes[len(runes)-num:])
	}
}

// KeepFirstLast returns an ObscureFunc that obscures all but the first numFirst and the last numLast runes of the string.
func KeepFirstLast(numFirst int, numLast int) ObscureFunc {
	return func(s string) string {
		runes := []rune(s)

		if len(runes) <= numFirst+numLast {
			return s
		}

		return string(runes[:numFirst]) + "***" + string(runes[len(runes)-numLast:])
	}
}

// obscure returns a partially or completely obscured version of str, according to outLevel and obscureFuncs.
// If no ObscureFunc is found for outLevel, the ObscureFunc for the next higher level is used.
// If there is no ObscureFunc for a higher level, Obscured() is used.
func obscure(str string, outLevel Level, obscureFuncs map[Level]ObscureFunc) string {
	highestLevel := Level(0)
	for l := range obscureFuncs {
		if l <= highestLevel {
			continue
		}

		highestLevel = l
	}

	obscureFunc := Obscured()

	for lvl := highestLevel; lvl >= outLevel; lvl-- {
		f, ok := obscureFuncs[lvl]
		if !ok {
			continue
		}

		obscureFunc = f
	}

	return obscureFunc(str)
}
