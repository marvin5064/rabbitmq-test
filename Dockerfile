FROM rabbitmq:3.7.7  

FROM golang:1.10

WORKDIR /go/src/github.com/marvin5064/rabbitmq-test

CMD make run