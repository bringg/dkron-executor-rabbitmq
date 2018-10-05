# build dkron-executor-rabbitmq
FROM golang:1.11.1-alpine AS builder

# install packages and download binaries
RUN apk --no-cache add git \
    && wget -O - https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# ensure dependencies
WORKDIR $GOPATH/src/github.com/bringg/dkron-executor-rabbitmq
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only -v

# compile and install
COPY . ./
RUN go install -v github.com/bringg/dkron-executor-rabbitmq

# final image
FROM dkron/dkron:v0.10.2
LABEL maintainer="Bringg Devops <devops@bringg.com>"

COPY --from=builder /go/bin/dkron-executor-rabbitmq /opt/local/dkron
