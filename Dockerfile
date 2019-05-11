FROM golang:latest AS builder

ARG DIR

RUN mkdir /build
ADD . /build
WORKDIR /build

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/main $DIR

FROM scratch
COPY --from=builder /build/bin /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080

CMD ["/app/main"]