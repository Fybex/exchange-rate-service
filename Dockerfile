FROM golang:1.22

WORKDIR /app

COPY . .

RUN go install -mod=mod github.com/githubnemo/CompileDaemon
RUN go mod download

EXPOSE 8000

ENTRYPOINT CompileDaemon --build="go build ./cmd/exchange-rate-service/main.go" --command="./main"