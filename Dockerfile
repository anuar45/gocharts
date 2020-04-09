### Go Build stage

FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /go/src/gocharts

COPY . .

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/gocharts . 


### Image Build stage

FROM alpine

RUN apk update && apk add --no-cache ca-certificates

COPY --from=builder /go/bin/gocharts /go/bin/gocharts
COPY --from=builder /go/src/gocharts/static /go/bin/

WORKDIR /go/bin

EXPOSE 80/tcp

ENTRYPOINT ["./gocharts"]