package scan

import (
	"os"
	"path/filepath"

	"pii-scanner/internal/detect"
	"pii-scanner/internal/model"
)

func Walk(root, scanID string, out chan<- model.FileTask) error {
	defer close(out)
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, e := d.Info()
		if e != nil {
			return nil
		}
		format := detect.DetectFormat(path)
		if format == "unsupported" {
			return nil
		}
		out <- model.FileTask{
			ScanID: scanID,
			Path:   path,
			Format: format,
			Size:   info.Size(),
		}
		return nil
	})
}