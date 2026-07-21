FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /temperature-api .

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /temperature-api /usr/local/bin/temperature-api

EXPOSE 8081
USER nobody:nobody
ENTRYPOINT ["/usr/local/bin/temperature-api"]
