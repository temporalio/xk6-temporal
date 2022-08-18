FROM golang:1.18 AS builder

WORKDIR /usr/src/k6

RUN go install go.k6.io/xk6/cmd/xk6@latest

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN xk6 build --output /usr/local/bin/k6 \
    --with github.com/grafana/xk6-output-prometheus-remote \
    --with github.com/temporalio/xk6-temporal=.

FROM scratch

COPY --from=builder /usr/local/bin/k6 /usr/local/bin/k6

CMD ["k6"]