# Base image for runtime
FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS base
USER app
WORKDIR /app
EXPOSE 8084

# Build stage
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
ARG BUILD_CONFIGURATION=Release
WORKDIR /src

# Kopirajte csproj fajl iz lokalnog direktorijuma u /src direktorijum unutar kontejnera
COPY ["Tours.csproj", "./"]

# Restore dependencies
RUN dotnet restore "Tours.csproj"

# Kopirajte sve ostale fajlove iz vašeg lokalnog direktorijuma u kontejner
COPY . .

# Build the project
WORKDIR "/src"
RUN dotnet build "Tours.csproj" -c $BUILD_CONFIGURATION -o /app/build

# Publish the project
FROM build AS publish
RUN dotnet publish "Tours.csproj" -c $BUILD_CONFIGURATION -o /app/publish /p:UseAppHost=false

# Final stage: runtime image
FROM base AS final
ENV ASPNETCORE_URLS=http://+:8084

WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "Tours.dll"]


