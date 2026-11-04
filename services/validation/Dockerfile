FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

COPY PharmacyValidation.fsproj ./
RUN dotnet restore -p:TargetFramework=net8.0

COPY . ./
RUN dotnet publish PharmacyValidation.fsproj -c Release -p:TargetFramework=net8.0 -o /app/publish

FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS final
WORKDIR /app
COPY --from=build /app/publish .

EXPOSE 5005
ENV PORT=5005
ENV ASPNETCORE_URLS=http://+:5005

ENTRYPOINT ["dotnet", "PharmacyValidation.dll"]
