package lang

import "fmt"

//go:generate stringer -type=Language
type Language int

func (l Language) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *Language) UnmarshalText(text []byte) error {
	const op = "runner.core.domain.Language.UnmarshalText"

	str := string(text)
	switch str {
	case "Python":
		*l = Python
	case "Go":
		*l = Go
	case "C":
		*l = C
	case "Cpp":
		*l = Cpp
	default:
		return fmt.Errorf("%s: unknown language: %s", op, str)
	}
	return nil
}

const (
	Python Language = iota
	Go
	C
	Cpp
)
