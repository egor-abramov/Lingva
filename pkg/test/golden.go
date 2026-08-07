package test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

var Update = flag.Bool("update", false, "update golden files")
var testDir = "testdata"
var testExtension = ".yml"

type Record[Req any, Resp any] struct {
	Input  Req  `yaml:"input"`
	Output Resp `yaml:"output"`
}

type Handler[Req any, Resp any] func(req Req) (Resp, error)

func RunTests[Req any, Resp any](t *testing.T, handler Handler[Req, Resp]) {
	err := filepath.WalkDir(testDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(d.Name()) != testExtension {
			return nil
		}

		relPath, err := filepath.Rel(testDir, path)
		if err != nil {
			return err
		}

		t.Run(relPath, func(t *testing.T) {
			fileBytes, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read golden file: %v", err)
			}

			var record Record[Req, Resp]
			if err := yaml.Unmarshal(fileBytes, &record); err != nil {
				t.Fatalf("Failed to unmarshal golden file: %v", err)
			}

			actualResp, err := handler(record.Input)
			if err != nil {
				t.Fatalf("Handler failed: %v", err)
			}

			if *Update {
				record.Output = actualResp
				updatedBytes, err := yaml.Marshal(record)
				if err != nil {
					t.Fatalf("Failed to marshal updated record: %v", err)
				}

				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatalf("Failed to create directory: %v", err)
				}

				if err := os.WriteFile(path, updatedBytes, 0644); err != nil {
					t.Fatalf("Failed to write golden file: %v", err)
				}
				return
			}

			expectedBytes, _ := yaml.Marshal(record.Output)
			actualBytes, _ := yaml.Marshal(actualResp)

			if string(expectedBytes) != string(actualBytes) {
				t.Errorf("Test mismatch!\nExpected:\n%s\nGot:\n%s", string(expectedBytes), string(actualBytes))
			}
		})
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk test directory: %v", err)
	}
}
