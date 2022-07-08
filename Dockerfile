FROM golang:1.18 as temporalite

RUN git clone -b rh-shard https://github.com/robholland/temporalite
RUN cd temporalite && go install ./cmd/temporalite
EXPOSE 7233
EXPOSE 8233

ENTRYPOINT ["temporalite", "start", "--ephemeral", "-n", "default", "--ip" , "0.0.0.0", "--metrics-port", "9000"]