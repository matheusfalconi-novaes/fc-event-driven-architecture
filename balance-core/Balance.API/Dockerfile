FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS base
WORKDIR /app

FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
ARG BUILD_CONFIGURATION=Release
WORKDIR /src

COPY ["Balance.Shared/*.csproj", "Balance.Shared/"]
RUN dotnet restore "Balance.Shared/Balance.Shared.csproj"

COPY ["Balance.Domain/*.csproj", "Balance.Domain/"]
RUN dotnet restore "Balance.Domain/Balance.Domain.csproj"

COPY ["Balance.Infrastructure/*.csproj", "Balance.Infrastructure/"]
RUN dotnet restore "Balance.Infrastructure/Balance.Infrastructure.csproj"

COPY ["Balance.API/*.csproj", "Balance.API/"]
RUN dotnet restore "Balance.API/Balance.API.csproj"

COPY Balance.Shared/. ./Balance.Shared/
COPY Balance.Domain/. ./Balance.Domain/
COPY Balance.Infrastructure/. ./Balance.Infrastructure/
COPY Balance.API/. ./Balance.API/

WORKDIR "/src/Balance.API"
RUN dotnet build "Balance.API.csproj" -c $BUILD_CONFIGURATION -o /app/build

FROM build AS publish
ARG BUILD_CONFIGURATION=Release
RUN dotnet publish "Balance.API.csproj" -c $BUILD_CONFIGURATION -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "Balance.API.dll"]
