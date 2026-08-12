package parsers

import (
	"bufio"
	"lingva/services/analyzer/internal/core/domain"
	"strconv"
	"strings"
)

func Gcc(stderr string) domain.AnalyzeResult {
	var diagnostics []domain.Diagnostic
	scanner := bufio.NewScanner(strings.NewReader(stderr))

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, ":", 5)
		if len(parts) < 4 || parts[0] != "<stdin>" {
			continue
		}

		lineNum, _ := strconv.Atoi(parts[1])
		var colNum int
		var severityStr, msg string

		if len(parts) >= 5 {
			colNum, _ = strconv.Atoi(parts[2])
			severityStr = strings.TrimSpace(parts[3])
			msg = strings.TrimSpace(parts[4])
		} else {
			severityStr = strings.TrimSpace(parts[2])
			msg = strings.TrimSpace(parts[3])
		}

		var severity domain.Severity
		if strings.Contains(severityStr, "error") {
			severity = domain.Error
		} else if strings.Contains(severityStr, "warning") {
			severity = domain.Warning
		} else {
			severity = domain.Info
		}

		diagnostics = append(diagnostics, domain.Diagnostic{
			Line:     int32(lineNum),
			Column:   int32(colNum),
			Message:  msg,
			Severity: severity,
		})
	}

	if diagnostics == nil {
		diagnostics = make([]domain.Diagnostic, 0)
	}
	return diagnostics
}
