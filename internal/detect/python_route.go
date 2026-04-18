package detect

func NeedsPython(format string) bool {
	switch format {
	case "pdf", "docx", "xlsx", "xlsm", "tif", "tiff", "jpg", "jpeg", "png", "gif":
		return true
	default:
		return false
	}
}