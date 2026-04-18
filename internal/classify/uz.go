package classify

import "pii-scanner/internal/model"

type UZConfig struct {
	UZ2GovIDThreshold   int
	UZ2PaymentThreshold int
	UZ3OrdinaryThreshold int
}

func DefaultUZConfig() UZConfig {
	return UZConfig{
		UZ2GovIDThreshold:   20,
		UZ2PaymentThreshold: 20,
		UZ3OrdinaryThreshold: 50,
	}
}

func DetermineUZWithReasons(findings []model.Finding, cfg UZConfig) (string, []string) {
	var ordinary, govID, payment, special, biometric int

	for _, f := range findings {
		switch f.Category {
		case "ordinary":
			ordinary += f.Count
		case "gov_id":
			govID += f.Count
		case "payment":
			payment += f.Count
		case "special":
			special += f.Count
		case "biometric":
			biometric += f.Count
		}
	}

	reasons := make([]string, 0, 3)

	if special > 0 {
		reasons = append(reasons, "обнаружены специальные категории ПДн")
	}
	if biometric > 0 {
		reasons = append(reasons, "обнаружены биометрические данные")
	}
	if special > 0 || biometric > 0 {
		return "UZ-1", reasons
	}

	if payment >= cfg.UZ2PaymentThreshold {
		reasons = append(reasons, "платежные данные превышают порог UZ-2")
	}
	if govID >= cfg.UZ2GovIDThreshold {
		reasons = append(reasons, "госидентификаторы превышают порог UZ-2")
	}
	if len(reasons) > 0 {
		return "UZ-2", reasons
	}

	if govID > 0 {
		reasons = append(reasons, "обнаружены госидентификаторы в небольшом объеме")
	}
	if ordinary >= cfg.UZ3OrdinaryThreshold {
		reasons = append(reasons, "обычные ПДн превышают порог UZ-3")
	}
	if len(reasons) > 0 {
		return "UZ-3", reasons
	}

	if ordinary > 0 {
		return "UZ-4", []string{"обнаружены только обычные ПДн в небольшом объеме"}
	}

	return "NO_PD", []string{"ПДн не обнаружены"}
}