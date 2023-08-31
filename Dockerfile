FROM golang

COPY . /server
WORKDIR /server

RUN go build -o start.exe -v ./cmd/app

CMD [ "./start" ]