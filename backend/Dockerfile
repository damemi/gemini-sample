# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /gemini-backend .

FROM alpine:3.20
RUN adduser -D -u 65532 app
USER 65532:65532
WORKDIR /app
COPY --from=build /gemini-backend ./
EXPOSE 8080
ENV LISTEN_ADDR=:8080
ENTRYPOINT ["/app/gemini-backend"]
