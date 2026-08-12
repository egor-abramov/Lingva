package parsers

import (
	"bufio"
	"lingva/services/analyzer/internal/core/domain"
	"strconv"
	"strings"
)

func Go(stderr string) domain.AnalyzeResult {
	var diagnostics []domain.Diagnostic
	scanner := bufio.NewScanner(strings.NewReader(stderr))

	for scanner.Scan() {
		line := scanner.Text()

		idx := strings.Index(line, "main.go:")
		if idx == -1 {
			continue
		}

		cleanLine := line[idx+len("main.go:"):]
		parts := strings.SplitN(cleanLine, ":", 3)

		if len(parts) < 2 {
			continue
		}

		lineNum, _ := strconv.Atoi(parts[0])
		var colNum int
		var msg string

		col, err := strconv.Atoi(parts[1])
		if err == nil {
			colNum = col
			if len(parts) > 2 {
				msg = strings.TrimSpace(parts[2])
			}
		} else {
			reParts := strings.SplitN(cleanLine, ":", 2)
			if len(reParts) > 1 {
				msg = strings.TrimSpace(reParts[1])
			}
		}

		diagnostics = append(diagnostics, domain.Diagnostic{
			Line:     int32(lineNum),
			Column:   int32(colNum),
			Message:  msg,
			Severity: domain.Error,
		})
	}

	if stderr != "" && len(diagnostics) == 0 {
		diagnostics = append(diagnostics, domain.Diagnostic{
			Line:     0,
			Column:   0,
			Message:  strings.TrimSpace(stderr),
			Severity: domain.Error,
		})
	}

	if diagnostics == nil {
		diagnostics = make([]domain.Diagnostic, 0)
	}
	return diagnostics
}
