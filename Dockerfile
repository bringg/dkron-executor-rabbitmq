# build dkron-executor-rabbitmq
FROM golang:1.10-alpine AS builder

ARG dkron_version=0.10.2

# install packages and download binaries
RUN apk --no-cache add git \
    && wget -O - https://raw.githubusercontent.com/golang/dep/master/install.sh | sh \
    && wget -O - https://github.com/victorcoder/dkron/releases/download/v${dkron_version}/dkron_${dkron_version}_linux_amd64.tar.gz | tar xzf - \
    && mv dkron_${dkron_version}_linux_amd64 /tmp/dkron

# ensure dependencies
WORKDIR $GOPATH/src/github.com/bringg/dkron-executor-rabbitmq
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only -v

# compile and install
COPY . ./
RUN go install -v github.com/bringg/dkron-executor-rabbitmq

# final image
FROM alpine
LABEL maintainer "Bringg Devops <devops@bringg.com>"

COPY --from=builder /tmp/dkron /opt/dkron
COPY --from=builder /go/bin/dkron-executor-rabbitmq /opt/dkron

ENV PATH=/opt/dkron:$PATH
EXPOSE 8080 8946
CMD ["dkron"]
