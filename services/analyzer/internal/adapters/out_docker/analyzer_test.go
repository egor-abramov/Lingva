package out_docker

import (
	"context"
	"lingva/pkg/test"
	"lingva/services/analyzer/internal/core/domain"
	"log/slog"
	"testing"

	"gopkg.in/yaml.v3"
)

type DoubleQuotedString string

func (s DoubleQuotedString) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: yaml.DoubleQuotedStyle,
		Value: string(s),
	}, nil
}

type DiagnosticTest struct {
	Line     int32              `yaml:"line"`
	Column   int32              `yaml:"column"`
	Message  DoubleQuotedString `yaml:"message"`
	Severity domain.Severity    `yaml:"severity"`
}

func TestDockerSandbox_Run(t *testing.T) {
	log := &slog.Logger{}
	analyzer, err := New(log)
	if err != nil {
		t.Fatal(err)
	}

	handler := func(req domain.AnalyzeJob) ([]DiagnosticTest, error) {
		resp, err := analyzer.Run(context.Background(), req)
		if err != nil {
			return nil, err
		}

		testResp := make([]DiagnosticTest, len(resp))
		for i, d := range resp {
			testResp[i] = DiagnosticTest{
				Line:     d.Line,
				Column:   d.Column,
				Message:  DoubleQuotedString(d.Message),
				Severity: d.Severity,
			}
		}
		return testResp, nil
	}
	test.RunTests(t, handler)
}
