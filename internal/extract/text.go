package extract

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExtractText(path, format string) (string, error) {
	switch format {
	case "txt":
		b, err := os.ReadFile(path)
		return string(b), err
	case "csv":
		return extractCSV(path)
	case "json":
		return extractJSON(path)
	case "html":
		return extractHTML(path)
	default:
		return "", nil
	}
}

func extractCSV(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	r := csv.NewReader(f)
	var sb strings.Builder
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return sb.String(), nil
		}
		sb.WriteString(strings.Join(rec, " "))
		sb.WriteString("\n")
	}
	return sb.String(), nil
}

func extractJSON(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return string(b), nil
	}
	var sb strings.Builder
	walkJSON(v, &sb)
	return sb.String(), nil
}

func walkJSON(v any, sb *strings.Builder) {
	switch x := v.(type) {
	case map[string]any:
		for k, v2 := range x {
			sb.WriteString(k + " ")
			walkJSON(v2, sb)
		}
	case []any:
		for _, it := range x {
			walkJSON(it, sb)
		}
	case string:
		sb.WriteString(x + " ")
	case float64, bool, int:
		sb.WriteString(" ")
	}
}

func extractHTML(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		return string(b), nil
	}
	return doc.Text(), nil
}