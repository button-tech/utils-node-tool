FROM golang:1.11-alpine

RUN mkdir /app
ADD . /app/

WORKDIR /app/

ENV GIN_MODE=release


RUN go build -o main ./server/main

EXPOSE 8545

CMD ["/app/main"]