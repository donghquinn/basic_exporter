# ============ BASE ===========
FROM golang:1.24rc2-alpine3.21 as base

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# ========= BUILDER ==========
FROM base as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o network_status_exporter .

# ========= RUNNER ==========
FROM golang:1.24rc2-alpine3.21 as release

WORKDIR /home/node

COPY --from=builder /app/network_status_exporter ./network_status_exporter

EXPOSE $APP_PORT

CMD [ "./network_status_exporter" ]