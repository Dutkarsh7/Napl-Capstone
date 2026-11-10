FROM mcr.microsoft.com/dotnet/sdk:10.0 AS build
WORKDIR /src

COPY PharmacyValidation.fsproj ./
RUN dotnet restore

COPY . ./
RUN dotnet publish PharmacyValidation.fsproj -c Release -o /app/publish

FROM mcr.microsoft.com/dotnet/aspnet:10.0 AS final
WORKDIR /app
COPY --from=build /app/publish .

EXPOSE 5005
ENV PORT=5005
ENV ASPNETCORE_URLS=http://+:5005

ENTRYPOINT ["dotnet", "PharmacyValidation.dll"]
