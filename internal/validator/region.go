package validator

import "strings"

var validISORegions = map[string]bool{
	"ID": true,
	"US": true,
	"GB": true,
	"AU": true,
	"SG": true,
	"MY": true,
	"JP": true,
	"KR": true,
	"IN": true,
	"DE": true,
	"FR": true,
	"CA": true,
	"BR": true,
	"NL": true,
	"SE": true,
}

func NormalizeRegion(region string) string {
	normalized := strings.ToUpper(strings.TrimSpace(region))
	if normalized == "" {
		return "ID"
	}

	if _, ok := validISORegions[normalized]; !ok {
		return "ID"
	}

	return normalized
}
