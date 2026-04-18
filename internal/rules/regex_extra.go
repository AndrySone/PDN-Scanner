package rules

import (
	"regexp"

	"pii-scanner/internal/model"
)

var (
	rePassportRF = regexp.MustCompile(`(?i)(?:паспорт[^0-9]{0,20})?(\d{2}\s?\d{2}\s?\d{6})`)
	reBik        = regexp.MustCompile(`\b04\d{7}\b`)
	reAccount    = regexp.MustCompile(`\b\d{20}\b`)
	reBirthDate  = regexp.MustCompile(`(?i)(?:дата\s*рождени[яе]|д\.р\.)[^0-9]{0,10}(\d{2}[.\-/]\d{2}[.\-/]\d{4})`)
	reMRZLike    = regexp.MustCompile(`\b[P|I|V|A|C][A-Z0-9<]{20,}\b`)
)

func DetectPIIExtra(text string) []model.Finding {
	type agg struct {
		count   int
		samples []string
	}
	m := map[string]*agg{}

	add := func(tp string, sample string) {
		if _, ok := m[tp]; !ok {
			m[tp] = &agg{}
		}
		m[tp].count++
		if len(m[tp].samples) < 2 {
			m[tp].samples = append(m[tp].samples, sample)
		}
	}

	for range rePassportRF.FindAllStringSubmatch(text, -1) {
		add("passport_rf", "**** ******")
	}
	for _, x := range reBik.FindAllString(text, -1) {
		_ = x
		add("bik", "04*******")
	}
	for _, x := range reAccount.FindAllString(text, -1) {
		_ = x
		add("bank_account", "********************")
	}
	for _, x := range reBirthDate.FindAllString(text, -1) {
		_ = x
		add("birth_date", "**.**.****")
	}
	for _, x := range reMRZLike.FindAllString(text, -1) {
		_ = x
		add("mrz", "P<********************")
	}

	out := make([]model.Finding, 0, len(m))
	for tp, v := range m {
		f := model.Finding{
			Type:          tp,
			Count:         v.count,
			Confidence:    0.92,
			MaskedSamples: v.samples,
		}
		switch tp {
		case "passport_rf", "mrz":
			f.Category = "gov_id"
		case "bik", "bank_account":
			f.Category = "payment"
		case "birth_date":
			f.Category = "ordinary"
		default:
			f.Category = "ordinary"
		}
		out = append(out, f)
	}
	return out
}