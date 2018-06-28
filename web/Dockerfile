FROM golang:alpine

COPY . .
RUN go build webserver.go

CMD ./webserver --port=$PORT --serving_host=$SERVING_HOST