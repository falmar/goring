FROM golang:1.12-alpine as build-env

WORKDIR /go-app
ADD . /go-app

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -o goring .

FROM alpine:3.9
WORKDIR /go-app
COPY --from=build-env /go-app/goring /go-app/goring
COPY public /go-app/public
COPY tpl /go-app/tpl
ENTRYPOINT ["./goring"]
