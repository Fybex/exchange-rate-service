# Exchange Rate Service

Service that provides exchange rate information and subscription functionality.

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
