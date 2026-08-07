package domain

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
	default:
		return fmt.Errorf("%s: unknown language: %s", op, str)
	}
	return nil
}

const (
	Python Language = iota
)

type ExecutionJob struct {
	Lang       Language `yaml:"lang"`
	TimeoutMS  int      `yaml:"-"`
	SourceCode string   `yaml:"code"`
}

type OutputChunk struct {
	Data     string `yaml:"data"`
	IsStdErr bool   `yaml:"-"`
}
