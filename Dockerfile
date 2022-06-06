FROM golang:1.18.2-alpine3.15 as buildenv

USER root

RUN mkdir /temp
WORKDIR /temp

RUN go mod init base
RUN go mod edit -require=github.com/VasKho/yagl@0.1.0
RUN go mod download

RUN rm -rf /temp
WORKDIR /go

CMD ["sh"]
