namespace PharmacyValidation.Api

open Microsoft.AspNetCore.Http
open Giraffe
open PharmacyValidation.Domain
open PharmacyValidation.Rules

module Handlers =

    let healthHandler: HttpHandler =
        fun (next: HttpFunc) (ctx: HttpContext) ->
            json {| status = "healthy" |} next ctx

    let validateHandler: HttpHandler =
        fun (next: HttpFunc) (ctx: HttpContext) ->
            task {
                try
                    let! raw = ctx.BindJsonAsync<RawPrescription>()
                    let result = PrescriptionRules.validatePrescription raw
                    let response = PrescriptionRules.toResponse result
                    return! json response next ctx
                with ex ->
                    let errResponse = {
                        status = "invalid"
                        missing = []
                        issues = [ { Field = "body"; Reason = sprintf "Malformed JSON or invalid request format: %s" ex.Message } ]
                    }
                    return! json errResponse next ctx
            }

    let webApp: HttpHandler =
        choose [
            GET >=> choose [
                route "/health" >=> healthHandler
            ]
            POST >=> choose [
                route "/validate/prescription" >=> validateHandler
            ]
            setStatusCode 404 >=> text "Not Found"
        ]
