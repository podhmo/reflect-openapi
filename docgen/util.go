package docgen

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	TRUNCATE_SIZE = 88
)

func toDocumentInfo(title, summary, description string) (di DocumentInfo) {
	// fmt.Fprintf(os.Stderr, "summary: %q\n", summary)
	// fmt.Fprintf(os.Stderr, "description: %q\n", description)
	// fmt.Fprintf(os.Stderr, "--\n")

	defer func() {
		if title != "" && strings.HasPrefix(di.Summary, title) {
			di.Summary = strings.TrimLeft(strings.TrimPrefix(di.Summary, title), " :")
		}
		if len(di.Summary) > TRUNCATE_SIZE {
			di.Summary = di.Summary[:TRUNCATE_SIZE]
		}
	}()

	parts := strings.Split(description, "\n")
	if len(parts) > 2 && strings.TrimSpace(parts[1]) == "" {
		di.Summary = strings.TrimSpace(parts[0])
		di.Description = strings.TrimSpace(strings.Join(parts[2:], "\n"))
		return
	} else if summary != "" {
		di.Summary = strings.TrimSpace(summary)
		di.Description = strings.TrimSpace(description)
		return
	} else {
		di.Summary = strings.TrimSpace(parts[0])
		di.Description = strings.TrimSpace(strings.Join(parts[1:], "\n"))
		return
	}
}

var (
	toDashRegex  = regexp.MustCompile(`[ \t]+`)
	toEmptyRegex = regexp.MustCompile(`[{/\.}]+`)
)

func toHtmlID(s string, xs ...string) string {
	if len(xs) > 0 {
		s = fmt.Sprintf("%s %s", s, strings.Join(xs, " "))
	}
	s = strings.ToLower(s)
	s = toEmptyRegex.ReplaceAllString(s, "")
	s = toDashRegex.ReplaceAllString(s, "-")
	return s
}
