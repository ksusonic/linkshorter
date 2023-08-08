FROM golang:latest
RUN mkdir /app

ENV GO111MODULE=on
ADD . /app/
WORKDIR /app

RUN go build -o main github.com/ksusonic/linkshorter/cmd

EXPOSE 8080
CMD ["/app/main"]
