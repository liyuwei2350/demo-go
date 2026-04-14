FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
ENV GO111MODULE=on
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o app .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /app/app /
EXPOSE 8080
ENTRYPOINT ["/app"]