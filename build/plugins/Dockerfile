FROM golang:alpine

RUN apk add build-base

WORKDIR $GOPATH/src/app

COPY internal/runtime/operators .

RUN go build -o plugins/and.so -buildmode=plugin internal/runtime/operators/and/main.go
RUN go build -o plugins/filter.so -buildmode=plugin internal/runtime/operators/filter/main.go
RUN go build -o plugins/more.so -buildmode=plugin internal/runtime/operators/more/main.go
RUN go build -o plugins/mul.so -buildmode=plugin internal/runtime/operators/mul/main.go
RUN go build -o plugins/or.so -buildmode=plugin internal/runtime/operators/or/main.go
RUN go build -o plugins/remainder.so -buildmode=plugin internal/runtime/operators/remainder/main.go

