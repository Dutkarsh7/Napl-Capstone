module PharmacyValidation.Tests.PrescriptionTests

open System
open Expecto
open PharmacyValidation.Domain
open PharmacyValidation.Rules.PrescriptionRules

let sampleNow = DateTime(2026, 10, 22, 12, 0, 0, DateTimeKind.Utc)

let validSample: RawPrescription = {
    patientName = Some "Alice Smith"
    doctorName = Some "Dr. John Watson"
    doctorRegNo = Some "MED-99482"
    date = Some "2026-10-15"
    medicines = Some [
        { name = Some "Amoxicillin 500mg"; dosage = Some "1 capsule 3x daily"; quantity = Some 15 }
    ]
}

[<Tests>]
let tests =
    testList "PrescriptionValidationRules" [

        testCase "Valid prescription succeeds" <| fun () ->
            let result = validatePrescriptionWithDate sampleNow validSample
            match result with
            | Valid rx ->
                Expect.equal rx.PatientName "Alice Smith" "Patient name must match"
                Expect.equal rx.DoctorName "Dr. John Watson" "Doctor name must match"
                Expect.equal rx.DoctorRegNo "MED-99482" "Doctor reg no must match"
                Expect.equal rx.Medicines.Length 1 "Medicine count must be 1"
            | other ->
                failwithf "Expected Valid, got %A" other

        testCase "Valid multi-medicine prescription" <| fun () ->
            let raw = {
                validSample with
                    medicines = Some [
                        { name = Some "Amoxicillin 500mg"; dosage = Some "1 cap TID"; quantity = Some 21 }
                        { name = Some "Paracetamol 500mg"; dosage = Some "1 tab PRN"; quantity = Some 10 }
                    ]
            }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Valid rx ->
                Expect.equal rx.Medicines.Length 2 "Should contain 2 medicines"
            | other ->
                failwithf "Expected Valid, got %A" other

        testCase "Valid boundary: quantity = 1" <| fun () ->
            let raw = {
                validSample with
                    medicines = Some [ { name = Some "Lisinopril 10mg"; dosage = Some "Daily"; quantity = Some 1 } ]
            }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Valid rx -> Expect.equal rx.Medicines.[0].Quantity 1 "Min quantity 1 is allowed"
            | other -> failwithf "Expected Valid, got %A" other

        testCase "Valid boundary: quantity = 30" <| fun () ->
            let raw = {
                validSample with
                    medicines = Some [ { name = Some "Atorvastatin 20mg"; dosage = Some "At night"; quantity = Some 30 } ]
            }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Valid rx -> Expect.equal rx.Medicines.[0].Quantity 30 "Max quantity 30 is allowed"
            | other -> failwithf "Expected Valid, got %A" other

        testCase "Incomplete: missing patient name" <| fun () ->
            let raw = { validSample with patientName = None }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Incomplete missing ->
                Expect.contains missing "patientName" "Should list patientName as missing"
            | other -> failwithf "Expected Incomplete, got %A" other

        testCase "Incomplete: empty patient name string" <| fun () ->
            let raw = { validSample with patientName = Some "   " }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Incomplete missing ->
                Expect.contains missing "patientName" "Should identify blank string as missing"
            | other -> failwithf "Expected Incomplete, got %A" other

        testCase "Incomplete: missing doctor name" <| fun () ->
            let raw = { validSample with doctorName = None }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Incomplete missing -> Expect.contains missing "doctorName" "Missing doctorName"
            | other -> failwithf "Expected Incomplete, got %A" other

        testCase "Incomplete: missing doctor registration number" <| fun () ->
            let raw = { validSample with doctorRegNo = None }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Incomplete missing -> Expect.contains missing "doctorRegNo" "Missing doctorRegNo"
            | other -> failwithf "Expected Incomplete, got %A" other

        testCase "Incomplete: missing date" <| fun () ->
            let raw = { validSample with date = None }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Incomplete missing -> Expect.contains missing "date" "Missing date"
            | other -> failwithf "Expected Incomplete, got %A" other

        testCase "Incomplete: missing medicines list" <| fun () ->
            let raw = { validSample with medicines = None }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Incomplete missing -> Expect.contains missing "medicines" "Missing medicines list"
            | other -> failwithf "Expected Incomplete, got %A" other

        testCase "Incomplete: empty medicines list" <| fun () ->
            let raw = { validSample with medicines = Some [] }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Incomplete missing -> Expect.contains missing "medicines" "Empty medicines list is incomplete"
            | other -> failwithf "Expected Incomplete, got %A" other

        testCase "Incomplete: multiple missing fields reported together" <| fun () ->
            let raw = { patientName = None; doctorName = None; doctorRegNo = None; date = None; medicines = None }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Incomplete missing ->
                Expect.equal missing.Length 5 "All 5 missing fields should be detected"
            | other -> failwithf "Expected Incomplete, got %A" other

        testCase "Invalid: future date" <| fun () ->
            let raw = { validSample with date = Some "2026-11-01" }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Invalid issues ->
                let hasFuture = issues |> List.exists (fun i -> i.Field = "date" && i.Reason.Contains("cannot be in the future"))
                Expect.isTrue hasFuture "Must reject future prescription date"
            | other -> failwithf "Expected Invalid, got %A" other

        testCase "Invalid: date older than 6 months" <| fun () ->
            let raw = { validSample with date = Some "2025-01-01" }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Invalid issues ->
                let hasOld = issues |> List.exists (fun i -> i.Field = "date" && i.Reason.Contains("older than 6 months"))
                Expect.isTrue hasOld "Must reject expired prescription older than 180 days"
            | other -> failwithf "Expected Invalid, got %A" other

        testCase "Invalid: unparseable date" <| fun () ->
            let raw = { validSample with date = Some "invalid-date-string" }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Invalid issues ->
                let hasFormatErr = issues |> List.exists (fun i -> i.Field = "date" && i.Reason.Contains("Invalid date format"))
                Expect.isTrue hasFormatErr "Must flag invalid date format"
            | other -> failwithf "Expected Invalid, got %A" other

        testCase "Invalid: quantity 0" <| fun () ->
            let raw = {
                validSample with
                    medicines = Some [ { name = Some "Ibuprofen"; dosage = Some "400mg"; quantity = Some 0 } ]
            }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Invalid issues ->
                let hasZero = issues |> List.exists (fun i -> i.Field = "medicines[0].quantity" && i.Reason.Contains("greater than 0"))
                Expect.isTrue hasZero "Must reject quantity 0"
            | other -> failwithf "Expected Invalid, got %A" other

        testCase "Invalid: negative quantity" <| fun () ->
            let raw = {
                validSample with
                    medicines = Some [ { name = Some "Ibuprofen"; dosage = Some "400mg"; quantity = Some -5 } ]
            }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Invalid issues ->
                let hasNeg = issues |> List.exists (fun i -> i.Field = "medicines[0].quantity")
                Expect.isTrue hasNeg "Must reject negative quantity"
            | other -> failwithf "Expected Invalid, got %A" other

        testCase "Invalid: quantity exceeds 30" <| fun () ->
            let raw = {
                validSample with
                    medicines = Some [ { name = Some "Amoxicillin"; dosage = Some "500mg"; quantity = Some 60 } ]
            }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Invalid issues ->
                let hasExcess = issues |> List.exists (fun i -> i.Field = "medicines[0].quantity" && i.Reason.Contains("exceeds maximum 30"))
                Expect.isTrue hasExcess "Must reject quantity > 30"
            | other -> failwithf "Expected Invalid, got %A" other

        testCase "Invalid: blank medicine name" <| fun () ->
            let raw = {
                validSample with
                    medicines = Some [ { name = Some "   "; dosage = Some "500mg"; quantity = Some 10 } ]
            }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Invalid issues ->
                let hasBlankName = issues |> List.exists (fun i -> i.Field = "medicines[0].name" && i.Reason.Contains("cannot be blank"))
                Expect.isTrue hasBlankName "Must reject blank medicine name"
            | other -> failwithf "Expected Invalid, got %A" other

        testCase "Invalid: blank dosage" <| fun () ->
            let raw = {
                validSample with
                    medicines = Some [ { name = Some "Amoxicillin 500mg"; dosage = Some ""; quantity = Some 10 } ]
            }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Invalid issues ->
                let hasBlankDosage = issues |> List.exists (fun i -> i.Field = "medicines[0].dosage" && i.Reason.Contains("cannot be blank"))
                Expect.isTrue hasBlankDosage "Must reject blank dosage"
            | other -> failwithf "Expected Invalid, got %A" other

        testCase "Invalid: restricted keyword (morphine) flags pharmacist attention" <| fun () ->
            let raw = {
                validSample with
                    medicines = Some [ { name = Some "Morphine Sulfate 15mg"; dosage = Some "1 tab q4h"; quantity = Some 20 } ]
            }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Invalid issues ->
                let hasWarning = issues |> List.exists (fun i -> i.Field = "medicines[0].name" && i.Reason.Contains("requires pharmacist attention"))
                Expect.isTrue hasWarning "Must flag restricted opioid keyword"
            | other -> failwithf "Expected Invalid, got %A" other

        testCase "Invalid: restricted keyword (fentanyl) flags pharmacist attention" <| fun () ->
            let raw = {
                validSample with
                    medicines = Some [ { name = Some "Fentanyl 25mcg Patch"; dosage = Some "Change q72h"; quantity = Some 5 } ]
            }
            let result = validatePrescriptionWithDate sampleNow raw
            match result with
            | Invalid issues ->
                let hasWarning = issues |> List.exists (fun i -> i.Field = "medicines[0].name" && i.Reason.Contains("fentanyl"))
                Expect.isTrue hasWarning "Must flag fentanyl"
            | other -> failwithf "Expected Invalid, got %A" other

        testCase "toResponse correctly maps Valid, Incomplete, and Invalid" <| fun () ->
            let validRes = toResponse (Valid {
                PatientName = "A"; DoctorName = "B"; DoctorRegNo = "C"; Date = DateTime.UtcNow; Medicines = []
            })
            Expect.equal validRes.status "valid" "Valid status"
            Expect.isEmpty validRes.missing "Valid missing is empty"
            Expect.isEmpty validRes.issues "Valid issues is empty"

            let incompRes = toResponse (Incomplete [ "patientName"; "date" ])
            Expect.equal incompRes.status "incomplete" "Incomplete status"
            Expect.equal incompRes.missing.Length 2 "Incomplete missing list"

            let invRes = toResponse (Invalid [ { Field = "date"; Reason = "bad" } ])
            Expect.equal invRes.status "invalid" "Invalid status"
            Expect.equal invRes.issues.Length 1 "Invalid issues list"
    ]
