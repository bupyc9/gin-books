FROM golang:1.26-alpine AS build

WORKDIR /app

COPY /src ./

RUN go build -o /usr/local/bin

CMD ["/usr/local/bin/books"]
