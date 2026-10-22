namespace PharmacyValidation

open System
open Microsoft.AspNetCore.Builder
open Microsoft.AspNetCore.Hosting
open Microsoft.Extensions.Hosting
open Microsoft.Extensions.DependencyInjection
open Giraffe
open PharmacyValidation.Api

module Program =

    [<EntryPoint>]
    let main args =
        let builder = WebApplication.CreateBuilder(args)

        let port =
            match Environment.GetEnvironmentVariable("VALIDATION_PORT") with
            | null | "" ->
                match Environment.GetEnvironmentVariable("PORT") with
                | null | "" -> "5005"
                | p -> p
            | p -> p

        builder.WebHost.UseUrls(sprintf "http://0.0.0.0:%s" port) |> ignore

        builder.Services.AddGiraffe() |> ignore

        let app = builder.Build()

        app.UseGiraffe(Handlers.webApp)

        Console.WriteLine(sprintf "F# Prescription Validation Service listening on port %s" port)
        app.Run()
        0
