FROM golang:1.18 as temporalite

RUN git clone -b rh-metrics-port https://github.com/robholland/temporalite
RUN cd temporalite && go install ./cmd/temporalite
EXPOSE 7233
EXPOSE 8233

ENTRYPOINT ["temporalite", "start", "--ephemeral", "-n", "default", "--ip" , "0.0.0.0", "--metrics-port", "9000"]

FROM golang:1.18 as worker

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY worker worker

RUN go build -v -o /usr/local/bin/worker ./worker

ENTRYPOINT ["worker"]
