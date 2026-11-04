package models

type ParsedRawMedicine struct {
	Name     *string `json:"name"`
	Dosage   *string `json:"dosage"`
	Quantity *int    `json:"quantity"`
}

type ParsedRawPrescription struct {
	PatientName *string             `json:"patientName"`
	DoctorName  *string             `json:"doctorName"`
	DoctorRegNo *string             `json:"doctorRegNo"`
	Date        *string             `json:"date"`
	Medicines   []ParsedRawMedicine `json:"medicines"`
}
