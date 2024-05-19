# Exchange Rate Service

**This project was developed as the second stage of selection for Genesis & KMA SOFTWARE ENGINEERING SCHOOL 4.0**

Service that provides USD to UAN exchange rate and notifies users daily about the rate changes.

## Installing

1. Clone the repository

```sh
git clone https://github.com/Fybex/exchange-rate-service.git
```

2. Navigate to the project directory

```sh
cd exchange-rate-service
```

3. Copy and edit `.env.example` to `.env`, setting actual values for the SMTP server, and database credentials.

```sh
cp .env.example .env
```

4. Run docker compose

```sh
docker-compose up -d
```

The service will be available at `http://localhost:8000`

## API Endpoints

`pkg/api/swagger.yaml` contains the API documentation.

## Running unit tests

```sh
docker-compose exec -it web go test ./... -v
```
