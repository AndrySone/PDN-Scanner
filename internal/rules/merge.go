package rules

import "pii-scanner/internal/model"

func MergeFindings(a, b []model.Finding) []model.Finding {
	type key struct {
		cat string
		tp  string
	}
	m := make(map[key]model.Finding)

	put := func(f model.Finding) {
		k := key{cat: f.Category, tp: f.Type}
		if ex, ok := m[k]; ok {
			ex.Count += f.Count
			if len(ex.MaskedSamples) < 2 && len(f.MaskedSamples) > 0 {
				ex.MaskedSamples = append(ex.MaskedSamples, f.MaskedSamples...)
				if len(ex.MaskedSamples) > 2 {
					ex.MaskedSamples = ex.MaskedSamples[:2]
				}
			}
			if f.Confidence > ex.Confidence {
				ex.Confidence = f.Confidence
			}
			m[k] = ex
			return
		}
		if len(f.MaskedSamples) > 2 {
			f.MaskedSamples = f.MaskedSamples[:2]
		}
		m[k] = f
	}

	for _, f := range a {
		put(f)
	}
	for _, f := range b {
		put(f)
	}

	out := make([]model.Finding, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}