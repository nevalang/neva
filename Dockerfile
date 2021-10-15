FROM golang:alpine

RUN apk add build-base

WORKDIR $GOPATH/src/app

COPY . .

RUN go build -o plugins/mul.so -buildmode=plugin internal/runtime/operators/mul/main.go
RUN go build -o /bin/nevabin cmd/server/main.go

EXPOSE 8090

RUN ["chmod", "+x", "/bin/nevabin"]

CMD ["/bin/nevabin"]
