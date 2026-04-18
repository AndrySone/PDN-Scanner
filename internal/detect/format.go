package detect

import "strings"

func DetectFormat(path string) string {
	p := strings.ToLower(path)
	switch {
	case strings.HasSuffix(p, ".csv"):
		return "csv"
	case strings.HasSuffix(p, ".json"):
		return "json"
	case strings.HasSuffix(p, ".txt"):
		return "txt"
	case strings.HasSuffix(p, ".html"), strings.HasSuffix(p, ".htm"):
		return "html"
	case strings.HasSuffix(p, ".pdf"):
		return "pdf"
	case strings.HasSuffix(p, ".docx"):
		return "docx"
	case strings.HasSuffix(p, ".xlsx"):
		return "xlsx"
	case strings.HasSuffix(p, ".xlsm"):
		return "xlsm"
	case strings.HasSuffix(p, ".tif"), strings.HasSuffix(p, ".tiff"):
		return "tif"
	case strings.HasSuffix(p, ".jpg"), strings.HasSuffix(p, ".jpeg"):
		return "jpg"
	case strings.HasSuffix(p, ".png"):
		return "png"
	case strings.HasSuffix(p, ".gif"):
		return "gif"
	case strings.HasSuffix(p, ".mp4"):
		return "mp4"
	default:
		return "unsupported"
	}
}