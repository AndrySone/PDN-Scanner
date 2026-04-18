package rules

import (
	"regexp"
	"strings"

	"pii-scanner/internal/model"
)

var (
	reEmail = regexp.MustCompile(`(?i)\b[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}\b`)
	rePhone = regexp.MustCompile(`(?:\+7|8)[\s\-\(]?\d{3}[\s\-\)]?\d{3}[\s\-]?\d{2}[\s\-]?\d{2}`)
	reCard  = regexp.MustCompile(`\b(?:\d[ -]*?){13,19}\b`)
	reSnils = regexp.MustCompile(`\b\d{3}-?\d{3}-?\d{3}\s?\d{2}\b`)
	reInn   = regexp.MustCompile(`\b\d{10}\b|\b\d{12}\b`)
)

func DetectPII(text string) []model.Finding {
	type agg struct {
		count   int
		samples []string
	}
	m := map[string]*agg{}

	add := func(key, sample string) {
		if _, ok := m[key]; !ok {
			m[key] = &agg{}
		}
		m[key].count++
		if len(m[key].samples) < 3 {
			m[key].samples = append(m[key].samples, sample)
		}
	}

	for _, x := range reEmail.FindAllString(text, -1) {
		add("email", maskEmail(x))
	}
	for _, x := range rePhone.FindAllString(text, -1) {
		add("phone", maskPhone(x))
	}
	for _, x := range reCard.FindAllString(text, -1) {
		if LuhnValid(x) {
			add("card_pan", maskCard(x))
		}
	}
	for _, x := range reSnils.FindAllString(text, -1) {
		if SnilsValid(x) {
			add("snils", "***-***-*** **")
		}
	}
	for _, x := range reInn.FindAllString(text, -1) {
		if InnValid(x) {
			add("inn", maskINN(x))
		}
	}

	out := make([]model.Finding, 0, len(m))
	for k, v := range m {
		f := model.Finding{
			Type:          k,
			Count:         v.count,
			Confidence:    0.98,
			MaskedSamples: v.samples,
		}
		switch k {
		case "email", "phone":
			f.Category = "ordinary"
		case "snils", "inn":
			f.Category = "gov_id"
		case "card_pan":
			f.Category = "payment"
		default:
			f.Category = "ordinary"
		}
		out = append(out, f)
	}
	return out
}

func maskEmail(s string) string {
	parts := strings.Split(s, "@")
	if len(parts) != 2 || len(parts[0]) < 2 {
		return "***@***"
	}
	return parts[0][:2] + "***@" + parts[1]
}
func maskPhone(_ string) string { return "+7******####" }
func maskCard(s string) string {
	d := onlyDigits(s)
	if len(d) < 4 {
		return "************"
	}
	return "************" + d[len(d)-4:]
}
func maskINN(s string) string {
	d := onlyDigits(s)
	if len(d) <= 2 {
		return "**"
	}
	return strings.Repeat("*", len(d)-2) + d[len(d)-2:]
}