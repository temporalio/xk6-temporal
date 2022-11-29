FROM golang:1.18 AS builder

WORKDIR /usr/src/k6

RUN go install go.k6.io/xk6/cmd/xk6@latest

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN xk6 build --output /usr/local/bin/k6 \
    --with github.com/grafana/xk6-output-prometheus-remote \
    --with github.com/temporalio/xk6-prometheus-client \
    --with github.com/temporalio/xk6-temporal=.

FROM alpine:3.16

WORKDIR /opt/k6

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/local/bin/k6 /usr/local/bin/k6
COPY ./examples /opt/k6/examples

CMD ["k6"]
