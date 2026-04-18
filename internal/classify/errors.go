package classify

import "strings"

func CategorizeErrors(errors []string) []string {
	set := map[string]struct{}{}

	for _, e := range errors {
		le := strings.ToLower(e)

		switch {
		case strings.Contains(le, "no such file or directory"),
			strings.Contains(le, "failed to open file"),
			strings.Contains(le, "package not found"):
			set["file_open_error"] = struct{}{}

		case strings.Contains(le, "unsupported_media"),
			strings.Contains(le, "unsupported for python worker"),
			strings.Contains(le, "unsupported image format/type"):
			set["unsupported_format"] = struct{}{}

		case strings.Contains(le, "tesseract"),
			strings.Contains(le, "ocr"),
			strings.Contains(le, "cannot identify image file"):
			set["ocr_error"] = struct{}{}

		case strings.Contains(le, "timeout"),
			strings.Contains(le, "context deadline exceeded"):
			set["timeout_error"] = struct{}{}

		case strings.Contains(le, "panic recovered"):
			set["internal_panic"] = struct{}{}

		case strings.Contains(le, "pdf"),
			strings.Contains(le, "docx"),
			strings.Contains(le, "xlsx"),
			strings.Contains(le, "zip"):
			set["corrupted_document"] = struct{}{}

		default:
			set["unknown_error"] = struct{}{}
		}
	}

	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	return out
}

func DetermineStatus(findingsCount int, errors []string) string {
	if len(errors) == 0 {
		return "ok"
	}
	if findingsCount > 0 {
		return "partial"
	}
	return "failed"
}