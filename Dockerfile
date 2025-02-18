# ============ BASE ===========
FROM golang:1.24rc2-alpine3.21 AS base

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# ========= BUILDER ==========
FROM base AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o network_status_exporter .

# ========= RUNNER ==========
FROM golang:1.24rc2-alpine3.21 AS release

WORKDIR /home/node

COPY --from=builder /app/network_status_exporter ./network_status_exporter

EXPOSE 8080

CMD [ "./network_status_exporter" ]