FROM golang:alpine
RUN mkdir /app
ADD . /app/
WORKDIR /app
CMD CGO_ENABLED=0 go test ./...
