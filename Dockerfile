FROM golang:1.19.1-alpine as builder

WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go build -o /go/bin/app

FROM alpine:3.16.2 as runner

WORKDIR /
ENV GIN_MODE=debug
COPY --from=builder /go/bin/app /app
EXPOSE 8000/tcp
CMD ["/app"]
