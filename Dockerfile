FROM golang:1.10-alpine AS builder

ARG dep_version=0.4.1
ARG dkron_branch=v0.10.0

RUN apk --no-cache add git \
    && wget https://github.com/golang/dep/releases/download/v${dep_version}/dep-linux-amd64 -O /usr/local/bin/dep \
    && chmod +x /usr/local/bin/dep \
    # build dkron
    && git clone https://github.com/victorcoder/dkron -b $dkron_branch $GOPATH/src/github.com/victorcoder/dkron \
    && cd $GOPATH/src/github.com/victorcoder/dkron \
    && dep ensure -v \
    && go install -v github.com/victorcoder/dkron

# build dkron-executor-rabbitmq
COPY . $GOPATH/src/github.com/bringg/dkron-executor-rabbitmq
RUN cd $GOPATH/src/github.com/bringg/dkron-executor-rabbitmq \
    && dep ensure -v \
    && go install -v github.com/bringg/dkron-executor-rabbitmq

# final image
FROM alpine
LABEL maintainer "Bringg Devops <devops@bringg.com>"

EXPOSE 8080 8946
CMD ["/usr/local/bin/dkron"]

COPY --from=builder /go/bin/dkron /go/bin/dkron-executor-rabbitmq /usr/local/bin/
