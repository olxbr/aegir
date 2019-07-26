FROM golang:1.12.7-alpine
WORKDIR /go/src/github.com/grupozap/aegir
ADD . .
RUN apk add --no-cache gcc git make musl-dev tar curl && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && \
    dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -o aegir .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /aegir
COPY --from=0 /go/src/github.com/grupozap/aegir/aegir .
ENTRYPOINT ["./aegir"]
