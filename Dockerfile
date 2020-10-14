FROM golang:1.15.2-alpine AS builder
WORKDIR /go/src/github.com/grupozap/aegir
RUN apk add --no-cache ca-certificates git
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v
RUN grep nobody /etc/passwd > /etc/passwd.nobody

FROM scratch
COPY --from=builder /etc/passwd.nobody /etc/passwd
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/github.com/grupozap/aegir/aegir /aegir
USER nobody
EXPOSE 8443
ENTRYPOINT ["/aegir"]
