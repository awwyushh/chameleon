FROM golang:latest as build

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o honeypot main.go dynamic.go

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /app/honeypot /app/honeypot

EXPOSE 8080
CMD ["/app/honeypot"]
