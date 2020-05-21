### Go Build stage

FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /go/src/topgomods

COPY . .

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.VERSION=$(git describe --tags)" -o /go/bin/topgomods cmd/topgomods/main.go

### Image Build stage

FROM alpine

RUN apk update && apk add --no-cache ca-certificates

COPY --from=builder /go/bin/topgomods /go/bin/topgomods
COPY --from=builder /go/src/topgomods/static /go/bin/static/

WORKDIR /go/bin

EXPOSE 8080/tcp

ENTRYPOINT ["./topgomods"]