FROM golang:1.13.0-alpine3.10
WORKDIR /go/src/github.com/grupozap/aegir
ADD . .
RUN apk add --no-cache gcc git make musl-dev tar curl ca-certificates && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && \
    dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -o aegir .
ENTRYPOINT ["./aegir"]
