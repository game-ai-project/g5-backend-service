FROM golang:1.19.1-alpine as builder

WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM scratch as runner

ENV GIN_MODE=debug
COPY --from=builder /go/bin/app /app
EXPOSE 8000/tcp
CMD ["/app"]
