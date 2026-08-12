package domain

import (
	"fmt"
	l "lingva/pkg/lang"
)

type AnalyzeJob struct {
	Lang      l.Language `yaml:"lang"`
	Code      string     `yaml:"code"`
	TimeoutMS int        `yaml:"-"`
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
