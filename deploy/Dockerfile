FROM golang:1.21.1

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o main cmd/app/main.go

CMD ["/app/main"]