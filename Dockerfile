FROM golang:1.15.2-alpine AS builder
WORKDIR /go/src/github.com/grupozap/aegir
RUN apk add --no-cache ca-certificates git
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v

FROM alpine:latest AS dry-app
RUN mkdir /aegir \
 &&  mkdir /aegir/tls
WORKDIR /aegir
COPY --from=builder /go/src/github.com/grupozap/aegir/aegir .
ENTRYPOINT ["./aegir"]
