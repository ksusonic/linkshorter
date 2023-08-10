FROM golang:1.20
RUN mkdir /app

ADD . /app/
WORKDIR /app

RUN go build -o main github.com/ksusonic/linkshorter/cmd

EXPOSE 8080
CMD ["/app/main"]
