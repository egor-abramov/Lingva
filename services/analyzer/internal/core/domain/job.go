package domain

import "fmt"

//go:generate stringer -type=Language
type Language int

func (l Language) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *Language) UnmarshalText(text []byte) error {
	const op = "analyzer.core.domain.Language.UnmarshalText"

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

type AnalyzeJob struct {
	Lang      Language `yaml:"lang"`
	Code      string   `yaml:"code"`
	TimeoutMS int      `yaml:"-"`
}

//go:generate stringer -type=Severity
type Severity int

func (s Severity) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *Severity) UnmarshalText(text []byte) error {
	const op = "analyzer.core.domain.Severity.UnmarshalText"

	str := string(text)
	switch str {
	case "Error":
		*s = Error
	case "Warning":
		*s = Warning
	case "Info":
		*s = Info
	default:
		return fmt.Errorf("%s: unknown severity: %s", op, str)
	}
	return nil
}

const (
	Error Severity = iota
	Warning
	Info
)

type Diagnostic struct {
	Line     int32    `yaml:"line"`
	Column   int32    `yaml:"column"`
	Message  string   `yaml:"message"`
	Severity Severity `yaml:"severity"`
}

type AnalyzeResult []Diagnostic
