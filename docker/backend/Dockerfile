FROM golang:1.26-alpine AS build

WORKDIR /app

COPY /src ./

RUN go build

FROM alpine:3.24

COPY --from=build /app/books /usr/local/bin

CMD ["/usr/local/bin/books"]
