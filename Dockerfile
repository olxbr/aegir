FROM golang:1.15.2-alpine AS builder
WORKDIR /go/src/github.com/grupozap/aegir
RUN apk add gcc git make musl-dev tar curl ca-certificates bash
ADD . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o aegir .


FROM alpine:latest AS dry-app
RUN mkdir /aegir \
 &&  mkdir /aegir/tls
WORKDIR /aegir
COPY --from=builder /go/src/github.com/grupozap/aegir/aegir .
ENTRYPOINT ["./aegir"]
