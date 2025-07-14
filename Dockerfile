FROM golang:1.23.9

WORKDIR /go/src

ENV TIME_ZONE Asia/Shanghai

RUN ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime

COPY ./cmd/api/app /go/src/app

# COPY ./web/dist /go/src/web