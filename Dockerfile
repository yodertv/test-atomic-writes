FROM golang:1.11-alpine
WORKDIR /go/src/test-atomic-writes
COPY . .
RUN go build -v
RUN mkdir /public
RUN /go/src/test-atomic-writes/test-atomic-writes > /public/index.html
