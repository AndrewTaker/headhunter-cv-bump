FROM golang:1.25-bookworm AS builder

WORKDIR /app

ENV CGO_ENABLED=1

COPY go.mod go.sum go.work ./
COPY ./pkg ./pkg
COPY ./apps ./apps

RUN apt update && apt install -y build-essential

RUN go work sync
RUN go build -o /app/cli ./apps/cli
RUN go build -o /app/scheduler ./apps/scheduler

FROM debian:bookworm

# RUN apt update --no-cache && apt add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/cli /app/cli
COPY --from=builder /app/scheduler /app/scheduler
COPY .env ./

CMD ["./cli"]

