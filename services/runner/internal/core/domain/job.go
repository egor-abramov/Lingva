package domain

import (
	l "lingva/pkg/lang"
)

type ExecutionJob struct {
	Lang       l.Language `yaml:"lang"`
	TimeoutMS  int        `yaml:"-"`
	SourceCode string     `yaml:"code"`
}

type OutputChunk struct {
	Data     string `yaml:"data"`
	IsStdErr bool   `yaml:"-"`
}
