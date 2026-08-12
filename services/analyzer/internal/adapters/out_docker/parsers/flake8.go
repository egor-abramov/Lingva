package parsers

import (
	"bufio"
	"fmt"
	"lingva/services/analyzer/internal/core/domain"
	"strconv"
	"strings"
)

func Flake8(stdout string, stderr string) domain.AnalyzeResult {
	var diagnostics []domain.Diagnostic
	scanner := bufio.NewScanner(strings.NewReader(stdout))

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, ":", 4)

		if len(parts) < 4 || parts[0] != "stdin" {
			continue
		}

		lineNum, _ := strconv.Atoi(parts[1])
		colNum, _ := strconv.Atoi(parts[2])
		msg := strings.TrimSpace(parts[3])

		var severity domain.Severity
		if strings.HasPrefix(msg, "E9") {
			severity = domain.Error
		} else if strings.HasPrefix(msg, "F") {
			severity = domain.Warning
		} else {
			severity = domain.Info
		}

		msg = strings.Join(strings.Split(msg, " ")[1:], " ")
		diagnostics = append(diagnostics, domain.Diagnostic{
			Line:     int32(lineNum),
			Column:   int32(colNum),
			Message:  msg,
			Severity: severity,
		})
	}

	if stderr != "" && len(diagnostics) == 0 {
		diagnostics = append(diagnostics, domain.Diagnostic{
			Line:     0,
			Column:   0,
			Message:  fmt.Sprintf("Internal Linter Error: %s", strings.TrimSpace(stderr)),
			Severity: domain.Error,
		})
	}

	if diagnostics == nil {
		diagnostics = make([]domain.Diagnostic, 0)
	}

	return diagnostics
}
