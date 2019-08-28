FROM golang:1.11-alpine AS builder
COPY main.go .
RUN go build -o /slowserver

FROM scratch
COPY --from=builder /slowserver /slowserver
ENTRYPOINT ["/slowserver"]