FROM golang:1.11-alpine AS builder
COPY main.go .
RUN go build -o /slowserver

FROM alpine
COPY --from=builder /slowserver /slowserver
ENTRYPOINT ["/slowserver"]