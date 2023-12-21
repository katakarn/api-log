FROM golang:1.19-alpine3.18 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main ./main.go

FROM debian:10

WORKDIR /dist

COPY --from=build /app/main /

EXPOSE 8080

ENTRYPOINT ["/main"]
