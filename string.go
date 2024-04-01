package gopii

import "strings"

// String is a representation of a string whose parts can be obscured.
type String struct {
	obscureFuncs map[PartType]map[Level]ObscureFunc
	parts        []part
}

// PartType is the type of a part of a string. Different part types can be obscured differently.
type PartType int

type part struct {
	s   string
	typ PartType
}

// NewString returns a new String containing the given parts. More parts can be added using Append or AppendParts.
// obscureFuncs determines how each part type is obscured at each output level.
func NewString(obscureFuncs map[PartType]map[Level]ObscureFunc, parts ...any) *String {
	s := &String{obscureFuncs: obscureFuncs}
	s.AppendParts(parts...)

	return s
}

// Append appends a new part to the String.
func (s *String) Append(str string, typ PartType) {
	s.parts = append(s.parts, part{s: str, typ: typ})
}

// AppendParts appends new parts to the String. parts must be in pairs of string and PartType.
func (s *String) AppendParts(parts ...any) {
	for i := 0; i < len(parts); i += 2 {
		s.Append(parts[i].(string), parts[i+1].(PartType)) //nolint:forcetypeassert // must conform to these types
	}
}

// String returns the string with its parts obscured according to outLevel and the String's obscureFuncs.
// For each part, if no ObscureFunc is found for outLevel, the ObscureFunc for the next higher level is used.
// If there is no ObscureFunc for a higher level, Obscured() is used.
func (s *String) String(outLevel Level) string {
	buf := strings.Builder{}

	for _, part := range s.parts {
		obscureFuncs, ok := s.obscureFuncs[part.typ]

		if !ok {
			buf.WriteString(Obscured()(part.s))
			continue
		}

		buf.WriteString(obscure(part.s, outLevel, obscureFuncs))
	}

	return buf.String()
}
