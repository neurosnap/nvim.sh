FROM golang:alpine

WORKDIR /app

ADD . /app

RUN CGO_ENABLED=0 go build -o nvim

EXPOSE 80

CMD ["/app/nvim"]
