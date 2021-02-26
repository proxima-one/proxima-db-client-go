FROM golang:alpine
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN apk add build-base
CMD ["go", "test", "./..."]
