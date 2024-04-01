package gopii

import "log/slog"

// SlogReplaceAttr returns a function that replaces *String values with obscured strings,
// according to outLevel.
// The returned function is intended to be used as slog.HandlerOptions.ReplaceAttr.
func SlogReplaceAttr(outLevel Level) func(groups []string, attr slog.Attr) slog.Attr {
	return func(_ []string, attr slog.Attr) slog.Attr {
		str, ok := attr.Value.Any().(*String)
		if !ok {
			return attr
		}

		return slog.String(attr.Key, str.String(outLevel))
	}
}
