FROM golang:1.6-alpine

COPY . /go/src/app

WORKDIR /go/src/app

RUN apk --no-cache add curl git && \
  go get ./ && \
  go build && go install -v

expose 8080 9020

CMD ["app"]
