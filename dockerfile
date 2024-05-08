FROM golang:1.21.5-bullseye AS build

RUN apt-get update

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd

RUN go build -o search-service

FROM busybox:latest

WORKDIR /search-service/cmd

COPY --from=build /app/cmd/search-service .

COPY --from=build /app/.env /search-service

EXPOSE 4003

CMD ["./search-service"]