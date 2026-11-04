namespace PharmacyValidation.Domain

open System
open System.Text.Json.Serialization

[<CLIMutable>]
type RawMedicine = {
    [<JsonPropertyName("name")>]
    name: string option
    [<JsonPropertyName("dosage")>]
    dosage: string option
    [<JsonPropertyName("quantity")>]
    quantity: int option
}

[<CLIMutable>]
type RawPrescription = {
    [<JsonPropertyName("patientName")>]
    patientName: string option
    [<JsonPropertyName("doctorName")>]
    doctorName: string option
    [<JsonPropertyName("doctorRegNo")>]
    doctorRegNo: string option
    [<JsonPropertyName("date")>]
    date: string option
    [<JsonPropertyName("medicines")>]
    medicines: RawMedicine list option
}

type ValidMedicine = {
    Name: string
    Dosage: string
    Quantity: int
}

type Prescription = {
    PatientName: string
    DoctorName: string
    DoctorRegNo: string
    Date: DateTime
    Medicines: ValidMedicine list
}

[<CLIMutable>]
type FieldIssue = {
    [<JsonPropertyName("field")>]
    Field: string
    [<JsonPropertyName("reason")>]
    Reason: string
}

type ValidationResult =
    | Valid of Prescription
    | Incomplete of missing: string list
    | Invalid of FieldIssue list

[<CLIMutable>]
type ValidationResponse = {
    [<JsonPropertyName("status")>]
    status: string
    [<JsonPropertyName("missing")>]
    missing: string list
    [<JsonPropertyName("issues")>]
    issues: FieldIssue list
}
