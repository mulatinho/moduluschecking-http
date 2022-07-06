FROM golang:1.18
ENV GOPATH /opt/go
ENV PATH $PATH:/opt/go/bin
ENV GIN_MODE=release

WORKDIR /opt/modulus
COPY main.go go.* /opt/modulus/

RUN groupadd modulus -g 800 && \
    useradd -g 800 -u 800 -d /opt/modulus modulus && \
    chown -R modulus.modulus /opt

USER modulus
RUN go mod download -x
COPY data/*.txt /opt/go/src/github.com/AntoineAugusti/moduluschecking/data/

CMD [ "go", "run", "main.go" ]
