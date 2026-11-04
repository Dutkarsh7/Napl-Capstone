namespace PharmacyValidation.Rules

open System
open System.Globalization
open PharmacyValidation.Domain

module PrescriptionRules =

    let restrictedKeywords = [
        "morphine"
        "fentanyl"
        "oxycodone"
        "ketamine"
        "methadone"
        "codeine"
        "controlled"
        "hydrocodone"
        "hydromorphone"
    ]

    let isBlank (s: string option) =
        match s with
        | None -> true
        | Some str when String.IsNullOrWhiteSpace(str) -> true
        | _ -> false

    let checkMissingFields (raw: RawPrescription) : string list =
        let mutable missing = []
        if isBlank raw.patientName then missing <- "patientName" :: missing
        if isBlank raw.doctorName then missing <- "doctorName" :: missing
        if isBlank raw.doctorRegNo then missing <- "doctorRegNo" :: missing
        if isBlank raw.date then missing <- "date" :: missing
        match raw.medicines with
        | None -> missing <- "medicines" :: missing
        | Some meds when List.isEmpty meds -> missing <- "medicines" :: missing
        | _ -> ()
        List.rev missing

    let parseDate (dateStr: string) : DateTime option =
        let formats = [| "yyyy-MM-dd"; "yyyy/MM/dd"; "MM/dd/yyyy"; "dd/MM/yyyy"; "yyyy-MM-ddTHH:mm:ssZ"; "o" |]
        match DateTime.TryParseExact(dateStr.Trim(), formats, CultureInfo.InvariantCulture, DateTimeStyles.AssumeUniversal ||| DateTimeStyles.AdjustToUniversal) with
        | true, dt -> Some dt
        | false, _ ->
            match DateTime.TryParse(dateStr.Trim(), CultureInfo.InvariantCulture, DateTimeStyles.AssumeUniversal ||| DateTimeStyles.AdjustToUniversal) with
            | true, dt -> Some dt
            | false, _ -> None

    let validatePrescriptionWithDate (now: DateTime) (raw: RawPrescription) : ValidationResult =
        let missing = checkMissingFields raw
        if not (List.isEmpty missing) then
            Incomplete missing
        else
            let mutable issues: FieldIssue list = []
            
            // Validate Date
            let dateStr = raw.date.Value.Trim()
            let parsedDateOpt = parseDate dateStr
            let validDate =
                match parsedDateOpt with
                | None ->
                    issues <- { Field = "date"; Reason = sprintf "Invalid date format: '%s'. Expected YYYY-MM-DD" dateStr } :: issues
                    None
                | Some dt ->
                    let today = now.Date
                    let rxDate = dt.Date
                    if rxDate > today then
                        issues <- { Field = "date"; Reason = sprintf "Prescription date cannot be in the future (%s)" (rxDate.ToString("yyyy-MM-dd")) } :: issues
                        None
                    elif (today - rxDate).TotalDays > 180.0 then
                        issues <- { Field = "date"; Reason = sprintf "Prescription is older than 6 months (issued %s)" (rxDate.ToString("yyyy-MM-dd")) } :: issues
                        None
                    else
                        Some dt

            // Validate Medicines
            let rawMeds = raw.medicines.Value
            let mutable validMeds: ValidMedicine list = []

            rawMeds
            |> List.iteri (fun idx med ->
                let medIdxName = sprintf "medicines[%d]" idx
                
                // Name check
                let medNameOpt =
                    match med.name with
                    | None ->
                        issues <- { Field = sprintf "%s.name" medIdxName; Reason = "Medicine name cannot be blank" } :: issues
                        None
                    | Some s when String.IsNullOrWhiteSpace(s) ->
                        issues <- { Field = sprintf "%s.name" medIdxName; Reason = "Medicine name cannot be blank" } :: issues
                        None
                    | Some s ->
                        let trimmed = s.Trim()
                        let lower = trimmed.ToLowerInvariant()
                        for keyword in restrictedKeywords do
                            if lower.Contains(keyword) then
                                issues <- { Field = sprintf "%s.name" medIdxName; Reason = sprintf "Controlled substance '%s': requires pharmacist attention" keyword } :: issues
                        Some trimmed

                // Dosage check
                let dosageOpt =
                    match med.dosage with
                    | None ->
                        issues <- { Field = sprintf "%s.dosage" medIdxName; Reason = "Dosage instructions cannot be blank" } :: issues
                        None
                    | Some s when String.IsNullOrWhiteSpace(s) ->
                        issues <- { Field = sprintf "%s.dosage" medIdxName; Reason = "Dosage instructions cannot be blank" } :: issues
                        None
                    | Some s -> Some (s.Trim())

                // Quantity check
                let quantityOpt =
                    match med.quantity with
                    | None ->
                        issues <- { Field = sprintf "%s.quantity" medIdxName; Reason = "Quantity is missing" } :: issues
                        None
                    | Some q when q <= 0 ->
                        issues <- { Field = sprintf "%s.quantity" medIdxName; Reason = sprintf "Quantity must be greater than 0, got %d" q } :: issues
                        None
                    | Some q when q > 30 ->
                        issues <- { Field = sprintf "%s.quantity" medIdxName; Reason = sprintf "Quantity exceeds maximum 30 units per line, got %d" q } :: issues
                        None
                    | Some q -> Some q

                match (medNameOpt, dosageOpt, quantityOpt) with
                | (Some n, Some d, Some q) ->
                    validMeds <- { Name = n; Dosage = d; Quantity = q } :: validMeds
                | _ -> ()
            )

            if not (List.isEmpty issues) then
                Invalid (List.rev issues)
            else
                match validDate with
                | Some dt ->
                    Valid {
                        PatientName = raw.patientName.Value.Trim()
                        DoctorName = raw.doctorName.Value.Trim()
                        DoctorRegNo = raw.doctorRegNo.Value.Trim()
                        Date = dt
                        Medicines = List.rev validMeds
                    }
                | None ->
                    Invalid (List.rev issues)

    let validatePrescription (raw: RawPrescription) : ValidationResult =
        validatePrescriptionWithDate DateTime.UtcNow raw

    let toResponse (result: ValidationResult) : ValidationResponse =
        match result with
        | Valid _ ->
            { status = "valid"; missing = []; issues = [] }
        | Incomplete missingFields ->
            { status = "incomplete"; missing = missingFields; issues = [] }
        | Invalid issuesList ->
            { status = "invalid"; missing = []; issues = issuesList }
