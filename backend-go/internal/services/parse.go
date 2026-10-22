package services

import (
	"regexp"
	"strconv"
	"strings"

	"smart-pharmacy/backend-go/internal/models"
)

var (
	patientRegex = regexp.MustCompile(`(?i)(?:patient(?:\s+name)?|patient:|name:)\s*([A-Za-z\s\.\'-]+)`)
	doctorRegex  = regexp.MustCompile(`(?i)(?:dr\.|doctor(?:\s+name)?|physician:?)\s*([A-Za-z\s\.\'-]+)`)
	regNoRegex   = regexp.MustCompile(`(?i)(?:reg(?:istration)?(?:\.|\s*no|\s*#)?|license(?:\s*no)?):\s*([A-Za-z0-9\-]+)`)
	dateRegex    = regexp.MustCompile(`(?i)(?:date:?)\s*(\d{4}[-/.]\d{1,2}[-/.]\d{1,2}|\d{1,2}[-/.]\d{1,2}[-/.]\d{4})`)

	// Line items for medicines: e.g. "Amoxicillin 500mg - 21" or "Tab Paracetamol 500mg Qty: 10" or "Metformin 500mg x 30"
	medItemRegex = regexp.MustCompile(`(?i)(?:(?:Rx|Tab|Cap)\s+)?([A-Za-z\s\-]+?)\s+(\d+(?:\.\d+)?\s*(?:mg|mcg|g|ml|IU|%))\s*(?:(?:x|qty|#|count|take)?\s*[:\-]?\s*(\d+))?`)
)

// ParsePrescriptionText extracts structured raw fields from OCR plain text
func ParsePrescriptionText(text string) models.ParsedRawPrescription {
	var result models.ParsedRawPrescription

	lines := strings.Split(text, "\n")

	// 1. Extract Patient Name
	for _, line := range lines {
		if m := patientRegex.FindStringSubmatch(line); len(m) > 1 {
			val := strings.TrimSpace(m[1])
			if val != "" {
				result.PatientName = &val
				break
			}
		}
	}

	// 2. Extract Doctor Name
	for _, line := range lines {
		if m := doctorRegex.FindStringSubmatch(line); len(m) > 1 {
			val := strings.TrimSpace(m[1])
			if !strings.HasPrefix(strings.ToLower(val), "dr.") {
				val = "Dr. " + val
			}
			result.DoctorName = &val
			break
		}
	}

	// 3. Extract Doctor Reg No
	for _, line := range lines {
		if m := regNoRegex.FindStringSubmatch(line); len(m) > 1 {
			val := strings.TrimSpace(m[1])
			if val != "" {
				result.DoctorRegNo = &val
				break
			}
		}
	}

	// 4. Extract Date
	for _, line := range lines {
		if m := dateRegex.FindStringSubmatch(line); len(m) > 1 {
			rawDate := strings.TrimSpace(m[1])
			normalized := normalizeDateString(rawDate)
			result.Date = &normalized
			break
		}
	}

	// 5. Extract Medicines
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if patientRegex.MatchString(line) || doctorRegex.MatchString(line) || regNoRegex.MatchString(line) || dateRegex.MatchString(line) {
			continue
		}

		if m := medItemRegex.FindStringSubmatch(line); len(m) > 2 {
			name := strings.TrimSpace(m[1])
			dosage := strings.TrimSpace(m[2])
			if name != "" && dosage != "" && !strings.EqualFold(name, "patient") && !strings.EqualFold(name, "date") {
				var qty *int
				if len(m) > 3 && m[3] != "" {
					if q, err := strconv.Atoi(m[3]); err == nil && q > 0 {
						qty = &q
					}
				}
				if qty == nil {
					defaultQty := 10
					qty = &defaultQty
				}

				result.Medicines = append(result.Medicines, models.ParsedRawMedicine{
					Name:     &name,
					Dosage:   &dosage,
					Quantity: qty,
				})
			}
		}
	}

	return result
}

func normalizeDateString(d string) string {
	d = strings.ReplaceAll(d, "/", "-")
	d = strings.ReplaceAll(d, ".", "-")
	parts := strings.Split(d, "-")
	if len(parts) == 3 {
		if len(parts[2]) == 4 {
			return parts[2] + "-" + padZero(parts[0]) + "-" + padZero(parts[1])
		}
		if len(parts[0]) == 4 {
			return parts[0] + "-" + padZero(parts[1]) + "-" + padZero(parts[2])
		}
	}
	return d
}

func padZero(s string) string {
	if len(s) == 1 {
		return "0" + s
	}
	return s
}
