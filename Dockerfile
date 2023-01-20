FROM golang:1.19.1-alpine AS builder

RUN apk add --no-cache --update git gcc g++ bash openssh tzdata

ENV GOPATH=/go

COPY . $GOPATH/src/github.com/tashfi04/printbin-server
WORKDIR $GOPATH/src/github.com/tashfi04/printbin-server

RUN chmod +x ./build.sh
RUN ./build.sh

RUN mv ./printbin-server /go/bin/
RUN mv ./config.yml /go/bin/
RUN mv ./config-secret.yml /go/bin/

FROM alpine:latest
RUN apk add --no-cache --update ca-certificates openssl && apk add --no-cache tzdata
COPY --from=0 /go/bin/printbin-server /usr/local/bin/printbin-server
COPY --from=0 /go/bin/config.yml /usr/local/bin/config.yml
COPY --from=0 /go/bin/config-secret.yml /usr/local/bin/config-secret.yml

ENV PRINTBIN_CONFIG_NAME config
ENV PRINTBIN_CONFIG_SECRET_NAME config-secret
ENV PRINTBIN_CONFIG_PATH /usr/local/bin/

CMD ["printbin-server", "serve", "-p", "8001"]